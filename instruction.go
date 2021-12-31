package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

// Web router
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("Post")
	muxRouter.HandleFunc("/nodes/register", handleRegisterNodes).Methods("Post")
	muxRouter.HandleFunc("/nodes", handleListNodes).Methods("GET")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	spew.Dump(blockchain)
	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	BPM int
	Ports []string
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

	for _, Port := range m.Ports {
		blockchain.Nodes = append(blockchain.Nodes, Node{Port})
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