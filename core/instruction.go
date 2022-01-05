package main

/*
import (
	"encoding/json"
	"io"
	"net/http"

	// "github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)


// Web router
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	// public
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("Post")
	muxRouter.HandleFunc("/nodes", handleRegisterNodes).Methods("Post")
	muxRouter.HandleFunc("/nodes", handleListNodes).Methods("GET")
	muxRouter.HandleFunc("/syncAll", handleSyncAll).Methods("GET")

	// private
	muxRouter.HandleFunc("/blocks", handleListBlocks).Methods("GET")
	muxRouter.HandleFunc("/syncBlocks", handleSyncBlocks).Methods("GET")
	muxRouter.HandleFunc("/syncNodes", handleSyncNodes).Methods("GET")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	BPM int
	URLs []string
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := blockchain.addBlockRecord(m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func handleRegisterNodes(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	for _, url := range m.URLs {
		node := Node{url}
		blockchain.Nodes[url] = node
	}
}

func handleListNodes(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.Nodes, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleListBlocks(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.Blocks, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleSyncBlocks(w http.ResponseWriter, r *http.Request) {
	blockchain.syncBlocks()
	// spew.Dump(blockchain)
	bytes, err := json.MarshalIndent(blockchain.Blocks, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleSyncAll(w http.ResponseWriter, r *http.Request) {
	blockchain.syncAll()
	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleSyncNodes(w http.ResponseWriter, r *http.Request) {
	blockchain.syncNodes()
	bytes, err := json.MarshalIndent(blockchain.Nodes, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
*/
