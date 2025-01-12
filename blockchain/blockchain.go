package blockchain

import (
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/database"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Blockchain struct {
	Blocks              []Block
	PendingTransactions []Transaction
	Difficulty          int
	MiningReward        float64
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() Blockchain {
	return Blockchain{
		Blocks:              []Block{NewGenesisBlock()},
		PendingTransactions: []Transaction{},
		Difficulty:          2,
		MiningReward:        50.0,
	}
}

// GetLatestBlock returns the latest block in the blockchain
func GetLatestBlock(blockchain Blockchain) Block {
	return blockchain.Blocks[len(blockchain.Blocks)-1]
}

// AddBlock adds a block to the blockchain
func (blockchain *Blockchain) AddBlock(block Block) {
	blockchain.Blocks = append(blockchain.Blocks, block)
}

// SaveToDatabase saves the blockchain to the database
func (bc *Blockchain) SaveToDatabase(kvstore *database.KVStore) error {
	// Save the blocks
	hashes := []string{}
	for _, block := range bc.Blocks {
		err := block.SaveToDatabase(kvstore)
		if err != nil {
			return err
		}
		hashes = append(hashes, block.Hash)
	}

	// Serialize the block hashes
	data, err := utils.SerializeHashes(hashes)
	if err != nil {
		return err
	}

	// Save the block hashes
	return kvstore.Put("blockchain", []byte("blockchain"), data)
}
