package blockchain

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
)

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(b *block.Block) error {
	// Validate the previous hash
	if err := bc.validatePrevHash(b); err != nil {
		return err
	}

	// Validate the block
	if err := b.Validate(); err != nil {
		return err
	}

	return nil
}

// validatePrevHash validates the previous hash
func (bc *Blockchain) validatePrevHash(b *block.Block) error {
	latestBlock := bc.GetLatestBlock()
	if b.PrevHash != latestBlock.BlockID {
		return fmt.Errorf("invalid previous hash: %s", b.PrevHash)
	}
	return nil
}
