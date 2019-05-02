//Package restclient is the rest client module talking to the RPC server of Harmony blockchain
package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "math/rand" // use it later
	"net/http"
	"time"

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

//PlayResp is the structure returned from /play rest call
type PlayResp struct {
	Players  []string `json:players`
	Balances []string `json:balances`
	Success  bool     `json:success`
}

// RPCMsg is a structure to exchange info between RPC client
type RPCMsg struct {
	Err  error
	Done bool
}

var (
	leaders    = make([]p2p.Peer, 0)
	rpcTimeout = 5 * time.Second
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
	//	return leaders[rand.Intn(len(leaders))]
	// FIXME: leo testing only
	return leaders[0]
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

// getClient is the generic get client for rest call
func getClient(url string, prefix string, result interface{}) error {
	fmt.Printf("getClient [%v]\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("[%v] can't new request", prefix)
	}

	ctx, cancel := context.WithTimeout(req.Context(), rpcTimeout)
	client := http.DefaultClient
	defer cancel()

	req = req.WithContext(ctx)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[%v] Do error: %s", prefix, err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("[%v] Status is not Ok", prefix)
	}

	contents, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("[%v] ReadAll failed: %v", prefix, err)
	}

	err = json.Unmarshal(contents, result)

	if err != nil {
		return fmt.Errorf("[%v] Unmarshal failed: %v", prefix, err)
	}

	return nil
}

// FundMe call /fundme rest call on leader
func FundMe(leader p2p.Peer, account string, done chan (RPCMsg)) {
	var err error
	var player Player
	var contents []byte
	var response *http.Response
	var ctx context.Context
	var cancel context.CancelFunc

	url := fmt.Sprintf("http://%s:%s/fundme?key=0x%s", leader.IP, leader.Port, account)
	fmt.Printf("FundMe: %v\n", url)

	client := http.DefaultClient

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("[FundMe] can't new request")
		goto DONE
	}

	ctx, cancel = context.WithTimeout(req.Context(), rpcTimeout)
	defer cancel()

	req = req.WithContext(ctx)

	response, err = client.Do(req)
	if err != nil {
		err = fmt.Errorf("[FundMe] GET Do error: %s", err)
		goto DONE
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("[FundMe] Status is not Ok")
		goto DONE
	}
	contents, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		err = fmt.Errorf("[FundMe] failed to read response: %v", err)
		goto DONE
	}

	err = json.Unmarshal(contents, &player)

	if err != nil {
		err = fmt.Errorf("[FundMe] failed to unmarshal result response: %v", err)
		goto DONE
	}

	if !player.Success {
		err = fmt.Errorf("[FundMe] Failed on blockchain")
		goto DONE
	}

DONE:
	done <- RPCMsg{
		Err:  err,
		Done: true,
	}
	return
}

// GetBalance call /balance rest call on leader
func GetBalance(leader p2p.Peer, account string, done chan (RPCMsg)) (uint64, error) {
	return 0, nil
}

// PlayGame calls /play rest call to enter the game and return the current level
func PlayGame(leader p2p.Peer, account string, amount string, done chan (RPCMsg)) {
	// FIXME: this is based on old RPC API
	//	url := fmt.Sprintf("http://%s:%s/play?key=0x%s&amount=%s", leader.IP, leader.Port, account, amount)
	url := fmt.Sprintf("http://%s:%s/play?key=%s", leader.IP, leader.Port, account)

	var play = new(PlayResp)

	err := getClient(url, "/play", play)

	done <- RPCMsg{
		Err:  err,
		Done: true,
	}
}

// GetRewards call /finish rest call to get rewards
func GetRewards(account string, height int64) (uint64, error) {
	return 0, nil
}
