package transaction

import "fmt"

type TransactionPool struct {
	Transactions map[string]*Transaction
}

// NewTransactionPool creates a new transaction pool
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: make(map[string]*Transaction),
	}
}

// AddTransaction adds a transaction to the pool
func (pool *TransactionPool) AddTransaction(tx Transaction) error {
	if pool.Transactions[tx.TransactionID] != nil {
		return fmt.Errorf("Transaction with ID %s already exists", tx.TransactionID)
	}
	pool.Transactions[tx.TransactionID] = &tx
	return nil
}

// RemoveTransaction removes a transaction from the pool
func (pool *TransactionPool) RemoveTransaction(txID string) error {
	if pool.Transactions[txID] == nil {
		return fmt.Errorf("Transaction with ID %s does not exist", txID)
	}
	delete(pool.Transactions, txID)
	return nil
}
