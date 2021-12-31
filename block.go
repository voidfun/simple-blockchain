package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Index int
	Timestamp string
	BPM int
	Hash string
	PrevHash string
}

type Node struct {
	Port string
}

func (n *Node) httpAddr() string {
	return "http://localhost:" + n.Port
}

type Blockchain struct {
	Blocks []Block
	Nodes []Node
}

func (bc *Blockchain) replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(bc.Blocks) {
		bc.Blocks = newBlocks
	}
}

func (bc *Blockchain) generateGenesisBlock() {
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	spew.Dump(genesisBlock)
	bc.Blocks = append(bc.Blocks, genesisBlock)
}

func (bc *Blockchain) generateBlock(BPM int) (Block, error) {
	t := time.Now()
	var newBlock Block
	oldBlock := bc.Blocks[len(bc.Blocks) - 1]
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

func (bc *Blockchain) addBlockRecord(BPM int) (bool, error) {
	newBlock, err := bc.generateBlock(BPM)
	if err != nil {
		return false, err
	}
	return bc.addBlock(newBlock), nil
}

func (bc *Blockchain) addBlock(newBlock Block) bool {
	oldBlock := bc.Blocks[len(bc.Blocks) - 1]
	if isBlockValid(newBlock, oldBlock) {
		newBlockchain := append(bc.Blocks, newBlock)
		bc.replaceChain(newBlockchain)
		spew.Dump(bc)
		return true
	} else {
		return false
	}
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write(([]byte(record)))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

