package main

import (
	"fmt"
	"os"

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
	_, err = torrent.GetTrackerResponse([20]byte{ // just a random peer ID for now
		'O', 'M', 'A', 'R', 'G', 'A', 'M', 'A', 'L', '1',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}, 6881)
	if err != nil {
		fmt.Println("Error getting tracker response:", err)
		return
	}
}

// we have the peers, now the handshake so we start accepting and sending messages
