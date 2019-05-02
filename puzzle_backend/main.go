package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

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
	ctx := appengine.NewContext(params.HTTPRequest)
	id := params.Email

	rpcDone := make(chan (restclient.RPCMsg))

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
			app_log.Criticalf(ctx, "handlePostReg registerAccount error: %v", err)
			return operations.NewPostRegServiceUnavailable().WithPayload(
				&operations.PostRegServiceUnavailableBody{
					Msg: "register account failure",
				},
			)
		}
		account = &player
		fmt.Printf("register new Account: %v for email: %v\n", account, id)
		select {
		case msg := <-rpcDone:
			if msg.Err != nil {
				return operations.NewPostRegGatewayTimeout().WithPayload(
					&operations.PostRegGatewayTimeoutBody{
						Msg: "fund me failure",
					},
				)
			}
			//TODO: send email to player
			go func() {
				fmt.Println("Sent email ..")
			}()

			return operations.NewPostRegCreated().WithAccessControlAllowOrigin("*").WithPayload(
				&operations.PostRegCreatedBody{
					Account: account.Address,
					Email:   id,
					Balance: "10000000000000000000", // TOOD: placeholder
				},
			)
		}

	} else {
		// we should find only one account, if more than one, just get the first one
		account = accounts[0]
		// TODO: leo
		//		go restclient.GetBalance(account.Leader, account.Address, rpcDone)
		//		fmt.Printf("found existing Account: %v for id: %v\n", account, id)

		return operations.NewPostRegOK().WithPayload(
			&operations.PostRegOKBody{
				Account: account.Address,
				Email:   id,
				Balance: "900000000000000000", // TOOD: placeholder
			},
		)
	}
}

func handlePostPlay(params operations.PostPlayParams) middleware.Responder {
	ctx := appengine.NewContext(params.HTTPRequest)

	key := params.AccountKey
	stake := params.Stake

	_ = stake

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		return operations.NewPostPlayNotFound()
	}
	account := accounts[0]
	fmt.Printf("player: %v is about to play\n", account.Address)

	err := restclient.EnterPuzzle(account.Address, fmt.Sprintf("%v", stake))
	if err != nil {
		app_log.Criticalf(ctx, "playHandler EnterPuzzle failed: %v", err)
		return operations.NewPostPlayGatewayTimeout().WithPayload(
			&operations.PostPlayGatewayTimeoutBody{
				Msg: "play failure",
			},
		)
	}

	return operations.NewPostPlayCreated()
}

func handlePostFinish(params operations.PostFinishParams) middleware.Responder {
	ctx := appengine.NewContext(params.HTTPRequest)

	key := params.AccountKey

	// find the existing account from firebase DB
	accounts := db.FindAccount("privkey", key)

	// can't play if player didn't register before
	if len(accounts) == 0 {
		return operations.NewPostPlayNotFound()
	}
	account := accounts[0]
	fmt.Printf("player: %v/%v is about to get paid\n", account.Address, params.Height)

	_, err := restclient.GetRewards(account.Address, *params.Height)
	if err != nil {
		app_log.Criticalf(ctx, "finishHandler GetRewards failed: %v", err)
		return operations.NewPostFinishGatewayTimeout().WithPayload(
			&operations.PostFinishGatewayTimeoutBody{
				Msg: "finish failure",
			},
		)
	}

	return operations.NewPostFinishOK().WithPayload(
		&operations.PostFinishOKBody{
			Reward: 5e+18,
		},
	)
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
