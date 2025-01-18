package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Block struct {
	BlockID      string
	PrevHash     string
	Timestamp    int64
	Nonce        int // Nonce is a number that miners use to change the hash of the block
	Transactions []*transaction.Transaction
}

// GenerateBlockID generates a unique ID for the block
func (b *Block) GenerateBlockID() string {
	data := fmt.Sprintf("%s%d%d", b.PrevHash, b.Timestamp, b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// NewBlock creates a new block with the given previous hash and transactions
func NewBlock(prevHash string, transactions []*transaction.Transaction) *Block {
	block := &Block{
		PrevHash:     prevHash,
		Timestamp:    0,
		Nonce:        0,
		Transactions: transactions,
	}
	block.BlockID = block.GenerateBlockID()
	return block
}
