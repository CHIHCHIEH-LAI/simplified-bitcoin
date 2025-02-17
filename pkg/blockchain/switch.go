package blockchain

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
)

// ShouldSwitchChain determines if the current chain should be replaced with a new chain
func (bc *Blockchain) ShouldSwitchChain(fork *Blockchain) error {
	// Validate the new chain
	if err := fork.Validate(); err != nil {
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
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the new chain
	if err := bc.ShouldSwitchChain(fork); err != nil {
		return err
	}

	// Find the last common block
	lca := bc.FindLastCommonBlock(fork)

	// Remove diverging blocks
	bc.RemoveDivergingBlocks(lca)

	// Append the new blocks
	bc.AppendForkBlocks(fork.Blocks, lca)

	return nil
}

// FindLastCommonBlock finds the last common block between two chains
func (bc *Blockchain) FindLastCommonBlock(fork *Blockchain) *block.Block {
	for i := len(bc.Blocks) - 1; i >= 0; i-- {
		for j := len(fork.Blocks) - 1; j >= 0; j-- {
			if bc.Blocks[i].BlockID == fork.Blocks[j].BlockID {
				return bc.Blocks[i] // LCA found
			}
		}
	}
	return nil
}

// RemoveDivergingBlocks removes blocks that diverge from the last common block
func (bc *Blockchain) RemoveDivergingBlocks(lca *block.Block) {
	for i := len(bc.Blocks) - 1; i >= 0; i-- {
		if lca != nil && bc.Blocks[i].BlockID == lca.BlockID {
			break
		}

		// Add transactions back to the mempool
		for _, tx := range bc.Blocks[i].Transactions {
			bc.Mempool.AddTransaction(tx)
		}

		// Remove the block
		bc.Blocks = bc.Blocks[:len(bc.Blocks)-1]
	}
}

// AppendBlocks appends new blocks to the blockchain
func (bc *Blockchain) AppendForkBlocks(blocks []*block.Block, lca *block.Block) {
	for i := 0; i < len(blocks); i++ {
		if lca != nil && blocks[i].BlockID != lca.BlockID {
			continue
		}

		start := i + 1
		if lca == nil {
			start = 0
		}

		// Append the rest of the blocks
		for j := start; j < len(blocks); j++ {
			bc.Blocks = append(bc.Blocks, blocks[j])

			// Remove transactions from the mempool
			for _, tx := range blocks[j].Transactions {
				bc.Mempool.RemoveTransaction(tx.TransactionID)
			}
		}

		return
	}
}
