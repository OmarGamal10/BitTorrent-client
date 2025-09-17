package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"io"
	"os"

	"github.com/jackpal/bencode-go"
)

// the application definition of a torrent file, comes from the serialization structs
type TorrentFile struct {
	Announce    string
	Name        string
	Length      int
	PieceLength int
	PieceHashes [][20]byte // SHA-1 hash is 20 bytes
	InfoHash    [20]byte   // SHA-1 hash of `info`, the fingerprint of the file to the tracker
}

// serialization structs
type bencodeTorrent struct {
	Announce string `bencode:"announce"`
	Info     info   `bencode:"info"`
}

type info struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

func (i info) hash() ([20]byte, error) {
	var buf bytes.Buffer // implements io.Writer
	err := bencode.Marshal(&buf, i)
	if err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func (i info) pieceHashes() ([][20]byte, error) {
	len := len(i.Pieces)
	if len%20 != 0 {
		return nil, errors.New("torrent is likely corrupted, invalid piece length")
	}
	numHashes := len / 20
	buf := []byte(i.Pieces)
	pieceHashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(pieceHashes[i][:], buf[i*20:(i+1)*20])
	}
	return pieceHashes, nil
}

func ReadTorrent(r io.Reader) (*TorrentFile, error) {

	bt := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bt)
	if err != nil {
		return nil, err
	}

	torrentFile, err := toTorrentFile(bt)
	if err != nil {
		return nil, err
	}
	return torrentFile, nil
}

func LoadTorrent(path string) (*TorrentFile, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ReadTorrent(r)
}

func toTorrentFile(bt bencodeTorrent) (*TorrentFile, error) {
	infoHash, err := bt.Info.hash()
	if err != nil {
		return nil, err
	}
	pieceHashes, err := bt.Info.pieceHashes()
	if err != nil {
		return nil, err
	}
	torrent := &TorrentFile{
		Announce:    bt.Announce,
		Name:        bt.Info.Name,
		Length:      bt.Info.Length,
		PieceLength: bt.Info.PieceLength,
		PieceHashes: pieceHashes,
		InfoHash:    infoHash,
	}
	return torrent, nil
}
