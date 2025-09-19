package torrentfile

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type Peer struct {
	IP   net.IP
	Port uint16 //normally listen on 6881, give up after 6889
}

// bencoded response from the tracker
// for now I will focus on http trackers
type TrackerResponse struct {
	Interval int    `bencode:"interval"` //interval in seconds to wait before re-requesting
	Peers    string `bencode:"peers"`    //encoded as a list of 6 byte strings, 4 bytes IP, 2 bytes port Big Endian
}

// socket address of a peer
func (p Peer) SocketAddress() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func unMarshal(peerBuf []byte) ([]Peer, error) {
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
	}
	return peers, nil
}

func (t TorrentFile) GetTrackerResponse(peerID [20]byte, port uint16) ([]Peer, error) {
	baseUrl, err := url.Parse(t.Announce)
	if err != nil {
		return nil, err
	}

	params := baseUrl.Query()
	params.Set("info_hash", string(t.InfoHash[:]))
	params.Set("peer_id", string(peerID[:]))
	params.Set("port", strconv.Itoa(int(port)))
	params.Set("uploaded", "0")
	params.Set("downloaded", "0")
	params.Set("left", strconv.Itoa(t.Length))
	params.Set("compact", "1")

	baseUrl.RawQuery = params.Encode()
	fmt.Println("Tracker URL:", baseUrl.String())
	res, err := http.Get(baseUrl.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // body wraps the underlying connection, so close it

	trackerRes := TrackerResponse{}
	err = bencode.Unmarshal(res.Body, &trackerRes)
	if err != nil {
		return nil, err
	}
	peers, err := unMarshal([]byte(trackerRes.Peers))
	if err != nil {
		return nil, err
	}
	for _, p := range peers {
		fmt.Println(p.SocketAddress())
	}
	return peers, nil
}
