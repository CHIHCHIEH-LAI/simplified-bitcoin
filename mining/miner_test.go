package mining

import (
	"testing"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	// Setup a mock transaction pool
	tx1 := &transaction.Transaction{Sender: "Alice", Recipient: "Bob", Amount: 10, Fee: 1}
	tx2 := &transaction.Transaction{Sender: "Charlie", Recipient: "Dave", Amount: 5, Fee: 0.5}
	transactions := []*transaction.Transaction{tx1, tx2}

	miner := NewMiner(&transactions, 2)

	// Get transactions
	transactions = miner.GetTransactions()

	// Verify the transactions are fetched correctly
	assert.Equal(t, 2, len(transactions), "Expected 2 transactions")
	assert.Equal(t, tx1, transactions[0], "First transaction should match")
	assert.Equal(t, tx2, transactions[1], "Second transaction should match")
}

// func TestStartMining(t *testing.T) {
// 	// Setup a mock blockchain and transaction pool
// 	mockBlockchain := blockchain.NewBlockchain()
// 	tx1 := &transaction.Transaction{Sender: "Alice", Recipient: "Bob", Amount: 10}
// 	transactionPool := []*transaction.Transaction{tx1}

// 	miner := &Miner{
// 		Transactions: &transactionPool,
// 		Blockchain:   mockBlockchain,
// 		Difficulty:   2,
// 		StopMining:   make(chan bool),
// 		MutexPool:    &sync.Mutex{},
// 	}

// 	// Start mining
// 	go miner.StartMining("MinerAddress123")

// 	// Wait for mining to complete
// 	select {
// 	case <-miner.StopMining:
// 		t.Fatal("Mining was stopped prematurely")
// 	default:
// 		// Verify that a block was mined and added to the blockchain
// 		if len(mockBlockchain.Chain) != 2 {
// 			t.Fatalf("Expected blockchain to have 2 blocks, got %d", len(mockBlockchain.Chain))
// 		}

// 		// Check the latest block
// 		latestBlock := mockBlockchain.GetLatestBlock()
// 		assert.Equal(t, "MinerAddress123", latestBlock.Transactions[0].Recipient, "Coinbase transaction recipient mismatch")
// 		assert.Equal(t, 50.0, latestBlock.Transactions[0].Amount, "Coinbase transaction reward mismatch")
// 	}
// }

// func TestPerformProofOfWork(t *testing.T) {
// 	// Setup a mock block
// 	block := &blockchain.Block{
// 		Transactions: []*transaction.Transaction{},
// 		Nonce:        0,
// 	}

// 	// Setup a miner
// 	mockBlockchain := blockchain.NewBlockchain()
// 	miner := &Miner{
// 		Blockchain: mockBlockchain,
// 		Difficulty: 1, // Set low difficulty for testing
// 	}

// 	// Perform proof of work
// 	miner.PerformProofOfWork(block)

// 	// Verify the block hash meets the difficulty target
// 	assert.True(t, strings.HasPrefix(block.BlockID, "0"), "Block hash does not meet difficulty target")
// }

// func TestStopMining(t *testing.T) {
// 	// Setup a mock blockchain and transaction pool
// 	mockBlockchain := blockchain.NewBlockchain()
// 	tx1 := &transaction.Transaction{Sender: "Alice", Recipient: "Bob", Amount: 10}
// 	transactionPool := []*transaction.Transaction{tx1}

// 	miner := &Miner{
// 		Transactions: &transactionPool,
// 		Blockchain:   mockBlockchain,
// 		Difficulty:   2,
// 		StopMining:   make(chan bool),
// 		MutexPool:    &sync.Mutex{},
// 	}

// 	// Start mining
// 	go miner.StartMining("MinerAddress123")

// 	// Stop mining after a short delay
// 	go func() {
// 		miner.Stop()
// 	}()

// 	// Wait and verify mining was stopped
// 	select {
// 	case <-miner.StopMining:
// 		t.Log("Mining process stopped as expected")
// 	default:
// 		t.Fatal("Mining process did not stop")
// 	}
// }
