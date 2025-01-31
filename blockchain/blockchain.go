package blockchain

import (
	"math"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Blockchain struct {
	BaseReward     float64
	BaseMiningTime int64
	Blocks         []*Block    `json:"blocks"` // Blocks in the blockchain
	mutex          *sync.Mutex // Mutex to protect the blockchain
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		BaseReward:     1000.0,
		BaseMiningTime: 10 * 60,
		Blocks:         []*Block{NewGenesisBlock()},
		mutex:          &sync.Mutex{},
	}
}

// NewBlock creates a new block with the given transactions
func (bc *Blockchain) NewBlock(transactions []*transaction.Transaction, miner string) *Block {
	prevHash := bc.GetLatestBlock().BlockID
	reward := bc.CalculateReward()
	difficulty := bc.CalculateDifficulty()
	return NewBlock(prevHash, transactions, miner, reward, difficulty)
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {

	// Validate the block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	bc.Blocks = append(bc.Blocks, block)

	return nil
}

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(block *Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the block ID
	if err := block.ValidateBlockID(); err != nil {
		return err
	}

	// Validate the previous hash
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	if err := block.ValidatePrevHash(prevBlock.BlockID); err != nil {
		return err
	}

	return nil
}

// GetLatestBlock returns the latest block in the blockchain
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	return bc.Blocks[len(bc.Blocks)-1]
}

// CalculateReward calculates the reward for the miner
func (bc *Blockchain) CalculateReward() float64 {
	return bc.BaseReward / float64(len(bc.Blocks))
}

// CalculateDifficulty calculates the difficulty for the miner
func (bc *Blockchain) CalculateDifficulty() int {
	if len(bc.Blocks) < 2 {
		return 5
	}

	lastBlockMinedTime := bc.GetLatestBlock().Timestamp - bc.Blocks[len(bc.Blocks)-2].Timestamp
	lastBlockMinedTime += 1
	difficulty := int(math.Max(5, float64(bc.BaseMiningTime/lastBlockMinedTime)))
	return difficulty
}
