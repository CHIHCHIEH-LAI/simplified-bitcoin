package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Index        int64         // Position in the chain
	Timestamp    int64         // Unix timestamp
	Transactions []Transaction // List of transactions
	PrevHash     string        // Hash of the previous block
	Hash         string        // Hash of the block
	Nonce        int           // Proof of work
}

// NewGenesisBlock creates the first block in the blockchain
func NewGenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		PrevHash:     "",
		Hash:         "genesis",
		Nonce:        0,
	}
}

// NewBlock creates a new block in the blockchain
func NewBlock(index int64, transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         "",
		Nonce:        0,
	}
	block.Hash = block.CalculateHash()
	return block
}

// CalculateHash calculates the hash of the block
func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d%d%s%d", b.Index, b.Timestamp, b.PrevHash, b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
