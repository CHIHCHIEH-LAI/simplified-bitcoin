package block

import (
	"encoding/json"
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Block struct {
	BlockID      string                     `json:"block_id"`     // Hash of the block
	PrevHash     string                     `json:"prev_hash"`    // Hash of the previous block
	MerkleRoot   string                     `json:"merkle_root"`  // Merkle root of the transactions
	Timestamp    int64                      `json:"timestamp"`    // Unix timestamp
	Nonce        int                        `json:"nonce"`        // Proof of work
	Difficulty   int                        `json:"difficulty"`   // Difficulty of the block
	Transactions []*transaction.Transaction `json:"transactions"` // List of transactions
}

// Hash returns the hash of the block
func (b *Block) Hash() string {
	data := fmt.Sprintf("%s%s%d%d%d", b.PrevHash, b.MerkleRoot, b.Timestamp, b.Nonce, b.Difficulty)
	return utils.Hash(data)
}

// GenerateBlockID generates a unique ID for the block
func (b *Block) GenerateBlockID() string {
	return b.Hash()
}

// NewBlock creates a new block with the given previous hash and transactions
func NewBlock(prevHash string, transactions []*transaction.Transaction, miner string, reward float64, difficulty int) *Block {
	// Create a coinbase transaction to reward the miner
	coinbaseTx := transaction.NewCoinbaseTransaction(miner, reward)
	transactions = append([]*transaction.Transaction{coinbaseTx}, transactions...)

	// Compute the Merkle root
	merkleRoot := ComputeMerkleRoot(transactions)

	block := &Block{
		PrevHash:     prevHash,
		MerkleRoot:   merkleRoot,
		Timestamp:    0,
		Nonce:        0,
		Difficulty:   difficulty,
		Transactions: transactions,
	}
	block.BlockID = block.GenerateBlockID()
	return block
}

// NewGenesisBlock creates the first block in the blockchain
func NewGenesisBlock() *Block {
	return NewBlock("", nil, "", 0, 0)
}

// Serialize serializes the block to a JSON string
func (b *Block) Serialize() (string, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("failed to serialize block: %v", err)
	}
	return string(data), nil
}

// DeserializeBlock deserializes a JSON string to a block
func DeserializeBlock(data string) (*Block, error) {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize block: %v", err)
	}
	return &block, nil
}

// ComputeMerkleRoot computes the Merkle Root for a list of transactions
func ComputeMerkleRoot(transactions []*transaction.Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	// Step 1: Get the hash of each transaction
	var transactionHashes []string
	for _, tx := range transactions {
		transactionHashes = append(transactionHashes, tx.Hash())
	}

	// Step 2: Compute the Merkle Root from the transaction hashes
	for len(transactionHashes) > 1 {
		var newLevel []string

		// Process pairs of hashes
		for i := 0; i < len(transactionHashes); i += 2 {
			if i+1 < len(transactionHashes) {
				// Pair of hashes
				combinedHash := utils.HashPair(transactionHashes[i], transactionHashes[i+1])
				newLevel = append(newLevel, combinedHash)
			} else {
				// Odd hash (duplicate the last hash)
				combinedHash := utils.HashPair(transactionHashes[i], transactionHashes[i])
				newLevel = append(newLevel, combinedHash)
			}
		}

		transactionHashes = newLevel
	}

	// The last remaining hash is the Merkle Root
	return transactionHashes[0]
}
