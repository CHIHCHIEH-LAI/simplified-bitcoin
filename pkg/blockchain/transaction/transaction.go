package transaction

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
)

type Transaction struct {
	TransactionID string  `json:"transaction_id"`
	Sender        string  `json:"sender"`
	Recipient     string  `json:"recipient"`
	Amount        float64 `json:"amount"`
	Fee           float64 `json:"fee"`
	Timestamp     int64   `json:"timestamp"`
	Signature     string  `json:"signature"`
}

// NewUnsignedTransaction creates a new unsigned transaction
func NewUnsignedTransaction(sender, recipient string, amount, fee float64) *Transaction {
	// Create a new transaction
	tx := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}

	// Generate the transaction ID
	tx.TransactionID = tx.GenerateTransactionID()

	return &tx
}

func NewCoinbaseTransaction(miner string, reward float64) *Transaction {
	// Create a new transaction
	tx := NewUnsignedTransaction("coinbase", miner, reward, 0)

	return tx
}

// GenerateTransactionID generates a unique ID for the transaction
func (tx *Transaction) GenerateTransactionID() string {
	data := fmt.Sprintf("%s%s%f%f%d%s", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp, tx.Signature)
	return utils.Hash(data)
}

// Hash generates the hash of the transaction
func (tx *Transaction) Hash() string {
	data := tx.GenerateDataForSigning()
	return utils.Hash(data)
}

// GenerateDataForSigning generates the data that needs to be signed
func (tx *Transaction) GenerateDataForSigning() string {
	return fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
}

// Serialize serializes the transaction into a string
func (tx *Transaction) Serialize() (string, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %v", err)
	}
	return string(data), nil
}

// DeserializeTransaction deserializes the transaction from a string
func DeserializeTransaction(data string) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %v", err)
	}
	return &tx, nil
}

// SortTransactionsByFee sorts the transactions by fee
func SortTransactionsByFee(transactions []*Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Fee > transactions[j].Fee
	})
}
