package mempool

import (
	"fmt"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
)

type Mempool struct {
	Transactions map[string]*transaction.Transaction // TransactionID -> Transaction
	Mutex        *sync.RWMutex                       // Mutex for the mempool
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make(map[string]*transaction.Transaction),
		Mutex:        &sync.RWMutex{},
	}
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
		mp.RemoveTransaction(tx.TransactionID)
	}
}

func (mp *Mempool) RemoveTransactions(transactionIds []string) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	for _, txID := range transactionIds {
		mp.RemoveTransaction(txID)
	}
}

// RemoveTransaction removes a transaction from the pool
func (mp *Mempool) RemoveTransaction(txID string) error {
	if mp.Transactions[txID] == nil {
		return fmt.Errorf("transaction with ID %s does not exist", txID)
	}
	delete(mp.Transactions, txID)
	return nil
}

// GetTopNRewardingTransactions returns the top N rewarding transactions
func (mp *Mempool) GetTopNRewardingTransactions(n int) []*transaction.Transaction {
	mp.Mutex.RLock()
	defer mp.Mutex.RUnlock()

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
