package transaction

import "fmt"

// Validate checks if the transaction is valid
func (tx *Transaction) Validate() error {
	// Check if the transaction ID is valid
	if tx.TransactionID == "" {
		return fmt.Errorf("transaction ID is required")
	}

	// Check if the sender is valid
	if tx.Sender == "" {
		return fmt.Errorf("sender is required")
	}

	// Check if the recipient is valid
	if tx.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}

	// Check if the amount is valid
	if tx.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	// Check if the fee is valid
	if tx.Fee < 0 {
		return fmt.Errorf("fee cannot be negative")
	}

	// Check if the timestamp is valid
	if tx.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be greater than zero")
	}

	// Check if the signature is valid
	if tx.Signature == "" {
		return fmt.Errorf("signature is required")
	}

	return nil
}
