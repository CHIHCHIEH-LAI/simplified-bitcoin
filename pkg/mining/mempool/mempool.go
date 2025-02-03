package mempool

import (
	"fmt"
	"log"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
)

type Mempool struct {
	Transactions map[string]*transaction.Transaction // TransactionID -> Transaction
	Mutex        *sync.Mutex                         // Mutex for the mempool
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make(map[string]*transaction.Transaction),
		Mutex:        &sync.Mutex{},
	}
}

// HandlenewTransaction handlkles a new transaction message
func (mp *Mempool) HandleNewTransaction(msg *message.Message) {
	// Deserialize the transaction
	tx, err := transaction.DeserializeTransaction(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize transaction: %v\n", err)
		return
	}

	// Validate the transaction
	if err := tx.Validate(); err != nil {
		log.Printf("Invalid transaction: %v\n", err)
		return
	}

	// Add the transaction to the pool
	mp.AddTransaction(tx)
}

// AddTransaction adds a transaction to the pool
func (mp *Mempool) AddTransaction(tx *transaction.Transaction) error {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	if mp.Transactions[tx.TransactionID] != nil {
		return fmt.Errorf("transaction with ID %s already exists", tx.TransactionID)
	}
	mp.Transactions[tx.TransactionID] = tx
	return nil
}

// RemoveTransactionsInBlock removes transactions in a block from the pool
func (mp *Mempool) RemoveTransactionsInBlock(block *block.Block) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	for _, tx := range block.Transactions {
		delete(mp.Transactions, tx.TransactionID)
	}
}

// RemoveTransaction removes a transaction from the pool
func (mp *Mempool) RemoveTransaction(txID string) error {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	if mp.Transactions[txID] == nil {
		return fmt.Errorf("transaction with ID %s does not exist", txID)
	}
	delete(mp.Transactions, txID)
	return nil
}

// GetTopNRewardingTransactions returns the top N rewarding transactions
func (mp *Mempool) GetTopNRewardingTransactions(n int) []*transaction.Transaction {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	// Convert the map to a slice
	txSlice := make([]*transaction.Transaction, 0, len(mp.Transactions))
	for _, tx := range mp.Transactions {
		txSlice = append(txSlice, tx)
	}

	// Sort the transactions by fee
	transaction.SortTransactionsByFee(txSlice)

	// Get the top N rewarding transactions
	if n > len(txSlice) {
		n = len(txSlice)
	}
	return txSlice[:n]
}
