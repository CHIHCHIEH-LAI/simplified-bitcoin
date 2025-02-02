package blockchain

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain/block"

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(b *block.Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the block ID
	if err := b.ValidateBlockID(); err != nil {
		return err
	}

	// Validate the previous hash
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	if err := b.ValidatePrevHash(prevBlock.BlockID); err != nil {
		return err
	}

	return nil
}
