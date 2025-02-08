package blockchain

import "fmt"

// ShouldSwitchChain determines if the current chain should be replaced with a new chain
func (bc *Blockchain) ShouldSwitchChain(fork *Blockchain) error {

	// Validate the new chain
	if err := bc.ValidateChain(fork); err != nil {
		return err
	}

	// Compare the cumulative difficulty of the two chains
	if fork.CumulativePoW <= bc.CumulativePoW {
		return fmt.Errorf("new chain has lower cumulative PoW")
	}

	return nil
}

// SwitchChain switches the current chain with a new chain
func (bc *Blockchain) SwitchChain(fork *Blockchain) error {

	// Validate the new chain
	if err := bc.ValidateChain(fork); err != nil {
		return err
	}

	// Compare the cumulative difficulty of the two chains
	if fork.CumulativePoW <= bc.CumulativePoW {
		return fmt.Errorf("new chain has lower cumulative PoW")
	}

	// Replace the chain
	bc.Blocks = fork.Blocks
	bc.CumulativePoW = fork.CumulativePoW

	return nil
}
