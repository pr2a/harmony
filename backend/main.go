package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/appengine/urlfetch"
	"gopkg.in/sendgrid/sendgrid-go.v2"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/harmony-one/demo-apps/backend/client"
	"github.com/harmony-one/demo-apps/backend/db"
	"github.com/harmony-one/demo-apps/backend/p2p"
	clientService "github.com/harmony-one/demo-apps/backend/service"
	"github.com/harmony-one/demo-apps/backend/utils"
	"google.golang.org/appengine"
	app_log "google.golang.org/appengine/log"
	"net/http"
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

// replaceInEmail keep track of replacement string in email text/html
type replaceInEmail struct {
	oldStr string
	newStr string
}

const (
	rpcRetry             = 3
	defaultConfigFile    = "./backend/.hmy/backend.ini"
	defaultProfile       = "default"
	port                 = "30000"
	winningEmailBody     = "Great work! You did it\nYou slayed dragons and conquered new land and you are the winner of the current Harmony lottery session.\nTo claim your Harmony tokens, join in our [discord, telegram, twitter]\nBut can you win again at http://lottery.harmony.one?"
	losingEmailBody      = "Try again at http://lottery.harmony.one!\nYou can do it!\nChak De Phatte!"
	winningEmailBodyHTML = ""
	losingEmailBodyHTML  = ""
	winnerEmailTitle     = "Feeling lucky? And the winner is..."
	losingEmailTitle     = "Feeling lucky? And the winner is..."

	winningEmailFolder = "email/win"
	losingEmailFolder  = "email/lose"
)

func printVersion(me string) {
	fmt.Fprintf(os.Stderr, "Harmony (C) 2019. %v, version %v-%v (%v %v)\n", path.Base(me), version, commit, builtBy, builtAt)
	os.Exit(0)
}

var (
	profile    = flag.String("profile", defaultProfile, "name of the profile")
	collection = flag.String("collection", "players", "name of collection")
	key        = flag.String("key", "./backend/keys/benchmark_account_key.json", "key filename")
	project    = flag.String("project", "benchmark-209420", "project ID of firebase")
	action     = flag.String("action", "player", "action of the program. Valid (player, reg, winner, notify, players, balances)")
	verbose    = flag.Bool("verbose", true, "verbose log print at every step")
	local      = flag.Bool("local", false, "Run locally")

	versionFlag = flag.Bool("version", false, "Output version info")

	db             *fdb.Fdb
	backendProfile *utils.BackendProfile
	leader         p2p.Peer
	mux            sync.Mutex
	allPlayers     []*fdb.Player
	winningInfo    *fdb.Winner
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

func processBalancesCommand(players []*fdb.Player, r *http.Request) {
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
				fmt.Printf("Balance in Shard %d:  %s/%v\n", shardID, convertBalanceIntoReadableFormat(balanceNonce.balance), balanceNonce.balance)
			}
			player.Balance.Set(balanceNonce.balance)
		}
	}
}

func convertBalanceIntoReadableFormat(balance *big.Int) string {
	curBalance := big.NewInt(0)
	curBalance.SetBytes(balance.Bytes())

	curBalance = curBalance.Div(curBalance, big.NewInt(params.GWei))
	strBalance := fmt.Sprintf("%d", curBalance.Uint64())

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

func findEmail(address string) string {
	for _, p := range allPlayers {
		if strings.ToLower(address) == strings.ToLower(p.Address) {
			return p.Email
		}
	}
	return ""
}

func sendRegEmail() {
	players := db.GetPlayers("email_key", "==", false)

	for i, p := range players {
		fmt.Printf("[sendRegEmail] %v => %v\n", i, p)
	}

	//TODO: send an email to user for the account registration
}

func sendWinningEmail(email string) {
	log.Printf("Sending email to winner: %v", email)
}

func pickWinner(r *http.Request) ([]string, []string) {
	ctx := appengine.NewContext(r)
	currentPlayers := getPlayer(r)
	existingPlayers := make([]*fdb.Player, 0)

	winningInfo = new(fdb.Winner)
	winningInfo.Session = getSession()
	go getAllPlayer(winningInfo.Session)

	for _, p := range currentPlayers {
		onePlayer := new(fdb.Player)
		onePlayer.Address = p.Address
		onePlayer.Balance = new(big.Int)
		existingPlayers = append(existingPlayers, onePlayer)
	}

	//Run the get winner smart contract
	winner, err := restclient.GetWinner(leader.IP, port)
	if err != nil {
		app_log.Criticalf(ctx, "GetWinner Error: %v", err)
	} else {
		app_log.Infof(ctx, "GetWinner: %v", winner)
	}

	// wait for the execution of smart contracts
	time.Sleep(15 * time.Second)

	processBalancesCommand(existingPlayers, r)

	winners := make([]string, 0)
	losers := make([]string, 0)

	for i, p := range existingPlayers {
		if p == nil {
			continue
		}
		email := findEmail(p.Address)
		// TODO: mark the winner explicitly in smart contract
		if p.Balance.Cmp(currentPlayers[i].Balance) > 0 {
			winners = append(winners, email)
			winningInfo.Address = p.Address
			z := p.Balance.Sub(p.Balance, currentPlayers[i].Balance)
			z = z.Div(z, big.NewInt(params.GWei))
			winningInfo.Amount = z.Int64()
		} else {
			losers = append(losers, email)
		}
	}
	sessionID := winningInfo.Session + 1
	// set current is_current to false
	db.UpdateSession()

	// add a new session id
	db.AddSession(sessionID)

	if *verbose {
		app_log.Infof(ctx, "Winner: %v\n", winners)
		app_log.Infof(ctx, "Loser: %v\n", losers)
		app_log.Infof(ctx, "Amount: %v\n", winningInfo.Amount)
		fmt.Printf("WINNER: %v\n", winners)
		fmt.Printf("LOSERS: %v\n", losers)
		fmt.Printf("Amount: %v\n", winningInfo.Amount)
	}

	db.AddWinner(winningInfo)

	return winners, losers
}

func getBalances(players []*fdb.Player) {
	processBalancesCommand(players, nil)
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

//getAllPlayer will get all players from DB
func getAllPlayer(sess int64) []*fdb.Player {
	mux.Lock()
	if sess > 0 {
		allPlayers = db.GetPlayers("session_id", "==", sess)
	} else {
		allPlayers = db.GetPlayers("", "", nil)
	}
	mux.Unlock()

	if *verbose {
		for i, p := range allPlayers {
			fmt.Printf("[getAllPlayer] %v => %v\n", i, p)
		}
	}

	return allPlayers
}

func getPlayer(r *http.Request) []*fdb.Player {
	var ctx context.Context

	if !*local {
		ctx = appengine.NewContext(r)
	}

	//Get a list of all current players
	players, err := restclient.GetPlayer(leader.IP, port)
	if err != nil {
		log.Fatalf("GetPlayer Error: %v", err)
		return nil
	}

	if *verbose {
		fmt.Printf("[getPlayer] REST: %v\n", players)
		if !*local {
			app_log.Infof(ctx, "[getPlayer]: %v\n", players)
		}
	}

	currentPlayers := fdb.NewPlayer(players)

	if *verbose && currentPlayers != nil {
		for i, p := range currentPlayers {
			fmt.Printf("[getPlayer:%v] account: %v, balances: %s/%v\n", i, p.Address, convertBalanceIntoReadableFormat(p.Balance), p.Balance)
			if !*local {
				app_log.Infof(ctx, "[getPlayer:%v] account: %v, balances: %s/%v\n", i, p.Address, convertBalanceIntoReadableFormat(p.Balance), p.Balance)
			}
		}
	}
	return currentPlayers
}

func notifyWinner(winnerEmails, nonWinnerEmails []string, r *http.Request) {
	noWinner := false
	if winningInfo == nil {
		winningInfo = new(fdb.Winner)
		noWinner = true
	}
	rep := []replaceInEmail{
		{
			"$SESSION$",
			fmt.Sprintf("%v", winningInfo.Session),
		},
		{
			"$AMOUNT$",
			fmt.Sprintf("%v", winningInfo.Amount/1000000000),
		},
	}

	winText, winHTML := getEmail(true, rep)
	loseText, loseHTML := getEmail(false, rep)

	if !noWinner {
		sendEmail(winnerEmails, winnerEmailTitle, winText, winHTML, r, false)
	}
	sendEmail(nonWinnerEmails, losingEmailTitle, loseText, loseHTML, r, true)
	return
}

func getEmail(win bool, r []replaceInEmail) (string, string) {
	folder := losingEmailFolder
	textStr := losingEmailBody
	htmlStr := losingEmailBodyHTML
	if win {
		folder = winningEmailFolder
		textStr = winningEmailBody
		htmlStr = winningEmailBodyHTML
	}

	text, err := ioutil.ReadFile(path.Join(folder, "email.txt"))
	if err != nil {
		fmt.Printf("read email text file error: %v", err)
	} else {
		textStr = string(text)
	}
	html, err := ioutil.ReadFile(path.Join(folder, "email.html"))
	if err != nil {
		fmt.Printf("read email html file error: %v", err)
	} else {
		htmlStr = string(html)
	}

	for _, s := range r {
		textStr = strings.Replace(textStr, s.oldStr, s.newStr, -1)
		htmlStr = strings.Replace(htmlStr, s.oldStr, s.newStr, -1)
	}

	return textStr, htmlStr
}

func sendEmail(recipients []string, title, body, htmlBody string, r *http.Request, batch bool) {
	ctx := appengine.NewContext(r)
	if len(recipients) == 0 {
		app_log.Infof(ctx, "Recipients list is empty")
		return
	}

	sg := sendgrid.NewSendGridClient("", "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")

	// Set http.Client to use the App Engine urlfetch client
	sg.Client = urlfetch.Client(ctx)

	if batch {
		// Batch for losers
		message := sendgrid.NewMail()
		message.AddTo("lottery@harmony.one")
		for _, email := range recipients {
			if email != "" {
				message.AddBcc(email)
			}
		}
		message.SetSubject(title)
		message.SetText(body)
		message.SetHTML(htmlBody)
		message.SetFrom("admin@harmony-lottery-app.appspotmail.com")

		app_log.Infof(ctx, "Sending email batch: %s", message)
		if err := sg.Send(message); err != nil {
			app_log.Errorf(ctx, "Couldn't send email", err)
		}
	} else {
		for _, email := range recipients {
			if email != "" {
				message := sendgrid.NewMail()
				message.AddTo(email)
				message.SetSubject(title)
				message.SetText(body)
				message.SetHTML(htmlBody)
				message.SetFrom("admin@harmony-lottery-app.appspotmail.com")

				app_log.Infof(ctx, "Sending email separately: %s", message)
				if err := sg.Send(message); err != nil {
					app_log.Errorf(ctx, "Couldn't send email", err)
				}
			}
		}
	}
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

	return backendProfile.RPCLeader[0]
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

	if *local {
		switch *action {
		case "reg":
			sendRegEmail()
		case "winner":
			pickWinner(nil)
		case "player":
			getPlayer(nil)
		case "players":
			getAllPlayer(0)
		case "notify":
			notifyWinner([]string{}, []string{}, nil)
		case "balances":
			allPlayers := getAllPlayer(0)
			getBalances(allPlayers)
		default:
			fmt.Printf("Wrong action: %v", action)
		}

		os.Exit(0)
	} else {
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/pickwinner", pickWinnerHandler)
		http.HandleFunc("/player", playerHandler)

		appengine.Main()
	}
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

// pickWinnerHandler responds to requests for the pick winner cron job.
func pickWinnerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/pickwinner" {
		http.NotFound(w, r)
		return
	}

	winners, losers := pickWinner(r)
	notifyWinner(winners, losers, r)
}

func playerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/player" {
		http.NotFound(w, r)
		return
	}
	getPlayer(r)
}
