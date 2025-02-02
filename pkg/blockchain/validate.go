package blockchain

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(b *block.Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the block ID
	if err := b.ValidateBlockID(); err != nil {
		return err
	}

	// Validate the previous hash
	latestBlock := bc.GetLatestBlock()
	if err := b.ValidatePrevHash(latestBlock.BlockID); err != nil {
		return err
	}

	return nil
}
