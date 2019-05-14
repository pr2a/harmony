package utils

// this module in utils handles the ini file read/write
import (
	"fmt"
	"strings"

	"github.com/harmony-one/demo-apps/backend/p2p"
	ini "gopkg.in/ini.v1"
)

// BackendProfile contains a section and key value pair map
type BackendProfile struct {
	Profile   string
	Bootnodes []string
	Shards    int
	RPCServer [][]p2p.Peer
	RPCLeader []p2p.Peer
}

func parseNodes(s string) (nodes []p2p.Peer) {
	for _, node := range strings.Fields(s) {
		v := strings.Split(node, ":")
		host := v[0]
        port := v[1]
        nodes = append(nodes, p2p.Peer{IP: host, Port: port})
	}
	return
}

// ReadBackendProfile reads an ini file and return BackendProfile
func ReadBackendProfile(fn string, profile string) (*BackendProfile, error) {
	cfg, err := ini.ShadowLoad(fn)
	if err != nil {
		return nil, err
	}
	config := new(BackendProfile)
	config.Profile = profile

	// get the profile section
	sec, err := cfg.GetSection(profile)
	if err != nil {
		return nil, err
	}

	if sec.HasKey("bootnode") {
		config.Bootnodes = sec.Key("bootnode").ValueWithShadows()
	} else {
		return nil, fmt.Errorf("can't find bootnode key")
	}

	if sec.HasKey("shards") {
		config.Shards = sec.Key("shards").MustInt()
		config.RPCServer = make([][]p2p.Peer, config.Shards)
		config.RPCLeader = make([]p2p.Peer, config.Shards)
	} else {
		return nil, fmt.Errorf("can't find shards key")
	}

	for i := 0; i < config.Shards; i++ {
		rpcSec, err := cfg.GetSection(fmt.Sprintf("%s.shard%v.rpc", profile, i))
		if err != nil {
			return nil, err
		}
		for _, n := range rpcSec.Key("rpc").ValueWithShadows() {
			config.RPCServer[i] = append(config.RPCServer[i], parseNodes(n)...)
		}
		var leaders []p2p.Peer
		if len(config.RPCServer[i]) == 0 {
			return nil, fmt.Errorf("shard %v has no RPC servers configured", i)
		}
		for _, l := range rpcSec.Key("leader").ValueWithShadows() {
			leaders = append(leaders, parseNodes(l)...)
		}
		switch len(leaders) {
		case 0:
            // Backward compatibility - use first rpc node listed as the leader
            config.RPCLeader[i] = config.RPCServer[i][0]
		case 1:
			config.RPCLeader[i] = leaders[0]
		default:
			return nil, fmt.Errorf("shard %v has multiple leaders configured", i)
		}
	}

	return config, nil

}
