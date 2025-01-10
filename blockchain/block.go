package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

func NewGenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    "2021-01-01",
		Transactions: []Transaction{},
		PrevHash:     "",
		Hash:         "",
		Nonce:        0,
	}
}

func NewBlock(index int, transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    "2021-01-01",
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         "",
		Nonce:        0,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%s%s%d", b.Index, b.Timestamp, b.PrevHash, b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
