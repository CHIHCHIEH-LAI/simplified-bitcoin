package wallet

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

// CreateTransaction creates a new transaction
func (w *Wallet) CreateTransaction(recipient string, amount float64, fee float64) (*transaction.Transaction, error) {
	// Create the transaction
	tx := transaction.NewUnsignedTransaction(w.GetAddress(), recipient, amount, fee)

	// Sign the transaction
	data := fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
	signature, err := w.Sign(data)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature

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
	w.Transmitter.Transmit(message)

	return nil
}
