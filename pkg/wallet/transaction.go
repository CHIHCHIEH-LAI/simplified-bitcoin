package wallet

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
)

// CreateTransaction creates a new transaction
func (w *Wallet) CreateTransaction(recipient string, amount float64, fee float64) (*transaction.Transaction, error) {
	// Create the transaction
	tx := transaction.NewUnsignedTransaction(w.GetAddress(), recipient, amount, fee)

	// Sign the transaction
	hash := tx.Hash()
	signature, err := w.Sign(hash)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature

	// Generate the transaction ID
	tx.TransactionID = tx.GenerateTransactionID()

	return tx, nil
}

// SendTransaction sends a transaction to a node in the network
func (w *Wallet) SendTransaction(tx *transaction.Transaction, selfAddress string, nodeAddress string) error {
	// Create and serialize the mwssage
	message, err := transaction.NewMessage(
		selfAddress,
		nodeAddress,
		tx,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction message: %v", err)
	}

	// Send the message to the node
	w.Transmitter.SendMessage(message)

	return nil
}
