package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     Hash                    `json:"block_hash"`
	Balances map[common.Address]uint `json:"balances"`
}

type TxAddReq struct {
	From     string `json:"from"`
	FromPwd  string `json:"from_pwd"`
	To       string `json:"to"`
	Gas      uint   `json:"gas"`
	GasPrice uint   `json:"gasPrice"`
	Value    uint   `json:"value"`
	Data     string `json:"data"`
}

type TxAddRes struct {
	Success bool `json:"success"`
}

type StatusRes struct {
	Hash        Hash                `json:"block_hash"`
	Number      uint64              `json:"block_number"`
	KnownPeers  map[string]PeerNode `json:"peers_known"`
	PendingTXs  []SignedTx          `json:"pending_txs"`
	NodeVersion string              `json:"node_version"`
	Account     common.Address      `json:"account"`
}

type SyncRes struct {
	Blocks []Block `json:"blocks"`
}

type AddPeerRes struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *State) {
	enableCors(&w)

	writeRes(w, BalancesRes{Hash(state.LatestBlockHash()), state.Balances})
}

func txAddHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	req := TxAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	from := NewAccount(req.From)

	if from.String() == common.HexToAddress("").String() {
		writeErrRes(w, fmt.Errorf("%s is an invalid 'from' sender", from.String()))
		return
	}

	if req.FromPwd == "" {
		writeErrRes(w, fmt.Errorf("password to decrypt the %s account is required. 'from_pwd' is empty", from.String()))
		return
	}

	nonce := node.state.GetNextAccountNonce(from)
	tx := NewTx(from, NewAccount(req.To), req.Gas, req.GasPrice, req.Value, nonce, req.Data)

	signedTx, err := SignTxWithKeystoreAccount(tx, from, req.FromPwd, GetKeystoreDirPath(node.dataDir))
	if err != nil {
		writeErrRes(w, err)
		return
	}

	err = node.AddPendingTX(signedTx, node.info)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{Success: true})
}

func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	enableCors(&w)

	res := StatusRes{
		Hash:        node.state.LatestBlockHash(),
		Number:      node.state.LatestBlock().Header.Number,
		KnownPeers:  node.knownPeers,
		PendingTXs:  node.getPendingTXsAsArray(),
		NodeVersion: node.nodeVersion,
		Account:     NewAccount(node.info.Account.String()),
	}

	writeRes(w, res)
}

func syncHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash := Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, err)
		return
	}

	blocks, err := GetBlocksAfter(hash, node.dataDir)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, SyncRes{Blocks: blocks})
}

func addPeerHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	peerIP := r.URL.Query().Get(endpointAddPeerQueryKeyIP)
	peerPortRaw := r.URL.Query().Get(endpointAddPeerQueryKeyPort)
	minerRaw := r.URL.Query().Get(endpointAddPeerQueryKeyMiner)
	versionRaw := r.URL.Query().Get(endpointAddPeerQueryKeyVersion)

	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		writeRes(w, AddPeerRes{false, err.Error()})
		return
	}

	peer := NewPeerNode(peerIP, peerPort, false, NewAccount(minerRaw), true, versionRaw)

	node.AddPeer(peer)

	fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())

	writeRes(w, AddPeerRes{true, ""})
}
