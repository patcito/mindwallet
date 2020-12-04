package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const DefaultAddress = "127.0.0.1:18081"

// Node represents a bitmonerod daemon.
type Node struct {
	client http.Client
	addr   string
}

// NewNode returns a new Node for a given address:port.
func NewNode(addr string) *Node {
	return &Node{addr: addr}
}

func (n *Node) rpcSend(path string, msg interface{}) (io.ReadSeeker, error) {
	var resp *http.Response
	var err error

	if msg != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(msg)
		if err == nil {
			resp, err = n.client.Post("http://"+n.addr+path, "application/json", &buf)
		}
	} else {
		resp, err = n.client.Post("http://"+n.addr+path, "application/json", nil)
	}
	if err != nil {
		return nil, err
	}

	check := new(struct {
		Status string `json:"status"`
	})

	rs, ok := resp.Body.(io.ReadSeeker)
	if !ok {
		b := make([]byte, resp.ContentLength)
		rs = bytes.NewReader(b)
		resp.Body.Read(b)
	}

	dec := json.NewDecoder(rs)
	err = dec.Decode(check)
	if err == nil {
		rs.Seek(0, 0)
		if check.Status != "OK" {
			err = errors.New(check.Status)
		}
	}
	return rs, err
}

// BlockchainHeight returns the height of the nodes blockchain.
func (n *Node) BlockchainHeight() (height int, err error) {
	resp, err := n.rpcSend("/getheight", nil)
	if err != nil {
		return 0, err
	}

	reply := &struct {
		Height *int `json:"height"`
	}{&height}
	dec := json.NewDecoder(resp)
	err = dec.Decode(&reply)
	return
}

//func (n *Node) BlocksFast(blockIds []string) error { }

/*
func (n *Node) Transactions(txsHashes []string) error {
	req := &struct {
		T []string `json:"txs_hashes"`
	}{txsHashes}

	_, err := n.rpcSend("/transactions", nil)
	return err
}
*/

// StartMining instructs the node to begin mining to Monero address using the given number of threads.
func (n *Node) StartMining(address string, threads int) error {
	req := &struct {
		A string `json:"miner_address"`
		T int    `json:"threads_count"`
	}{address, threads}

	_, err := n.rpcSend("/start_mining", req)
	return err
}

type NodeInfo struct {
	Height              int `json:"height"`
	Difficulty          int `json:"difficulty"`
	TxCount             int `json:"tx_count"`
	TxPoolSize          int `json:"tx_pool_size"`
	AltBlocksCount      int `json:"alt_blocks_count"`
	OutgoingConnections int `json:"outgoing_connections_count"`
	IncomingConnections int `json:"incoming_connections_count"`
	WhitePeerlistSize   int `json:"white_peerlist_size"`
	GreyPeerlistSize    int `json:"grey_peerlist_size"`
}

func (n *Node) Info() (ni *NodeInfo, err error) {
	ni = new(NodeInfo)
	r, err := n.rpcSend("/getinfo", nil)
	if err == nil {
		dec := json.NewDecoder(r)
		err = dec.Decode(ni)
	}
	return
}

// StopMining instructs the node to cease mining.
func (n *Node) StopMining() error {
	_, err := n.rpcSend("/stop_mining", nil)
	return err
}

func (n *Node) BlockCount() (count int, err error) {
	v := &struct {
		C *int `json:"count"`
	}{&count}

	_, err = n.rpcSend("/getblockcount", v)
	return
}
