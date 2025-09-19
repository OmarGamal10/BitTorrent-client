package torrentfile

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
	"github.com/omargamal10/BitTorrent-client/peers"
)

// bencoded response from the tracker
// for now I will focus on http trackers
type TrackerResponse struct {
	Interval int    `bencode:"interval"` //interval in seconds to wait before re-requesting
	Peers    string `bencode:"peers"`    //encoded as a list of 6 byte strings, 4 bytes IP, 2 bytes port Big Endian
}

func (t TorrentFile) GetTrackerResponse(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	baseUrl, err := url.Parse(t.Announce)
	if err != nil {
		return nil, err
	}

	params := baseUrl.Query()
	params.Set("info_hash", string(t.InfoHash[:]))
	params.Set("peer_id", string(peerID[:])) // random
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
	peers, err := peers.UnMarshal([]byte(trackerRes.Peers))
	if err != nil {
		return nil, err
	}
	return peers, nil
}
