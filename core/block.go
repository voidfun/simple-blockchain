package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	json "github.com/json-iterator/go"
)

const BlockReward = 100

type BlockHeader struct {
	Parent Hash           `json:"parent"`
	Number uint64         `json:"number"`
	Nonce  uint32         `json:"nonce"`
	Time   uint64         `json:"time"`
	Miner  common.Address `json:"miner"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []SignedTx  `json:"payload"`
}

func NewBlock(parent Hash, number uint64, nonce uint32, time uint64, miner common.Address, txs []SignedTx) Block {
	return Block{BlockHeader{parent, number, nonce, time, miner}, txs}
}

func (b *Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}

func (b *Block) GasReward() uint {
	reward := uint(0)

	for _, tx := range b.TXs {
		reward += tx.GasCost()
	}

	return reward
}

func IsBlockHashValid(hash Hash, miningDifficulty uint) bool {
	zeroesCount := uint(0)

	for i := uint(0); i < miningDifficulty; i++ {
		if fmt.Sprintf("%x", hash[i]) == "0" {
			zeroesCount++
		}
	}

	if fmt.Sprintf("%x", hash[miningDifficulty]) == "0" {
		return false
	}

	return zeroesCount == miningDifficulty
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}
