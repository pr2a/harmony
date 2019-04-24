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

	"github.com/jinzhu/copier"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/harmony-one/demo-apps/backend/client"
	fdb "github.com/harmony-one/demo-apps/backend/db"
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
	profile    = flag.String("profile", "default", "name of the profile")
	collection = flag.String("collection", "players", "name of collection")
	key        = flag.String("key", "./keys/benchmark_account_key.json", "key filename")
	project    = flag.String("project", "lottery-demo-leo", "project ID of firebase")
	ip         = flag.String("server_ip", "34.222.210.98", "the IP address of the server")
	action     = flag.String("action", "reg", "action of the program. Valid (reg, winner)")

	versionFlag = flag.Bool("version", false, "Output version info")

	db            *fdb.Fdb
	walletProfile *utils.WalletProfile
)

// FetchBalance fetches account balance of specified address from the Harmony network
func FetchBalance(address common.Address) map[uint32]AccountState {
	result := make(map[uint32]AccountState)
	for i := 0; i < walletProfile.Shards; i++ {
		balance := big.NewInt(0)
		var nonce uint64

		result[uint32(i)] = AccountState{balance, 0}

		for retry := 0; retry < rpcRetry; retry++ {
			server := walletProfile.RPCServer[i][rand.Intn(len(walletProfile.RPCServer[i]))]
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

func processBalancesCommand(addresses []string, balances []uint64) {
	for i, address := range addresses {
		addr := common.HexToAddress(address)
		fmt.Printf("Address: %s\n", addr)
		// assuming number of shard is 1
		for shardID, balanceNonce := range FetchBalance(addr) {
			fmt.Printf("    Balance in Shard %d:  %s, nonce: %v \n", shardID, convertBalanceIntoReadableFormat(balanceNonce.balance), balanceNonce.nonce)
			balances[i] = balanceNonce.balance.Uint64()
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
	var currentPlayers restclient.Player
	copier.Copy(&currentPlayers, player)

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

	processBalancesCommand(currentPlayers.Players, currentPlayers.Balances)

	for i, p := range currentPlayers.Players {
		fmt.Printf("account: %v, balances: %v\n", p, currentPlayers.Balances[i])
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
