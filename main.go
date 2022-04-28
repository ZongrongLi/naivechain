package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1465154705,
	Data:         []byte("my genesis block!!"),
	Hash:         "816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7",
}

type Block struct {
	Index        int64  `json:"index"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previousHash"`
	Timestamp    int64  `json:"timestamp"`
	Data         []byte `json:"data"`
}

type PeerMsg struct {
	Data []byte `json:"data"`
}

var (
	localhostAddr = "ws://localhost:6001"
	blockchain    []*Block
	peersAddr     = []string{
		"ws://localhost:6002",
	}
	sockets []*websocket.Conn
)

func initConnection(peerAddr string) *websocket.Conn {
	ws, err := websocket.Dial(peerAddr, "", "")
	if err != nil {
		panic("fatal")
	}
	return ws
}

func broadcast(block *Block) {
	for _, peerWs := range sockets {
		t, err := json.Marshal(block)
		if err != nil {
			fmt.Println(err)
		}
		peerWs.Write(t)
	}
}

func generateandBroadCast(data []byte) {
	b := genesisBlock
	b.Data = data
	b.PreviousHash = blockchain[len(blockchain)-1].Hash
	blockchain = append(blockchain, b)
	broadcast(b)
}
func handlerBlock(msg *PeerMsg) {
	b := Block{}
	json.Unmarshal(msg.Data, &b)
	blockchain = append(blockchain, &b)
}

func initPeerConnection(ws *websocket.Conn) {
	go func() {
		for {
			msg := PeerMsg{}
			websocket.Message.Receive(ws, &msg)
			fmt.Println("s", msg)
			handlerBlock(&msg)
		}
	}()
}

// 产生区块发布到每个节点 维护各个节点的block 同步
func main() {
	fmt.Println("Hello, playground")
	ws := initConnection(localhostAddr)
	sockets = append(sockets, ws)
	initPeerConnection(ws)

	http.ListenAndServe(localhostAddr, nil)
}
