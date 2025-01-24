package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Transaction struct {
	TransactionID string
	Sender        string
	Recipient     string
	Amount        float64
	Fee           float64
	Timestamp     int64
	Signature     string
}

// NewUnsignedTransaction creates a new unsigned transaction
func NewUnsignedTransaction(sender, recipient string, amount, fee float64) (*Transaction, error) {
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

	return &tx, nil
}

// GenerateTransactionID generates a unique ID for the transaction
func (tx *Transaction) GenerateTransactionID() string {
	data := fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Serialize serializes the transaction into a string
func (tx *Transaction) Serialize() string {
	return fmt.Sprintf("%s|%s|%s|%f|%f|%d|%s", tx.TransactionID, tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp, tx.Signature)
}

// DeserializeTransaction deserializes the transaction from a string
func DeserializeTransaction(data string) (Transaction, error) {
	parts := strings.Split(data, "|")
	if len(parts) != 7 {
		return Transaction{}, fmt.Errorf("invalid transaction format")
	}

	tx := Transaction{
		TransactionID: parts[0],
		Sender:        parts[1],
		Recipient:     parts[2],
	}

	amount, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse amount: %v", err)
	}
	tx.Amount = amount

	fee, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse fee: %v", err)
	}
	tx.Fee = fee

	timestamp, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return Transaction{}, fmt.Errorf("failed to parse timestamp: %v", err)
	}
	tx.Timestamp = timestamp

	tx.Signature = parts[6]

	return tx, nil
}
