package blockchain

import (
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Blockchain struct {
	Blocks []*Block    // Chain of blocks
	mutex  *sync.Mutex // Mutex to lock the blockchain
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
}

func (bc *Blockchain) NewBlock(transactions []*transaction.Transaction, miner string, reward float64) *Block {
	return NewBlock(bc.GetLatestBlock().BlockID, transactions, miner, reward)
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
