package transaction

import (
	"fmt"
	"log"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

type Mempool struct {
	Transactions map[string]*Transaction // TransactionID -> Transaction
	Mutex        *sync.Mutex             // Mutex for the mempool
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make(map[string]*Transaction),
		Mutex:        &sync.Mutex{},
	}
}

// HandlenewTransaction handlkles a new transaction message
func (mp *Mempool) HandleNewTransaction(msg *message.Message) {
	// Deserialize the transaction
	tx, err := DeserializeTransaction(msg.Payload)
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
func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	if mp.Transactions[tx.TransactionID] != nil {
		return fmt.Errorf("Transaction with ID %s already exists", tx.TransactionID)
	}
	mp.Transactions[tx.TransactionID] = tx
	return nil
}

// RemoveTransaction removes a transaction from the pool
func (mp *Mempool) RemoveTransaction(txID string) error {
	if mp.Transactions[txID] == nil {
		return fmt.Errorf("Transaction with ID %s does not exist", txID)
	}
	delete(mp.Transactions, txID)
	return nil
}
