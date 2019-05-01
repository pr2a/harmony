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
	"sync"

	"github.com/harmony-one/demo-apps/backend/client"
	"github.com/harmony-one/demo-apps/backend/db"
	"github.com/harmony-one/demo-apps/backend/p2p"
	"github.com/harmony-one/demo-apps/backend/utils"
	"google.golang.org/appengine"
	app_log "google.golang.org/appengine/log"
)

type respEnter struct {
	Address string `json:address`
	Level   uint   `json:level`
	Balance uint64 `json:balance`
}

type respFinish struct {
	Level   int    `json:level`
	Rewards uint64 `json:rewards`
}

type respReg struct {
	Account string `json:account`
	Private string `json:private`
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
)

const (
	minimalFee = 1
	gameFee    = 1
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

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/enter", enterHandler)
	http.HandleFunc("/finish", finishHandler)
	http.HandleFunc("/test", testHandler)

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

	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
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
	accounts := db.FindAccount(id)

	var wg sync.WaitGroup
	// register the new account
	if len(accounts) == 0 {
		// generate the key
		address, priv := utils.GenereateKeys()
		leader := restclient.PickALeader()

		wg.Add(1)
		go restclient.FundMe(leader, address, wg)

		player := fdb.PzPlayer{
			Email:   id,
			CosID:   "133",
			PrivKey: priv,
			Address: address,
			Leader:  leader.IP,
			Port:    leader.Port,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			app_log.Criticalf(ctx, "enterHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		account = &player
		fmt.Printf("register new Account: %v for email: %v\n", account, id)
	} else {
		// we should find only one account, if more than one, just get the first one
		account := accounts[0]
		fmt.Printf("found Account: %v for id: %v\n", account, id)
	}

	resp := respReg{
		Account: account.Address,
		Private: account.PrivKey,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("can't marshal enter resp: %s\n", resp)
		http.Error(w, "Can't marshal enter response", http.StatusInternalServerError)
		return
	}
	res := string(bytes)
	io.WriteString(w, res)
}

func enterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/api/v1/enter" {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()

	emails, ok := q["email"]
	if !ok {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}
	email := emails[0]

	var account *fdb.PzPlayer
	// find the existing account from firebase DB
	accounts := db.FindAccount(email)

	var wg sync.WaitGroup
	// register the new account
	if len(accounts) == 0 {
		// generate the key
		address, priv := utils.GenereateKeys()
		leader := restclient.PickALeader()

		wg.Add(1)
		go restclient.FundMe(leader, address, wg)

		player := fdb.PzPlayer{
			Email:   email,
			CosID:   "133",
			PrivKey: priv,
			Address: address,
			Leader:  leader.IP,
			Port:    leader.Port,
		}
		err := db.RegisterAccount(&player)
		if err != nil {
			app_log.Criticalf(ctx, "enterHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		account = &player
		fmt.Printf("register new Account: %v for email: %v\n", account, email)
	} else {
		// we should find only one account, if more than one, just get the first one
		account := accounts[0]
		fmt.Printf("found Account: %v for email: %v\n", account, email)
	}

	wg.Wait()
	balance, err := restclient.GetBalance(account.Address)
	if err != nil {
		app_log.Criticalf(ctx, "enterHandler GetBalance error: %v", err)
		http.Error(w, "Can't GetBalance, please retry", http.StatusInternalServerError)
		return
	}
	if balance < minimalFee {
		http.Error(w, "Not enough balance to play, please get more token", http.StatusBadRequest)
		return
	}

	level, err := restclient.EnterPuzzle(account.Address, gameFee)
	if err != nil {
		app_log.Criticalf(ctx, "enterHandler EnterPuzzle error: %v", err)
		http.Error(w, "Can't Enter Game, please retry", http.StatusInternalServerError)
	}

	resp := respEnter{
		Address: account.Address,
		Level:   level,
		Balance: balance,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("can't marshal enter resp: %s\n", resp)
		http.Error(w, "Can't marshal enter response", http.StatusInternalServerError)
		return
	}
	res := string(bytes)
	io.WriteString(w, res)
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
		fmt.Printf("can't marshal finish resp: %s\n", resp)
		http.Error(w, "Can't marshal finish response", http.StatusInternalServerError)
		return
	}
	res := string(bytes)
	io.WriteString(w, res)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/test" {
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
		emails, ok := q["email"]
		if !ok {
			http.Error(w, "missing email params", http.StatusBadRequest)
			break
		}
		accounts := db.FindAccount(emails[0])
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
			app_log.Criticalf(ctx, "enterHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
		res = fmt.Sprintf("accounts: %v\n", account)
	}
	io.WriteString(w, res)
}
