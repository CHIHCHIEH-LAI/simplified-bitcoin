package blockchain

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
)

// Validate validates the blockchain
func (bc *Blockchain) Validate() error {
	// Validate the cumulative PoW
	if err := bc.validateCumulativePoW(); err != nil {
		return err
	}

	// Validate the blocks
	for i, b := range bc.Blocks[1:] {
		if err := bc.ValidateBlock(b, i); err != nil {
			return fmt.Errorf("invalid block: %v", err)
		}
	}

	return nil
}

// validateCumulativePoW validates the cumulative PoW
func (bc *Blockchain) validateCumulativePoW() error {
	cumulativePoW := bc.CalculateCumulativePoW()
	if bc.CumulativePoW != cumulativePoW {
		return fmt.Errorf("invalid cumulative PoW: %d", bc.CumulativePoW)
	}
	return nil
}

// ValidateNewBlock validates the new block
func (bc *Blockchain) ValidateNewBlock(b *block.Block) error {
	// Validate the previous hash
	height := len(bc.Blocks) - 1
	if err := bc.validatePrevHash(b, height); err != nil {
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

// ValidateBlock validates the block
func (bc *Blockchain) ValidateBlock(b *block.Block, height int) error {
	// Validate the previous hash
	if err := bc.validatePrevHash(b, height); err != nil {
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
func (bc *Blockchain) validatePrevHash(b *block.Block, height int) error {
	if height == 0 {
		return nil
	}

	if height >= len(bc.Blocks) {
		return fmt.Errorf("invalid height: %d", height)
	}

	prevBlock := bc.Blocks[height]
	if b.PrevHash != prevBlock.BlockID {
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
	reward := bc.CalculateReward(b.Transactions)
	coinbaseTx := b.Transactions[0]
	if coinbaseTx.Amount != reward {
		return fmt.Errorf("invalid reward: %f", coinbaseTx.Amount)
	}

	return nil
}

func (bc *Blockchain) ValidateTransaction(tx *transaction.Transaction) error {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	// Validate the transaction
	if err := tx.Validate(); err != nil {
		return err
	}

	// Validate the unspent transaction outputs
	if err := bc.validateUTXOs(tx); err != nil {
		return err
	}

	return nil
}

// validateUTXOs validates the unspent transaction outputs
func (bc *Blockchain) validateUTXOs(tx *transaction.Transaction) error {
	// Get the unspent transaction outputs
	utxos := bc.calculateUTXOs(tx.Sender)

	// Validate the sender's balance
	if utxos < tx.Amount+tx.Fee {
		return fmt.Errorf("insufficient balance: %f", utxos)
	}

	return nil
}
