package peers

import (
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

// just for formatting
func (p Peer) socketAddress() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

// Takes in the compact peer list (string encoded) and returns a list of peers
func UnMarshal(peerBuf []byte) ([]Peer, error) {
	if len(peerBuf)%6 != 0 {
		return nil, fmt.Errorf("corrupted peer list")
	}
	numPeers := len(peerBuf) / 6
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		ip := net.IPv4(peerBuf[i*6], peerBuf[i*6+1], peerBuf[i*6+2], peerBuf[i*6+3]) // idk if it even supports ipv6
		port := uint16(peerBuf[i*6+4])<<8 | uint16(peerBuf[i*6+5])                   // big endian
		peers[i] = Peer{
			IP:   ip,
			Port: port,
		}
		// fmt.Println(peers[i].socketAddress())
	}
	return peers, nil
}
