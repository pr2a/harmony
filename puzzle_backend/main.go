package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-openapi/swag"
	"google.golang.org/appengine"
	app_log "google.golang.org/appengine/log"

	restclient "github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
	"github.com/harmony-one/demo-apps/backend/p2p"
	"github.com/harmony-one/demo-apps/backend/utils"
	"github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi"
	_ "github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi"
	"github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi/operations"
	_ "github.com/harmony-one/demo-apps/puzzle_backend/swagger/restapi/operations"
)

type respEnter struct {
	Address string `json:"address"`
	Level   uint   `json:"level"`
	Balance uint64 `json:"balance"`
}

type respFinish struct {
	Level   int    `json:"level"`
	Rewards uint64 `json:"rewards"`
}

type respReg struct {
	Account string `json:"account"`
	Private string `json:"private"`
}

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

	rpcDone chan (restclient.RPCMsg)
)

const (
	minimalFee = 1
	adminKey   = "e401343197a852f361e38ce6b46c99f1d6d1f80499864c6ae7effee42b46ab6b"
	dbKeyFile  = "./puzzle_backend/keys/benchmark_account_key.json"
	dbProject  = "benchmark-209420"
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

	rpcDone = make(chan (restclient.RPCMsg))

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
	defer server.Shutdown()

	server.Port = 30000 // TODO ek – parametrize this

	api.PostRegHandler = operations.PostRegHandlerFunc(handlePostReg)
	api.PostPlayHandler = operations.PostPlayHandlerFunc(handlePostPlay)
	api.PostFinishHandler = operations.PostFinishHandlerFunc(handlePostFinish)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

	appengine.Main()
}

func handlePostReg(params operations.PostRegParams) middleware.Responder {
	return operations.NewPostRegCreated().WithAccessControlAllowOrigin(
		"*", // TODO ek – tighten this up!
	).WithPayload(
		&operations.PostRegCreatedBody{
			Account: "0x0000000000000000000000000000000000000000",
			Email:   "ek@harmony.one",
		},
	)
}

func handlePostPlay(params operations.PostPlayParams) middleware.Responder {
	return operations.NewPostPlayCreated()
}

func handlePostFinish(params operations.PostFinishParams) middleware.Responder {
	return operations.NewPostFinishOK().WithPayload(
		&operations.PostFinishOKBody{
			Reward: 5e+18,
		},
	)
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/api/v1/reg" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	ids, ok := q["id"]
	if !ok {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	id := ids[0]

	var account *fdb.PzPlayer
	// find the existing account from firebase DB
	accounts := db.FindAccount("email", id)

	// register the new account
	if len(accounts) == 0 { // didn't find the account
		// generate the key
		address, priv := utils.GenereateKeys()
		leader := restclient.PickALeader()

		go restclient.FundMe(leader, address, rpcDone)

		player := fdb.PzPlayer{
			Email:   id,
			CosID:   "133", //FIXME: this has to be an id
			PrivKey: priv,
			Address: address,
			Leader:  leader.IP,
			Port:    leader.Port,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			app_log.Criticalf(ctx, "regHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		account = &player
		fmt.Printf("register new Account: %v for email: %v\n", account, id)
	} else {
		// we should find only one account, if more than one, just get the first one
		account = accounts[0]
		fmt.Printf("found Account: %v for id: %v\n", account, id)
		go func() {
			msg := restclient.RPCMsg{
				Err:  nil,
				Done: true,
			}
			rpcDone <- msg
		}()
	}

	var res string

	select {
	case msg := <-rpcDone:
		if msg.Err != nil {
			http.Error(w, "FundMe Error. Please retry", http.StatusInternalServerError)
			return
		}
		resp := respReg{
			Account: account.Address,
			Private: account.PrivKey,
		}

		bytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("can't marshal reg resp: %s\n", resp)
			http.Error(w, "warning: can't marshal response", http.StatusNoContent)
			return
		}
		res = string(bytes)
		fmt.Printf("res: %v\n", res)
		break
		//TODO: send email to player
	}

	io.WriteString(w, res)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/api/v1/play" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	keys, ok := q["key"]
	if !ok {
		http.Error(w, "missing key params", http.StatusBadRequest)
		return
	}
	key := keys[0]

	stakes, ok := q["stake"]
	if !ok {
		http.Error(w, "missing stake params", http.StatusBadRequest)
		return
	}
	stake := stakes[0]

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		http.Error(w, "can't find the registered player", http.StatusBadRequest)
		return

	}

	// we should find only one account, if more than one, just get the first one
	account := accounts[0]
	app_log.Infof(ctx, "player: %v is about to play", account.Address)

	// calling the play smart contract
	level, err := restclient.EnterPuzzle(account.Address, stake)
	if err != nil {
		app_log.Criticalf(ctx, "playHandler EnterPuzzle failed: %v", err)
		http.Error(w, "can't play the game, please retry", http.StatusInternalServerError)
		return
	}

	resp := respEnter{
		Address: account.Address,
		Level:   level,
		Balance: 0,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("can't marshal enter resp: %v\n", resp)
		http.Error(w, "Can't marshal enter response", http.StatusInternalServerError)
		return
	}
	res := string(bytes)
	_, _ = io.WriteString(w, res)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/finish" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	var ok bool

	newlevels, ok := q["newlevel"]
	if !ok {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}
	newlevel, err := strconv.Atoi(newlevels[0])
	if err != nil {
		http.Error(w, "wrong parameters", http.StatusBadRequest)
		return
	}

	accounts, ok := q["account"]
	if !ok {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}
	account := accounts[0]
	keys, ok := q["key"]
	if !ok || keys[0] != adminKey {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	rewards, err := restclient.GetRewards(account, newlevel)
	if err != nil {
		app_log.Criticalf(ctx, "finishHandler GetRewards error: %v", err)
		http.Error(w, "Can't Get Rewards", http.StatusInternalServerError)
		return
	}

	resp := respFinish{
		Level:   newlevel,
		Rewards: rewards,
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("can't marshal finish resp: %v\n", resp)
		http.Error(w, "Can't marshal finish response", http.StatusInternalServerError)
		return
	}
	res := string(bytes)
	_, _ = io.WriteString(w, res)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
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
		app_log.Infof(ctx, "accounts: %v", accounts)
		res = fmt.Sprintf("accounts: %v\n", accounts)
	case "RegisterAccount":
		account, priv := utils.GenereateKeys()
		emails, ok := q["email"]
		if !ok {
			http.Error(w, "missing email params", http.StatusBadRequest)
			break
		}
		app_log.Infof(ctx, "accounts: %v/%v", account, priv)
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
			app_log.Criticalf(ctx, "playHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		res = fmt.Sprintf("accounts: %v\n", account)
	}
	io.WriteString(w, res)
}
