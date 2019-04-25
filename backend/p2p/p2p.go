package p2p

import (
	"fmt"
	"net"
)

// Peer is the object for a p2p peer (node)
type Peer struct {
	IP              string         // IP address of the peer
	Port            string         // Port number of the peer
}

func (p Peer) String() string {
	return fmt.Sprintf("%s", net.JoinHostPort(p.IP, p.Port))
}
