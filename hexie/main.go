package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/harmony-one/demo-apps/backend/client"
	"github.com/harmony-one/demo-apps/backend/p2p"
	"github.com/harmony-one/demo-apps/backend/utils"
	"google.golang.org/appengine"
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
	defaultConfigFile = ".hmy/backend.ini"
	defaultProfile    = "default"
	defaultPort       = "30000"
	leader            p2p.Peer
	backendProfile    *utils.BackendProfile

	profile     = flag.String("profile", defaultProfile, "name of the profile")
	versionFlag = flag.Bool("version", false, "Output version info")
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
	if r.URL.Path != "/enter" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "EnterHandler!")
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/finish" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintln(w, "FinishHandler!")
}
