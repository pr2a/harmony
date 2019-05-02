package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-openapi/swag"

	restclient "github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
	"github.com/harmony-one/demo-apps/backend/p2p"
	"github.com/harmony-one/demo-apps/backend/utils"
	"github.com/harmony-one/demo-apps/puzzle_backend/swagger/models"
	"github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi"
	_ "github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi"
	"github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi/operations"
	_ "github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi/operations"
)

var (
	version string
	builtBy string
	builtAt string
	commit  string
)

func printVersion(me string) {
	fmt.Fprintf(os.Stderr, "Harmony (C) 2019. %v, version %v-%v (%v %v)\n", path.Base(me), version, commit, builtBy, builtAt)
	os.Exit(0)
}

var (
	defaultConfigFile = "./puzzle_backend/.hmy/backend.ini"
	defaultProfile    = "default"
	defaultPort       = "30000"
	leader            p2p.Peer
	backendProfile    *utils.BackendProfile

	db *fdb.Fdb

	profile     = flag.String("profile", defaultProfile, "name of the profile")
	versionFlag = flag.Bool("version", false, "Output version info")
)

const (
	adminKey  = "e401343197a852f361e38ce6b46c99f1d6d1f80499864c6ae7effee42b46ab6b"
	dbKeyFile = "./puzzle_backend/keys/benchmark_account_key.json"
	dbProject = "benchmark-209420"
)

// readProfile read the ini file and return the leader's IP
func readProfile(profile string) p2p.Peer {
	fmt.Printf("Using %s profile for backend\n", profile)
	var err error
	backendProfile, err = utils.ReadBackendProfile(defaultConfigFile, profile)
	if err != nil {
		fmt.Printf("Read backend profile error: %v\nExiting ...\n", err)
		os.Exit(2)
	}

	return backendProfile.RPCServer[0][0]
}

func main() {
	flag.Parse()
	if *versionFlag {
		printVersion(os.Args[0])
	}

	var err error
	db, err = fdb.NewFdb(dbKeyFile, dbProject)

	if err != nil || db == nil {
		log.Fatalf("Failed to create Fdb client: %v", err)
		os.Exit(1)
	}

	// Close FDB when done.
	defer db.CloseFdb()

	leader = readProfile(*profile)

	leaders := make([]p2p.Peer, 0)
	for _, ldr := range backendProfile.RPCServer {
		leaders = append(leaders, p2p.Peer{IP: ldr[0].IP, Port: defaultPort})
	}
	restclient.SetLeaders(leaders)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHarmonyPuzzleAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.EnabledListeners = []string{"http"}
	defer func() { _ = server.Shutdown() }()

	server.Port = 30000 // TODO ek – parametrize this

	api.PostRegHandler = operations.PostRegHandlerFunc(handlePostReg)
	api.PostPlayHandler = operations.PostPlayHandlerFunc(handlePostPlay)
	api.PostFinishHandler = operations.PostFinishHandlerFunc(handlePostFinish)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

	//appengine.Main()
}

type cosGetUIDRequestBody struct {
	Token     string `json:"token"`     // Temporary COS login token.
	Timestamp int64  `json:"timestamp"` // Current UNIX timestamp, in ms.
	ClientID  string `json:"client_id"` // App-specific client ID.
}

type cosGetUIDResponseBodyData struct {
	UID string `json:"string"` // COS UID.
}

type cosGetUIDResponseBody struct {
	Status  int64                     `json:"status"`  // Status code.
	Message string                    `json:"message"` // Status message.
	Data    cosGetUIDResponseBodyData `json:"data"`    // Data (“meat”).
}

type getUIDError struct {
	Status  int64
	Message string
}

func (e *getUIDError) Error() string {
	return fmt.Sprintf("get_uid failed with status=%#v message=%#v",
		e.Status, e.Message)
}

func getUID(token string) (string, error) {
	ts := int64(time.Now().UnixNano() / 1000000)
	reqBodyBytes, err := json.Marshal(cosGetUIDRequestBody{
		Token:     token,
		Timestamp: ts,
		ClientID:  "3",
	})
	if err != nil {
		return "", err
	}

	h := md5.New()
	h.Write([]byte("5VCjkUpHkueWo77S1TJC8d3dAgDry0pitRIGliIbucE=")) // TODO ek parametrize
	h.Write(reqBodyBytes)
	if _, err := fmt.Fprint(h, ts); err != nil {
		return "", err
	}
	auth := hexutil.Encode(h.Sum(nil))[2:]

	cosClient := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"http://qa.contentos.io/api/v1/open/get_uid",
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	res, err := cosClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() { _ = res.Body.Close() }()
	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resBody := cosGetUIDResponseBody{}
	err = json.Unmarshal(resBodyBytes, &resBody)
	if err != nil {
		return "", err
	}
	if resBody.Status != 1 {
		return "", &getUIDError{resBody.Status, resBody.Message}
	}
	return resBody.Data.UID, nil
}

func handlePostReg(params operations.PostRegParams) middleware.Responder {
	uid, err := getUID(params.Token)
	if err != nil {
		middleware.Logger.Printf("handlePostReg: getUID returned %#v", err)
		return operations.NewPostRegUnauthorized()
	}

	rpcDone := make(chan (restclient.RPCMsg))

	var account *fdb.PzPlayer
	// find the existing account from firebase DB
	accounts := db.FindAccount("cosid", uid)
	var resFunc func(payload *models.PostRegResponse) middleware.Responder

	var balance string

	// register the new account
	if len(accounts) == 0 { // didn't find the account
		// generate the key
		address, priv := utils.GenereateKeys()
		leader := restclient.PickALeader()

		go restclient.FundMe(leader, address, rpcDone)

		player := fdb.PzPlayer{
			Email:   "",
			CosID:   uid,
			PrivKey: priv,
			Address: address,
			Leader:  leader.IP,
			Port:    leader.Port,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			middleware.Logger.Printf("handlePostReg registerAccount error: %v", err)
			return operations.NewPostRegServiceUnavailable().WithPayload(
				&operations.PostRegServiceUnavailableBody{
					Msg: "register account failure",
				},
			)
		}
		account = &player
		middleware.Logger.Printf("handlePostReg: register new Account: %v for cosid: %v", account, uid)
		if msg := <-rpcDone; msg.Err != nil {
			return operations.NewPostRegGatewayTimeout().WithPayload(
				&operations.PostRegGatewayTimeoutBody{
					Msg: "fund me failure",
				},
			)
		}

		//TODO: send email to player
		go func() {
			middleware.Logger.Printf("Sent email ..")
		}()

		resFunc = func(payload *models.PostRegResponse) middleware.Responder {
			return operations.NewPostRegCreated().
				WithAccessControlAllowOrigin("*").WithPayload(payload)
		}
		// TODO ek – change this later.
		//  We should probably let FE query the balance independently,
		//  and not query the balance here.
		//  A balance query will fail here because fundme is likely to be
		//  still in flight.
		balance = "10000000000000000000"
	} else {
		// we should find only one account, if more than one, just get the first one
		account = accounts[0]

		leader = p2p.Peer{
			IP:   account.Leader,
			Port: account.Port,
		}

		chanBalanceMsg := make(chan restclient.AccountBalanceMsg)
		go restclient.GetBalance(leader, account.Address, chanBalanceMsg)
		balanceMsg := <-chanBalanceMsg
		if balanceMsg.Err != nil {
			middleware.Logger.Printf("get balance failure: %#v", balanceMsg.Err)
			return operations.NewPostRegGatewayTimeout().WithPayload(
				&operations.PostRegGatewayTimeoutBody{
					Msg: "get balance failure",
				},
			)
		}

		balance = balanceMsg.Balance
		resFunc = func(payload *models.PostRegResponse) middleware.Responder {
			return operations.NewPostRegOK().WithPayload(payload)
		}
	}

	return resFunc(
		&models.PostRegResponse{
			Account: account.Address,
			UID:     uid,
			Balance: balance,
		},
	)
}

func handlePostPlay(params operations.PostPlayParams) middleware.Responder {
	key := params.AccountKey
	stake := params.Stake

	rpcDone := make(chan (restclient.RPCMsg))

	_ = stake

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		return operations.NewPostPlayNotFound()
	}
	account := accounts[0]
	middleware.Logger.Printf("player: %v is about to play", account.Address)
	leader := p2p.Peer{
		IP:   account.Leader,
		Port: account.Port,
	}

	go restclient.PlayGame(leader, key, fmt.Sprintf("%v", stake), rpcDone)

	select {
	case msg := <-rpcDone:

		if msg.Err != nil {
			middleware.Logger.Printf("playHandler PlayGame failed: %v", msg.Err)
			return operations.NewPostPlayGatewayTimeout().WithPayload(
				&operations.PostPlayGatewayTimeoutBody{
					Msg: "play failure",
				},
			)
		}
		break
	}

	return operations.NewPostPlayCreated()
}

func handlePostFinish(params operations.PostFinishParams) middleware.Responder {
	key := params.AccountKey

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		return operations.NewPostPlayNotFound()
	}

	account := accounts[0]
	middleware.Logger.Printf("player: %v/%v is about to get paid", account.Address, params.Height)

	leader := p2p.Peer{
		IP:   account.Leader,
		Port: account.Port,
	}

	rpcDone := make(chan (restclient.RPCMsg))
	go restclient.PayOut(leader, key, *params.Height, params.Sequence, rpcDone)
	msg := <-rpcDone

	if msg.Err != nil {
		middleware.Logger.Printf("/finish PayOut failed: %v", msg.Err)
		return operations.NewPostFinishGatewayTimeout().WithPayload(
			&operations.PostFinishGatewayTimeoutBody{
				Msg: "payout failure",
			},
		)
	}

	return operations.NewPostFinishOK().WithPayload(
		&operations.PostFinishOKBody{
			Reward: "",
			Txid:   msg.Txid,
		},
	)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/test" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	var ok bool
	var res string

	function, ok := q["function"]
	if !ok {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	switch function[0] {
	case "FindAccount":
		keys, ok := q["key"]
		if !ok {
			http.Error(w, "missing key params", http.StatusBadRequest)
			break
		}
		values, ok := q["value"]
		if !ok {
			http.Error(w, "missing value params", http.StatusBadRequest)
			break
		}
		accounts := db.FindAccount(keys[0], values[0])
		middleware.Logger.Printf("accounts: %v", accounts)
		res = fmt.Sprintf("accounts: %v\n", accounts)
	case "RegisterAccount":
		account, priv := utils.GenereateKeys()
		emails, ok := q["email"]
		if !ok {
			http.Error(w, "missing email params", http.StatusBadRequest)
			break
		}
		middleware.Logger.Printf("accounts: %v/%v", account, priv)
		player := fdb.PzPlayer{
			Email:   emails[0],
			CosID:   "133",
			PrivKey: priv,
			Address: account,
			Leader:  "192.168.192.1",
			Port:    defaultPort,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			middleware.Logger.Printf("playHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		res = fmt.Sprintf("accounts: %v\n", account)
	}
	io.WriteString(w, res)
}
