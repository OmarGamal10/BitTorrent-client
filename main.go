package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/omargamal10/BitTorrent-client/handshake"
	"github.com/omargamal10/BitTorrent-client/torrentfile"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	path := dir + "/pybenc/debian-12.10.0-amd64-netinst.iso.torrent"
	torrent, err := torrentfile.LoadTorrent(path)
	if err != nil {
		fmt.Println("Error reading torrent file:", err)
		return
	}
	peers, err := torrent.GetTrackerResponse([20]byte{ // just a random peer ID for now
		'O', 'M', 'A', 'R', 'G', 'A', 'M', 'A', 'L', '1',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}, 6881)
	if err != nil {
		fmt.Println("Error getting tracker response:", err)
		return
	}
	// will try to send a handshake for first ip
	hs := handshake.New(torrent.InfoHash, [20]byte{
		'O', 'M', 'A', 'R', 'G', 'A', 'M', 'A', 'L', '1',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	})
	conn, err := net.DialTimeout("tcp", peers[0].IP.String()+":"+strconv.Itoa(int(peers[0].Port)), time.Millisecond*5000)
	if err != nil {
		fmt.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write(hs.Serialize())
	if err != nil {
		fmt.Println("Error sending handshake:", err)
		return
	}
	fmt.Println("Handshake sent to", peers[0].IP.String()+":"+strconv.Itoa(int(peers[0].Port)))
	respBuf := make([]byte, 68)
	_, err = conn.Read(respBuf)
	if err != nil {
		fmt.Println("Error reading handshake response:", err)
		return
	}
	handshakeResp, err := handshake.Deserialize(respBuf)
	if err != nil {
		fmt.Println("Error deserializing handshake response:", err)
		return
	}
	fmt.Printf("Received handshake response: %+v\n", handshakeResp)
	// check if the peerId is the same one sent by client
}

// we have the handshake we should expect messages now
