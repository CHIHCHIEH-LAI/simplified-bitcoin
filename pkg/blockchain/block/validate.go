package block

import "fmt"

// Validate validates the block
func (b *Block) Validate() error {
	// Validate the block ID
	if err := b.validateBlockID(); err != nil {
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
