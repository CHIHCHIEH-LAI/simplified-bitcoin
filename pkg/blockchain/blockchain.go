package blockchain

import (
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
)

type Blockchain struct {
	BaseReward float64
	Blocks     []*block.Block `json:"blocks"` // Blocks in the blockchain
	mutex      *sync.Mutex    // Mutex to protect the blockchain
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		BaseReward: 1000.0,
		Blocks:     []*block.Block{block.NewGenesisBlock()},
		mutex:      &sync.Mutex{},
	}
}

// NewBlock creates a new block with the given transactions
func (bc *Blockchain) NewBlock(transactions []*transaction.Transaction, miner string) *block.Block {
	prevHash := bc.GetLatestBlock().BlockID
	reward := bc.CalculateReward()
	difficulty := bc.CalculateDifficulty()
	return block.NewBlock(prevHash, transactions, miner, reward, difficulty)
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *block.Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}

	bc.Blocks = append(bc.Blocks, block)

	bc.printBlockchain()

	return nil
}

// GetLatestBlock returns the latest block in the blockchain
func (bc *Blockchain) GetLatestBlock() *block.Block {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	return bc.Blocks[len(bc.Blocks)-1]
}

// CalculateReward calculates the reward for the miner
func (bc *Blockchain) CalculateReward() float64 {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	return bc.BaseReward / float64(len(bc.Blocks))
}

// CalculateDifficulty calculates the difficulty for the miner
func (bc *Blockchain) CalculateDifficulty() int {
	return 5
}

// printBlockchain prints the blockchain
func (bc *Blockchain) printBlockchain() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	for _, block := range bc.Blocks {
		block.PrintBlock()
	}
}
