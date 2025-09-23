package handshake

import (
	"errors"
	"fmt"
)

// 8 reserved btyes after pstrlen
type handshake struct {
	pstr     string // always "BitTorrent protocol" unquoted
	infoHash [20]byte
	peerID   [20]byte
}

// to enfore the Pstr
func New(infoHash [20]byte, peerId [20]byte) *handshake {
	pstr := "BitTorrent protocol"
	return &handshake{pstr: pstr, infoHash: infoHash, peerID: peerId}
}

func (h handshake) Serialize() []byte {
	buf := make([]byte, (1 + len(h.pstr) + 8 + 20 + 20))
	buf[0] = byte(len(h.pstr))
	copy(buf[1:], []byte(h.pstr))
	copy(buf[1+len(h.pstr):], make([]byte, 8)) // 8 reserved bytes
	copy(buf[1+len(h.pstr)+8:], h.infoHash[:])
	copy(buf[1+len(h.pstr)+8+20:], h.peerID[:])
	return buf
}

func Deserialize(buf []byte) (*handshake, error) {
	pstrl := int(buf[0])
	if pstrl != len("BitTorrent protocol") {
		fmt.Println("pstr is ", pstrl)
		return nil, errors.New("invalid pstr length")
	}
	pstr := string(buf[1 : 1+pstrl])
	if pstr != "BitTorrent protocol" {
		return nil, errors.New("invalid pstr")
	}
	var infoHash [20]byte
	var peerId [20]byte
	copy(infoHash[:], buf[1+19+8:])
	copy(peerId[:], buf[1+19+8+20:])
	return New(infoHash, peerId), nil
}
