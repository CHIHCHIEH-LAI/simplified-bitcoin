package blockchain

import (
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
)

type Blockchain struct {
	BaseReward    float64
	Blocks        []*block.Block `json:"blocks"` // Blocks in the blockchain
	mutex         *sync.RWMutex  // Mutex to protect the blockchain
	CumulativePoW int            `json:"cumulativePoW"` // Tracks total proof-of-work (sum of difficulties)
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	genesisBlock := block.NewGenesisBlock()
	return &Blockchain{
		BaseReward:    1000.0,
		Blocks:        []*block.Block{genesisBlock},
		mutex:         &sync.RWMutex{},
		CumulativePoW: genesisBlock.Difficulty,
	}
}

// NewBlock creates a new block with the given transactions
func (bc *Blockchain) NewBlock(transactions []*transaction.Transaction, miner string) *block.Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

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
	bc.CumulativePoW += block.Difficulty

	return nil
}

// GetLatestBlock returns the latest block in the blockchain
func (bc *Blockchain) GetLatestBlock() *block.Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// CalculateReward calculates the reward for the miner
func (bc *Blockchain) CalculateReward() float64 {
	return bc.BaseReward / float64(len(bc.Blocks))
}

// CalculateDifficulty calculates the difficulty for the miner
func (bc *Blockchain) CalculateDifficulty() int {
	return 10
}

// printBlockchain prints the blockchain
func (bc *Blockchain) PrintBlockchain() {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for _, block := range bc.Blocks {
		block.PrintBlock()
	}
}
