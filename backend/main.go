package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/jinzhu/copier"

	restclient "github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
	p2p "github.com/harmony-one/demo-apps/backend/p2p"
	clientService "github.com/harmony-one/demo-apps/backend/service"
	utils "github.com/harmony-one/demo-apps/backend/utils"
)

var (
	version string
	builtBy string
	builtAt string
	commit  string
)

// AccountState includes the balance and nonce of an account
type AccountState struct {
	balance *big.Int
	nonce   uint64
}

const (
	rpcRetry          = 3
	defaultConfigFile = ".hmy/backend.ini"
	defaultProfile    = "default"
	port              = "30000"
)

func printVersion(me string) {
	fmt.Fprintf(os.Stderr, "Harmony (C) 2019. %v, version %v-%v (%v %v)\n", path.Base(me), version, commit, builtBy, builtAt)
	os.Exit(0)
}

var (
	profile    = flag.String("profile", defaultProfile, "name of the profile")
	collection = flag.String("collection", "players", "name of collection")
	key        = flag.String("key", "./keys/leo_account_key.json", "key filename")
	project    = flag.String("project", "lottery-demo-leo", "project ID of firebase")
	action     = flag.String("action", "player", "action of the program. Valid (player, reg, winner, notify, players, balances)")
	verbose    = flag.Bool("verbose", true, "verbose log print at every step")

	versionFlag = flag.Bool("version", false, "Output version info")

	db             *fdb.Fdb
	backendProfile *utils.BackendProfile
	leader         p2p.Peer
	allPlayers     []*fdb.Player
)

// FetchBalance fetches account balance of specified address from the Harmony network
func FetchBalance(address common.Address) map[uint32]AccountState {
	result := make(map[uint32]AccountState)
	for i := 0; i < backendProfile.Shards; i++ {
		balance := big.NewInt(0)
		var nonce uint64

		result[uint32(i)] = AccountState{balance, 0}

		for retry := 0; retry < rpcRetry; retry++ {
			server := backendProfile.RPCServer[i][rand.Intn(len(backendProfile.RPCServer[i]))]
			client, err := clientService.NewClient(server.IP, server.Port)
			if err != nil {
				continue
			}

			response, err := client.GetBalance(address)
			if err != nil {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			balance.SetBytes(response.Balance)
			nonce = response.Nonce
			break
		}
		result[uint32(i)] = AccountState{balance, nonce}
	}
	return result
}

func processBalancesCommand(players []*fdb.Player) {
	for _, player := range players {
		if player == nil {
			continue
		}
		addr := common.HexToAddress(player.Address)
		if *verbose {
			fmt.Printf("Address: %s\n", player.Address)
		}
		// assuming number of shard is 1
		for shardID, balanceNonce := range FetchBalance(addr) {
			if *verbose {
				fmt.Printf("Balance in Shard %d:  %s, nonce: %v \n", shardID, convertBalanceIntoReadableFormat(balanceNonce.balance), balanceNonce.nonce)
			}
			player.Balance = balanceNonce.balance
		}
	}
}

func convertBalanceIntoReadableFormat(balance *big.Int) string {
	balance = balance.Div(balance, big.NewInt(params.GWei))
	strBalance := fmt.Sprintf("%d", balance.Uint64())

	bytes := []byte(strBalance)
	hasDecimal := false
	for i := 0; i < 11; i++ {
		if len(bytes)-1-i < 0 {
			bytes = append([]byte{'0'}, bytes...)
		}
		if bytes[len(bytes)-1-i] != '0' && i < 9 {
			hasDecimal = true
		}
		if i == 9 {
			newBytes := append([]byte{'.'}, bytes[len(bytes)-i:]...)
			bytes = append(bytes[:len(bytes)-i], newBytes...)
		}
	}
	zerosToRemove := 0
	for i := 0; i < len(bytes); i++ {
		if hasDecimal {
			if bytes[len(bytes)-1-i] == '0' {
				bytes = bytes[:len(bytes)-1-i]
				i--
			} else {
				break
			}
		} else {
			if zerosToRemove < 5 {
				bytes = bytes[:len(bytes)-1-i]
				i--
				zerosToRemove++
			} else {
				break
			}
		}
	}
	return string(bytes)
}

func sendRegEmail() {
	players := db.GetPlayers("email_key", "==", false)

	for i, p := range players {
		fmt.Printf("[sendRegEmail] %v => %v\n", i, p)
	}

	//TODO: send an email to user for the account registration
}

func pickWinner() {
	currentPlayers := getPlayer()
	existingPlayers := make([]*fdb.Player, 0)

	go getAllPlayer()

	for _, p := range currentPlayers {
		onePlayer := fdb.Player{}
		copier.Copy(&onePlayer, p)
		existingPlayers = append(existingPlayers, &onePlayer)
	}

	if *verbose {
		fmt.Printf("currentPlayers: %v\n", currentPlayers)
		fmt.Printf("existingPlayers: %v\n", existingPlayers)
	}

	//Run the get winner smart contract
	winner, err := restclient.GetWinner(leader.IP, port)
	if err != nil {
		log.Fatalf("GetWinner Error: %v", err)
	} else {
		fmt.Printf("Winner: %v\n", winner)
	}

	// wait for the execution of smart contracts
	time.Sleep(5 * time.Second)

	processBalancesCommand(existingPlayers)

	for i, p := range existingPlayers {
		if p == nil {
			continue
		}
		if p.Balance.Cmp(currentPlayers[i].Balance) > 0 {
			fmt.Printf("%s is the winner\n", p.Address)
			// TODO: send the email to winner
		} else {
			fmt.Printf("%s is NOT the winner\n", p.Address)
		}
	}

	return
}

func getBalances(players []*fdb.Player) {
	processBalancesCommand(players)
	if *verbose {
		for _, p := range players {
			fmt.Printf("[pickWinner] new players account: %v, balances: %s\n", p.Address, convertBalanceIntoReadableFormat(p.Balance))
		}
	}
}

func getSession() int64 {
	session := db.GetSession(false)
	if len(session) > 0 {
		if *verbose {
			fmt.Printf("Current Session ID: %v\n", session[0].ID)
		}
		return session[0].ID
	}
	fmt.Printf("[getSession] ERROR: get No Session\n")
	return 0
}

func getAllPlayer() []*fdb.Player {
	//Get all player from DB
	allPlayers = db.GetPlayers("", "", nil)
	if *verbose {
		for i, p := range allPlayers {
			fmt.Printf("[getAllPlayer] %v => %v\n", i, p)
		}
	}

	return allPlayers
}

func getPlayer() []*fdb.Player {
	//Get a list of all current players
	players, err := restclient.GetPlayer(leader.IP, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
		return nil
	}

	if *verbose {
		fmt.Printf("[getPlayer] REST: %v\n", players)
	}

	currentPlayers := fdb.NewPlayer(players)

	if *verbose && currentPlayers != nil {
		for i, p := range currentPlayers {
			fmt.Printf("[getPlayer:%v] account: %v, balances: %v/%v\n", i, p.Address, convertBalanceIntoReadableFormat(p.Balance), p.Balance)
		}
	}
	return currentPlayers
}

func notifyWinner() {
	// TODO: send email to winner
	return
}

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
	db, err = fdb.NewFdb(*key, *project)

	if err != nil || db == nil {
		log.Fatalf("Failed to create Fdb client: %v", err)
		os.Exit(1)
	}

	// Close FDB when done.
	defer db.CloseFdb()
	leader = readProfile(*profile)

	switch *action {
	case "reg":
		sendRegEmail()
	case "winner":
		pickWinner()
	case "player":
		getPlayer()
	case "players":
		getAllPlayer()
	case "notify":
		notifyWinner()
	case "balances":
		allPlayers := getAllPlayer()
		getBalances(allPlayers)
	default:
		fmt.Printf("Wrong action: %v", action)
	}

	os.Exit(0)
}
