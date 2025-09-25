package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/omargamal10/BitTorrent-client/bitfield"
	"github.com/omargamal10/BitTorrent-client/handshake"
	"github.com/omargamal10/BitTorrent-client/message"
	"github.com/omargamal10/BitTorrent-client/peers"
	"github.com/omargamal10/BitTorrent-client/torrentfile"
)

// represents a connection with a peer from handshake till sever
// again, this is simplified as I have NO WAY of testing seeding, so I can't implement it
type Connection struct {
	conn       io.ReadWriter
	bf         bitfield.Bitfield
	peer       peers.Peer
	choked     bool
	interested bool
}

func ConnectToPeers(peerList []peers.Peer, tf torrentfile.TorrentFile) []Connection {
	connections := make([]Connection, 0)
	for _, p := range peerList {
		connection, err := connectToPeer(p, tf)
		if err != nil {
			fmt.Println(err)
		}

		if connection != nil {
			connections = append(connections, *connection)
		}
	}
	js, err := json.Marshal(connections)
	fmt.Printf("%s\n", js)
	if err != nil {
		fmt.Println("Error marshaling connections:", err)
	}
	return connections
}

func connectToPeer(p peers.Peer, tf torrentfile.TorrentFile) (*Connection, error) {
	conn, err := net.DialTimeout("tcp", p.SocketAddress(), time.Millisecond*2000)
	if err != nil {
		return nil, fmt.Errorf("could not connect to peer: %s", p.SocketAddress())
	}
	hs := handshake.New(tf.InfoHash, [20]byte{
		'O', 'M', 'A', 'R', 'G', 'A', 'M', 'A', 'L', '1',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	})

	conn.Write(hs.Serialize())
	resp, err := handshake.Deserialize(conn)

	if err != nil {
		return nil, fmt.Errorf("Could not get a handshake back from peer: %s", p.SocketAddress())
	}
	if resp.Validate(tf.InfoHash) != nil {
		return nil, errors.New("Poisoned torrent, info hash does not match")
	}
	fmt.Println("connectted to peer: ", p.SocketAddress())

	return newPeer(conn, p)
}

func newPeer(conn io.ReadWriter, p peers.Peer) (*Connection, error) {
	msg, err := message.Deserialize(conn)
	fmt.Println("serialized message from peer: ", p.SocketAddress())
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//expecting a bitfield
	if msg.Id != message.Bitfield {
		fmt.Println("not a bitfield")
		return nil, errors.New("first message must be a bitfield")
	}
	return &Connection{
		conn:       conn,
		bf:         msg.Payload,
		peer:       p,
		choked:     true,
		interested: false,
	}, nil
}
