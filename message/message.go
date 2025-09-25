package message

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type messageId int

const (
	Choke messageId = iota
	Unchoke
	Interested
	NotInterested
	Have
	Bitfield
	Request
	Piece
	Cancel
)

type Message struct {
	Id      messageId
	Payload []byte
}

// these are for payload formatting
func (m *Message) PayloadHave(idx int) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(m.Payload, uint32(idx))
	m.Payload = buf
}

// as of now VPN won't let me upload, I can only be a leecher, this will not be used
// func (m *Message) PayloadPiece()

// leave bitfield for now

// should i make length always 16kb? mmm
// begin and length are offsets within A PIECE
func (m *Message) PayloadRequest(idx int, begin int, length int) {
	populateCancelAndRequest(m, idx, begin, length)
}

// idk what this is for yet
func (m *Message) PayloadCancel(idx int, begin int, length int) {
	populateCancelAndRequest(m, idx, begin, length)
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	if m.Id == Choke || m.Id == Unchoke || m.Id == Interested || m.Id == NotInterested {
		return []byte{0, 0, 0, 1, byte(m.Id)}
	}
	length := len(m.Payload) + 1
	buf := make([]byte, 4+1+length)
	binary.BigEndian.PutUint32(buf, uint32(length))
	buf[5] = byte(m.Id)
	copy(buf[5:], m.Payload)
	return buf
}

func Deserialize(conn io.Reader) (*Message, error) {
	var length uint32
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not read message length")
		// will probably cut the connection or something then
	}
	buf := make([]byte, length)
	io.ReadFull(conn, buf)
	fmt.Printf("message is: %v\n", buf)
	return &Message{
		Id:      messageId(buf[0]), // already read the length prefix
		Payload: buf[1:],
	}, nil
}

func populateCancelAndRequest(m *Message, idx int, begin int, length int) {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:4], uint32(idx))
	binary.BigEndian.PutUint32(buf[4:8], uint32(begin))
	binary.BigEndian.PutUint32(buf[8:12], uint32(length))
	m.Payload = buf
}
