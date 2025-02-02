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

	// Validate the difficulty
	if err := bc.validateDifficulty(b); err != nil {
		return err
	}

	// Validate the reward
	if err := bc.validateReward(b); err != nil {
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

// validateDifficulty validates the difficulty
func (bc *Blockchain) validateDifficulty(b *block.Block) error {
	if b.Difficulty != bc.CalculateDifficulty() {
		return fmt.Errorf("invalid difficulty: %d", b.Difficulty)
	}
	return nil
}

// validateReward validates the reward
func (bc *Blockchain) validateReward(b *block.Block) error {
	reward := bc.CalculateReward()
	coinbaseTx := b.Transactions[0]
	if coinbaseTx.Amount != reward {
		return fmt.Errorf("invalid reward: %f", coinbaseTx.Amount)
	}

	return nil
}
