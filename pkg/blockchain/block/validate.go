package block

import (
	"fmt"
	"strings"
)

// Validate validates the block
func (b *Block) Validate() error {
	// Validate the block ID
	if err := b.validateBlockID(); err != nil {
		return err
	}

	// Validate the difficulty
	if err := b.validateDifficulty(); err != nil {
		return err
	}

	// Validate the transactions
	if err := b.validateTransactions(); err != nil {
		return err
	}

	// Validate the Merkle root
	if err := b.validateMerkleRoot(); err != nil {
		return err
	}

	return nil
}

// ValidateBlockID validates the block ID
func (b *Block) validateBlockID() error {
	if b.BlockID != b.GenerateBlockID() {
		return fmt.Errorf("invalid block ID")
	}

	return nil
}

// ValidateDifficulty validates the difficulty
func (b *Block) validateDifficulty() error {
	prefix := strings.Repeat("0", b.Difficulty)
	if !strings.HasPrefix(b.BlockID, prefix) {
		return fmt.Errorf("invalid difficulty")
	}

	return nil
}

// ValidateTransactions validates the transactions
func (b *Block) validateTransactions() error {
	for _, tx := range b.Transactions {
		if err := tx.Validate(); err != nil {
			return fmt.Errorf("invalid transaction: %v", err)
		}
	}

	return nil
}

// ValidateMerkleRoot validates the Merkle root
func (b *Block) validateMerkleRoot() error {
	merkleRoot := ComputeMerkleRoot(b.Transactions)
	if b.MerkleRoot != merkleRoot {
		return fmt.Errorf("invalid Merkle root")
	}

	return nil
}
