//Package restclient is the rest client module talking to the RPC server of Harmony blockchain
package restclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/harmony-one/demo-apps/backend/p2p"
)

const (
	adminKey = "27978f895b11d9c737e1ab1623fde722c04b4f9ccb4ab776bf15932cc72d7c66"
)

//Winner is the structure returned from /winner rest call
type Winner struct {
	Players  string `json:players`
	Balances string `json:balances`
	Success  bool   `json:success`
}

//Player is the structure returned from /result rest call
//All the players in the lottery
type Player struct {
	Players  []string `json:players`
	Balances []string `json:balances`
	Success  bool     `json:success`
}

var (
	leaders = make([]p2p.Peer, 0)
)

//SetLeaders set the leader ip and port
func SetLeaders(l []p2p.Peer) {
	for _, p := range l {
		leaders = append(leaders, p)
	}
}

//GetLeaders return the list of existing leaders
func GetLeaders() []p2p.Peer {
	return leaders
}

// PickALeader return a random leader from the leader list
func PickALeader() p2p.Peer {
	return leaders[rand.Intn(len(leaders))]
}

//GetWinner return the result of a rest api call
func GetWinner(ip, port string) (*Winner, error) {
	url := fmt.Sprintf("http://%s:%s/winner", ip, port)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[GetWinner] GET winner error: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[GetWinner] can't get winner data")
	}

	contents, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("[GetWinner] failed to read response: %v", err)
	}

	var winner Winner
	err = json.Unmarshal(contents, &winner)

	if err != nil {
		return nil, fmt.Errorf("[GetWinner] failed to unmarshal winner response: %v", err)
	}

	if !winner.Success {
		return &winner, fmt.Errorf("[GetWinner] Failed on blockchain")
	}

	return &winner, nil
}

//GetPlayer return the result of a rest api call
func GetPlayer(ip, port string) (*Player, error) {
	url := fmt.Sprintf("http://%s:%s/result?key=%s", ip, port, adminKey)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[GetPlayer] GET result error: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[GetPlayer] can't get result data")
	}

	contents, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("[GetPlayer] failed to read response: %v", err)
	}

	var player Player
	err = json.Unmarshal(contents, &player)

	if err != nil {
		return nil, fmt.Errorf("[GetPlayer] failed to unmarshal result response: %v", err)
	}

	if !player.Success {
		return &player, fmt.Errorf("[GetPlayer] Failed on blockchain")
	}

	return &player, nil
}

// FundMe call /fundme rest call on leader
func FundMe(leader p2p.Peer, account string, wg sync.WaitGroup) error {
	defer wg.Done()
	url := fmt.Sprintf("http://%s:%s/fundme?key=0x%s", leader.IP, leader.Port, account)
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("[FundMe] GET result error: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("[FundMe] can't get result data")
	}
	contents, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("[FundMe] failed to read response: %v", err)
	}

	var player Player
	err = json.Unmarshal(contents, &player)

	if err != nil {
		return fmt.Errorf("[FundMe] failed to unmarshal result response: %v", err)
	}

	if !player.Success {
		return fmt.Errorf("[FundMe] Failed on blockchain")
	}

	return nil
}

// GetBalance call /balance rest call on leader
func GetBalance(account string) (uint64, error) {
	return 0, nil
}

// EnterPuzzle calls /enter rest call to enter the game and return the current level
func EnterPuzzle(account string, amount uint64) (uint, error) {
	return 0, nil
}

// GetRewards call /finish rest call to get rewards
func GetRewards(account string, level int) (uint64, error) {
	return 0, nil
}
