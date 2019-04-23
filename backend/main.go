package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/harmony-one/demo-apps/backend/client"
	"github.com/harmony-one/demo-apps/backend/db"
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

	versionFlag = flag.Bool("version", false, "Output version info")
)

const (
	port = "30000"
)

func main() {
	flag.Parse()
	if *versionFlag {
		printVersion(os.Args[0])
	}

	fdb, err := fdb.NewFdb(*key, *project)

	if err != nil {
		log.Fatalf("Failed to create Fdb client: %v", err)
		os.Exit(1)
	}

	// Close FDB when done.
	defer fdb.CloseFdb()

	fdb.GetPlayers(true)

	player, err := restclient.GetPlayer(*ip, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
	} else {
		fmt.Printf("Player: %v\n", player)
	}

	winner, err := restclient.GetWinner(*ip, port)
	if err != nil {
		log.Fatalf("GetWinner Error: %v", err)
	} else {
		fmt.Printf("Winner: %v\n", winner)
	}

	// wait for the execution of smart contracts
	time.Sleep(2 * time.Second)

	player, err = restclient.GetPlayer(*ip, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
	} else {
		fmt.Printf("Player: %v\n", player)
	}

	os.Exit(0)
}
