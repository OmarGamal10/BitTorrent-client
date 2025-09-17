package main

import (
	"encoding/json"
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
	pretty, _ := json.MarshalIndent(torrent, "", "  ")
	fmt.Println(string(pretty))
}
