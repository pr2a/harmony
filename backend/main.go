package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
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
	collection = flag.String("collection", "players", "name of collection")
	key        = flag.String("key", "./keys/benchmark_account_key.json", "key filename")
	project    = flag.String("project", "lottery-demo-leo", "project ID of firebase")
	ip         = flag.String("server_ip", "34.222.210.98", "the IP address of the server")
	action     = flag.String("action", "reg", "action of the program. Valid (reg, winner)")

	versionFlag = flag.Bool("version", false, "Output version info")

	db *fdb.Fdb
)

const (
	port = "30000"
)

func sendRegEmail() {
	if db == nil {
		fmt.Printf("[sendRegEmail] fdb is nil")
		return
	}

	players := db.GetPlayers("email_key", "==", false)

	for i, p := range players {
		fmt.Printf("[sendRegEmail] %v => %v\n", i, p)
	}

	//TODO: send an email to user for the account registration
}

func pickWinner() {
	if db == nil {
		fmt.Printf("[pickWinner] fdb is nil")
		return
	}

	//Get a list of all current players
	player, err := restclient.GetPlayer(*ip, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
	} else {
		fmt.Printf("Player: %v\n", player)
	}

	/*
		session := db.GetSession(false)
		if len(session) > 0 {
			fmt.Printf("Current Session ID: %v\n", session[0].ID)
		} else {
			fmt.Printf("Get No Session\n")
		}
	*/

	//Get all player from DB
	players := db.GetPlayers("", "", nil)
	for i, p := range players {
		fmt.Printf("%v => %v\n", i, p)
	}

	//TODO: get balances of all players

	//Run the get winner smart contract
	winner, err := restclient.GetWinner(*ip, port)
	if err != nil {
		log.Fatalf("GetWinner Error: %v", err)
	} else {
		fmt.Printf("Winner: %v\n", winner)
	}

	// wait for the execution of smart contracts
	time.Sleep(2 * time.Second)

	//TODO: check the balances of all players again

	player, err = restclient.GetPlayer(*ip, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
	} else {
		fmt.Printf("Player: %v\n", player)
	}

	return
}

func main() {
	flag.Parse()
	if *versionFlag {
		printVersion(os.Args[0])
	}

	var err error
	db, err = fdb.NewFdb(*key, *project)

	if err != nil {
		log.Fatalf("Failed to create Fdb client: %v", err)
		os.Exit(1)
	}

	// Close FDB when done.
	defer db.CloseFdb()

	switch *action {
	case "reg":
		sendRegEmail()
	case "winner":
		pickWinner()
	default:
		fmt.Printf("Wrong action: %v", action)
	}

	os.Exit(0)
}
