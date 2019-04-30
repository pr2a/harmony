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
	defaultConfigFile = ".hmy/backend.ini"
	defaultProfile    = "default"
	defaultPort       = "30000"
	leader            p2p.Peer
	backendProfile    *utils.BackendProfile

	profile     = flag.String("profile", defaultProfile, "name of the profile")
	versionFlag = flag.Bool("version", false, "Output version info")
)

const (
	minimalFee = 1
	gameFee    = 1
	adminKey   = "e401343197a852f361e38ce6b46c99f1d6d1f80499864c6ae7effee42b46ab6b"
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

	http.HandleFunc("/enter", enterHandler)
	http.HandleFunc("/finish", finishHandler)

	leader = readProfile(*profile)

	//Get a list of all current players
	_, err := restclient.GetPlayer(leader.IP, defaultPort)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
		return
	}

	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func enterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/enter" {
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

	// find the existing account from firebase DB
	account, leader := fdb.FindAccount(email)
	_ = leader

	// register the new account
	if account == "" {
		// generate the key
		account, _ = utils.GenereateAccount(email)
		err := restclient.FundMe(account)
		if err != nil {
			app_log.Criticalf(ctx, "enterHandler FundMe error: %v", err)
			http.Error(w, "FundMe Error, please retry", http.StatusInternalServerError)
			// TODO: retry
			return
		}
		leaders := restclient.GetLeaders()
		err = fdb.RegisterAccount(email, account, leaders[0].IP)
		if err != nil {
			app_log.Criticalf(ctx, "enterHandler registerAccount error: %v", err)
			http.Error(w, "Register Account, please retry", http.StatusInternalServerError)
		}
	}
	fmt.Printf("found Account: %v for email: %v\n", account, email)

	balance, err := restclient.GetBalance(account)
	if err != nil {
		app_log.Criticalf(ctx, "enterHandler GetBalance error: %v", err)
		http.Error(w, "Can't GetBalance, please retry", http.StatusInternalServerError)
		return
	}
	if balance < minimalFee {
		http.Error(w, "Not enough balance to play, please get more token", http.StatusBadRequest)
		return
	}

	level, err := restclient.EnterPuzzle(account, gameFee)
	if err != nil {
		app_log.Criticalf(ctx, "enterHandler EnterPuzzle error: %v", err)
		http.Error(w, "Can't Enter Game, please retry", http.StatusInternalServerError)
	}

	resp := respEnter{
		Address: account,
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
