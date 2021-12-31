package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

var blockchain Blockchain = Blockchain{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	blockchain.Blocks = make([]Block, 0)
	blockchain.Nodes = make([]Node, 0)
	blockchain.Nodes = append(blockchain.Nodes, Node{port})
	blockchain.generateGenesisBlock()
	spew.Dump(blockchain)

	log.Fatal(run(port))
}

func run(Port string) error {
	mux := makeMuxRouter()
	log.Println("Listening on ", Port)

	s := &http.Server{
		Addr: ":" + Port,
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
