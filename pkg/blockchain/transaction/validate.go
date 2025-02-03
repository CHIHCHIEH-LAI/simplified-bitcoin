package transaction

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
)

// Validate checks if the transaction is valid
func (tx *Transaction) Validate() error {
	// Check if the transaction ID is valid
	if err := tx.validateTransactionID(); err != nil {
		return err
	}

	// Check if the sender is valid
	if err := tx.validateSender(); err != nil {
		return err
	}

	// Check if the recipient is valid
	if err := tx.validateRecipient(); err != nil {
		return err
	}

	// Check if the amount is valid
	if err := tx.validateAmount(); err != nil {
		return err
	}

	// Check if the fee is valid
	if err := tx.validateFee(); err != nil {
		return err
	}

	// Check if the timestamp is valid
	if err := tx.validateTimestamp(); err != nil {
		return err
	}

	// Check if the signature is valid
	if err := tx.validateSignature(); err != nil {
		return err
	}

	return nil
}

// validateTransactionID checks if the transaction ID is valid
func (tx *Transaction) validateTransactionID() error {
	if tx.TransactionID != tx.GenerateTransactionID() {
		return fmt.Errorf("invalid transaction ID")
	}
	return nil
}

// TODO: Implement validateSender
// validateSender checks if the sender is valid
func (tx *Transaction) validateSender() error {
	return nil
}

// TODO: Implement validateRecipient
// validateRecipient checks if the recipient is valid
func (tx *Transaction) validateRecipient() error {
	return nil
}

// validateAmount checks if the amount is valid
func (tx *Transaction) validateAmount() error {
	if tx.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	return nil
}

// validateFee checks if the fee is valid
func (tx *Transaction) validateFee() error {
	if tx.Fee < 0 {
		return fmt.Errorf("fee cannot be negative")
	}
	return nil
}

// validateTimestamp checks if the timestamp is valid
func (tx *Transaction) validateTimestamp() error {
	currentTime := utils.GetCurrentTimeInUnix()

	// Check if the timestamp is in the future
	if tx.Timestamp > currentTime {
		return fmt.Errorf("timestamp cannot be in the future")
	}

	// Check if the timestamp is too old
	if currentTime-tx.Timestamp > 60 {
		return fmt.Errorf("timestamp is too old")
	}

	return nil
}

// validateSignature checks if the signature is valid
func (tx *Transaction) validateSignature() error {
	data := tx.GenerateDataForSigning()
	if !utils.VerifySignature(tx.Sender, data, tx.Signature) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}
