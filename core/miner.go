package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type PendingBlock struct {
	parent Hash
	number uint64
	time   uint64
	miner  common.Address
	txs    []SignedTx
}

func NewPendingBlock(parent Hash, number uint64, miner common.Address, txs []SignedTx) PendingBlock {
	return PendingBlock{parent, number, uint64(time.Now().Unix()), miner, txs}
}

func Mine(ctx context.Context, pb PendingBlock, miningDifficulty uint) (Block, error) {
	if len(pb.txs) == 0 {
		return Block{}, fmt.Errorf("mining empty blocks is not allowed")
	}

	start := time.Now()
	attempt := 0
	var block Block
	var hash Hash
	var nonce uint32

	for !IsBlockHashValid(hash, miningDifficulty) {
		select {
		case <-ctx.Done():
			fmt.Println("Mining cancelled!")

			return Block{}, fmt.Errorf("mining cancelled. %s", ctx.Err())
		default:
		}

		attempt++
		nonce = generateNonce()

		if attempt%1000000 == 0 || attempt == 1 {
			fmt.Printf("Mining %d Pending TXs. Attempt: %d\n", len(pb.txs), attempt)
		}

		block = NewBlock(pb.parent, pb.number, nonce, pb.time, pb.miner, pb.txs)
		blockHash, err := block.Hash()
		if err != nil {
			return Block{}, fmt.Errorf("couldn't mine block. %s", err.Error())
		}

		hash = blockHash
	}

	fmt.Printf("\nMined new Block '%x' using PoW 🎉🎉🎉\n", hash)
	fmt.Printf("\tHeight: '%v'\n", block.Header.Number)
	fmt.Printf("\tNonce: '%v'\n", block.Header.Nonce)
	fmt.Printf("\tCreated: '%v'\n", block.Header.Time)
	fmt.Printf("\tMiner: '%v'\n", block.Header.Miner.String())
	fmt.Printf("\tParent: '%v'\n\n", block.Header.Parent.Hex())

	fmt.Printf("\tAttempt: '%v'\n", attempt)
	fmt.Printf("\tTime: %s\n\n", time.Since(start))

	return block, nil
}

func generateNonce() uint32 {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Uint32()
}
