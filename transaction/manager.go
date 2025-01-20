package transaction

import "fmt"

type TransactionManager struct {
	TransactionPool map[string]*Transaction
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		TransactionPool: make(map[string]*Transaction),
	}
}

// AddTransaction adds a transaction to the pool
func (mgr *TransactionManager) AddTransaction(tx Transaction) error {
	if mgr.TransactionPool[tx.TransactionID] != nil {
		return fmt.Errorf("Transaction with ID %s already exists", tx.TransactionID)
	}
	mgr.TransactionPool[tx.TransactionID] = &tx
	return nil
}

// RemoveTransaction removes a transaction from the pool
func (mgr *TransactionManager) RemoveTransaction(txID string) error {
	if mgr.TransactionPool[txID] == nil {
		return fmt.Errorf("Transaction with ID %s does not exist", txID)
	}
	delete(mgr.TransactionPool, txID)
	return nil
}
