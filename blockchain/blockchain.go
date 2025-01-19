package blockchain

type Blockchain struct {
	Blocks []*Block // Chain of blocks
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {
	// Validate the block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}

	bc.Blocks = append(bc.Blocks, block)
	return nil
}

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(block *Block) error {
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
