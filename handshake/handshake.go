package handshake

import (
	"errors"
	"io"
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

// handshakes are of fixed length, could've passed a buffer directly but for consistency
func Deserialize(conn io.Reader) (*handshake, error) {
	buf := make([]byte, 1+19+8+20+20)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, errors.New("Couldn't read handshake off open tcp connection")
	}
	pstrlen := int(buf[0])
	if int(buf[0]) != len("BitTorrent protocol") {
		return nil, errors.New("invalid pstr length")
	}

	pstr := string(buf[1 : 1+pstrlen])
	if pstr != "BitTorrent protocol" {
		return nil, errors.New("invalid pstr")
	}

	var infoHash [20]byte
	var peerId [20]byte
	copy(infoHash[:], buf[1+19+8:])
	copy(peerId[:], buf[1+19+8+20:])
	return New(infoHash, peerId), nil
}

func (h handshake) Validate(infoHash [20]byte) error {
	if h.infoHash != infoHash {
		return errors.New("infohash does not match")
	}
	return nil
}
