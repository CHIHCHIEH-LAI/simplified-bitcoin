package blockchain

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Block struct {
	BlockID      string
	PrevHash     string
	Timestamp    int64
	Nonce        int // Nonce is a number that miners use to change the hash of the block
	Transactions []*transaction.Transaction
}

// Hash returns the hash of the block
func (b *Block) Hash() string {
	data := fmt.Sprintf("%s%d%d", b.PrevHash, b.Timestamp, b.Nonce)
	return utils.Hash(data)
}

// GenerateBlockID generates a unique ID for the block
func (b *Block) GenerateBlockID() string {
	return b.Hash()
}

// NewBlock creates a new block with the given previous hash and transactions
func NewBlock(prevHash string, transactions []*transaction.Transaction, miner string, reward float64) *Block {
	// Create a coinbase transaction to reward the miner
	coinbaseTx := transaction.NewCoinbaseTransaction(miner, reward)
	transactions = append([]*transaction.Transaction{coinbaseTx}, transactions...)

	block := &Block{
		PrevHash:     prevHash,
		Timestamp:    0,
		Nonce:        0,
		Transactions: transactions,
	}
	block.BlockID = block.GenerateBlockID()
	return block
}

// NewGenesisBlock creates the first block in the blockchain
func NewGenesisBlock() *Block {
	return NewBlock("", nil, "", 0)
}

// ValidateBlockID validates the block ID
func (b *Block) ValidateBlockID() error {
	if b.BlockID != b.GenerateBlockID() {
		return fmt.Errorf("invalid block ID")
	}

	return nil
}

// ValidatePrevHash validates the previous hash
func (b *Block) ValidatePrevHash(prevHash string) error {
	if b.PrevHash != prevHash {
		return fmt.Errorf("invalid previous hash")
	}

	return nil
}
