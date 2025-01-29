package transaction

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type TransactionManager struct {
	Sender          string
	TransactionPool map[string]*Transaction
	Mutex           *sync.Mutex
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(address string) *TransactionManager {
	return &TransactionManager{
		Sender:          address,
		TransactionPool: make(map[string]*Transaction),
		Mutex:           &sync.Mutex{},
	}
}

func (mgr *TransactionManager) HandleNewTransaction(msg *message.Message) {
	// Check if the transaction is stale
	if math.Abs(float64(msg.Timestamp-utils.GetCurrentTimeInUnix())) > TIME_VALID_TX_THRESHOLD {
		log.Printf("Received stale transaction message\n")
		return
	}

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
	mgr.AddTransaction(tx)
}

// AddTransaction adds a transaction to the pool
func (mgr *TransactionManager) AddTransaction(tx *Transaction) error {
	mgr.Mutex.Lock()
	defer mgr.Mutex.Unlock()

	if mgr.TransactionPool[tx.TransactionID] != nil {
		return fmt.Errorf("Transaction with ID %s already exists", tx.TransactionID)
	}
	mgr.TransactionPool[tx.TransactionID] = tx
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
