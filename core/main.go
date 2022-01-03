package main

import (
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

var blockchain Blockchain = Blockchain{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	url := "http://" + host + ":" + port
	blockchain.Blocks = make([]Block, 0)
	node := Node{url}
	blockchain.SelfNode = node
	blockchain.Nodes = make(map[string]Node)
	blockchain.Nodes[url] = node
	blockchain.generateGenesisBlock()
	// spew.Dump(blockchain)

	log.Fatal(run(port))
}

func run(port string) error {
	mux := makeMuxRouter()
	log.Println("Listening on ", port)

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
