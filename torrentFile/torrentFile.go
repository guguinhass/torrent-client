import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"
	"github.com/guguinhass/torrent-client/p2p"
	"github.com/jackpal/bencode-go"
)

const Port uint16 = 6881

type TorrentFile struct {
	Announce string
	InfoHash [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length int
	Name string
}
type bencodeInfo struct {
	Pieces string `bencode:"pieces"`
	PiecesLength int `bencode:"piece lenght"`
	Length int `bencode:"length"`
	Name string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string `bencode:"announce"`
	Info bencodeInfo `bencode:"info"`
}
func (t *TorrentFile) DownloadToFile (path string) error{
	var peerID [20]byte
	_, err:= rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := t.requestPeers(peerID, Port)
	if err != nil {
		return err
	}

	torrent ;= p2p.Torrent{
		Peers:	peers,
		PeerID: peerID,
		InfoHash: t.InfoHash,
		PiecesHashes: t.PieceHashes,
		PieceLenght: t.PieceLenght,
		Length: t.Length,
		Name: t.Name,
	}

	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile := os.Create(path)
	if err != nil{
		return err
	}

	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func Open(r io.Reader) (*bencodeTorrent, error){
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}

func (i *bencodeInfo) splitPiecesHashes() ([][20]byte, error){
	hashlen := 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)
	
	for i:= 0; i< numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (i *bencodeInfo) hash() ([20]byte, error){
	var buf bytes.Buffer
	err := bencode.Marshall(&buf, *i)
	if err != nill {
		return [20]byte{}, err
	}
	h:= sha1.Sum(buf.Bytes())
	return h, nil
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPiecesHashes()
	if err!= nil {
		return TorrentFil{}, err
	}
	t := TorrentFile{
		Announce: bto.Announce,
		InfoHash: infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Lenght: bto.Info.Length,
		Name: bto.Info.Name,
	}
	return t, nil
}






