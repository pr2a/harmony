package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/go-openapi/swag"
	"google.golang.org/appengine"
	app_log "google.golang.org/appengine/log"

	restclient "github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
	"github.com/harmony-one/demo-apps/backend/p2p"
	"github.com/harmony-one/demo-apps/backend/utils"
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

	http.HandleFunc("/reg", handlePostReg)
	http.HandleFunc("/play", handlePostPlay)
	http.HandleFunc("/finish", handlePostFinish)

	leader = readProfile(*profile)

	leaders := make([]p2p.Peer, 0)
	for _, ldr := range backendProfile.RPCServer {
		leaders = append(leaders, p2p.Peer{IP: ldr[0].IP, Port: defaultPort})
	}
	restclient.SetLeaders(leaders)

	appengine.Main()
}

type cosGetUIDRequestBody struct {
	Token     string `json:"token"`     // Temporary COS login token.
	Timestamp int64  `json:"timestamp"` // Current UNIX timestamp, in ms.
	ClientID  string `json:"client_id"` // App-specific client ID.
}

type cosGetUIDResponseBodyData struct {
	UID string `json:"uid"` // COS UID.
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

type msgBody struct {
	Msg string `json:"msg"`
}

type postRegResponseBody struct {
	Account string `json:"address"`
	PrivKey string `json:"privkey"`
	UID     string `json:"uid"`
	Txid    string `json:"txid"`
	Balance string `json:"balance"`
}

func handlePostReg(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/reg" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	tokens, ok := q["token"]
	if !ok {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}
	token := tokens[0]
	uid, err := getUID(token)
	if err != nil {
		app_log.Infof(ctx, "handlePostReg: getUID returned %#v", err)
		http.Error(w, "", http.StatusUnauthorized)
	}
	app_log.Infof(ctx, "handlePostReg: UID %#v logging in", uid)

	rpcDone := make(chan (restclient.RPCMsg))

	var account *fdb.PzPlayer
	// find the existing account from firebase DB
	accounts := db.FindAccount("cosid", uid)
	var resCode int

	// register the new account
	var resBody postRegResponseBody
	if len(accounts) == 0 { // didn't find the account
		// generate the key
		resBody.Account, resBody.PrivKey = utils.GenereateKeys()
		// TODO ek – fix this later somehow...
		resBody.Balance = "10000000000000000000"
		leader := restclient.PickALeader()

		go restclient.FundMe(leader, resBody.Account, rpcDone)

		player := fdb.PzPlayer{
			Email:   "",
			CosID:   uid,
			PrivKey: resBody.PrivKey,
			Address: resBody.Account,
			Leader:  leader.IP,
			Port:    leader.Port,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			app_log.Infof(ctx, "handlePostReg registerAccount error: %v", err)
			http.Error(w, "register account failure", http.StatusServiceUnavailable)
			return
		}
		account = &player
		app_log.Infof(ctx, "handlePostReg: register new Account: %v for cosid: %v", account, uid)
		if msg := <-rpcDone; msg.Err != nil {
			http.Error(w, "fund me failure", http.StatusGatewayTimeout)
			return
		} else {
			resBody.Txid = msg.Txid
		}

		//TODO: send email to player
		go func() {
			app_log.Infof(ctx, "Sent email ..")
		}()

		resCode = http.StatusCreated
	} else {
		// we should find only one account, if more than one, just get the first one
		account = accounts[0]
		resBody.Account = account.Address
		resBody.PrivKey = account.PrivKey
		resCode = http.StatusOK

		leader = p2p.Peer{
			IP:   account.Leader,
			Port: account.Port,
		}

		chanBalanceMsg := make(chan restclient.AccountBalanceMsg)
		go restclient.GetBalance(leader, account.Address, chanBalanceMsg)
		balanceMsg := <-chanBalanceMsg
		if balanceMsg.Err != nil {
			app_log.Infof(ctx, "get balance failure: %#v", balanceMsg.Err)
			http.Error(w, "get balance failure", http.StatusGatewayTimeout)
			return
		}
		resBody.Balance = balanceMsg.Balance
	}

	jsonResp(ctx, w, resCode, resBody)
}

type postPlayResponseBody struct {
	Txid string `json:"txid"`
}

func handlePostPlay(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/play" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	keys, ok := q["accountKey"]
	if !ok {
		http.Error(w, "missing account key", http.StatusBadRequest)
		return
	}
	key := keys[0]
	stakes, ok := q["stake"]
	if !ok {
		http.Error(w, "missing account key", http.StatusBadRequest)
		return
	}
	stake := stakes[0]

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		http.NotFound(w, r)
		return
	}
	account := accounts[0]
	app_log.Infof(ctx, "player: %v is about to play", account.Address)
	leader := p2p.Peer{
		IP:   account.Leader,
		Port: account.Port,
	}

	rpcDone := make(chan (restclient.RPCMsg))
	go restclient.PlayGame(leader, key, fmt.Sprintf("%v", stake), rpcDone)
	msg := <-rpcDone
	if msg.Err != nil {
		app_log.Infof(ctx, "playHandler PlayGame failed: %v", msg.Err)
		http.Error(w, "play failure", http.StatusGatewayTimeout)
		return
	}

	jsonResp(ctx, w, http.StatusCreated, &postPlayResponseBody{
		Txid: msg.Txid,
	})
}

type postFinishResponseBody struct {
	Reward string `json:"reward"`
	Txid   string `json:"txid"`
}

func handlePostFinish(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/play" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	keys, ok := q["accountKey"]
	if !ok {
		http.Error(w, "missing account key", http.StatusBadRequest)
		return
	}
	key := keys[0]

	heightStrs, ok := q["height"]
	if !ok {
		http.Error(w, "missing height", http.StatusBadRequest)
		return
	}
	height, err := strconv.ParseInt(heightStrs[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid height", http.StatusBadRequest)
		return
	}

	sequences, ok := q["sequence"]
	if !ok {
		http.Error(w, "missing sequence", http.StatusBadRequest)
		return
	}
	sequence := sequences[0]

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		http.NotFound(w, r)
		return
	}

	account := accounts[0]
	app_log.Infof(ctx, "player: %v/%v is about to get paid", account.Address, height)

	leader := p2p.Peer{
		IP:   account.Leader,
		Port: account.Port,
	}

	rpcDone := make(chan (restclient.RPCMsg))
	go restclient.PayOut(leader, key, height, sequence, rpcDone)
	msg := <-rpcDone

	if msg.Err != nil {
		app_log.Infof(ctx, "/finish PayOut failed: %v", msg.Err)
		http.Error(w, "payout failure", http.StatusGatewayTimeout)
		return
	}

	rpcEndDone := make(chan (restclient.RPCMsg))
	go restclient.EndGame(leader, key, rpcEndDone)
	msgEnd := <-rpcEndDone

	if msgEnd.Err != nil {
		app_log.Infof(ctx, "/finish EndGame failed: %v", msgEnd.Err)
		http.Error(w, "endgame failure", http.StatusGatewayTimeout)
		return
	}

	jsonResp(ctx, w, http.StatusOK, postFinishResponseBody{
		Reward: "",
		Txid:   msg.Txid,
	})
}

func jsonResp(
	ctx context.Context, w http.ResponseWriter, code int, res interface{},
) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept")
	resBytes, err := json.Marshal(res)
	if err != nil {
		app_log.Errorf(ctx, "cannot marshal response %#v: %v", res, err)
		http.Error(w, "", http.StatusInternalServerError)
	}
	_, _ = w.Write(resBytes)
}
