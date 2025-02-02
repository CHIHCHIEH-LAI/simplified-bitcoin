package block

import "fmt"

// ValidateBlockID validates the block ID
func (b *Block) ValidateBlockID() error {
	if b.BlockID != b.GenerateBlockID() {
		return fmt.Errorf("invalid block ID")
	}

	return nil
}

// ValidatePrevHash validates the previous hash
func (b *Block) ValidatePrevHash(prevHash string) error {
	if b.PrevHash != prevHash {
		return fmt.Errorf("invalid previous hash")
	}

	return nil
}
