package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/database"
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

// SaveToDatabase saves the block to the database
func (b *Block) SaveToDatabase(kvstore *database.KVStore) error {
	data, err := b.Serialize()
	if err != nil {
		return err
	}

	key := []byte(b.Hash)
	return kvstore.Put("blocks", key, data)

}

// LoadBlockFromDatabase loads a block from the database using its hash
func LoadBlockFromDatabase(kvstore *database.KVStore, hash string) (Block, error) {
	key := []byte(hash)
	data, err := kvstore.Get("blocks", key)
	if err != nil {
		return Block{}, err
	}

	return DeserializeBlock(data)
}

// Serialize serializes the block into a byte slice
func (b *Block) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

// DeserializeBlock deserializes a block from a byte slice
func DeserializeBlock(data []byte) (Block, error) {
	var block Block
	err := json.Unmarshal(data, &block)
	return block, err
}
