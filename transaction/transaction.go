package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func (tx *Transaction) GenerateTransactionID() string {
	data := fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
