package wallet

import (
	"fmt"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

// CreateTransaction creates a new transaction
func (w *Wallet) CreateTransaction(recipient string, amount float64, fee float64) (*transaction.Transaction, error) {
	// Create the transaction
	tx := transaction.Transaction{
		Sender:    w.GetAddress(),
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
	}

	// Generate the transaction ID
	tx.TransactionID = tx.GenerateTransactionID()

	// Sign the transaction
	data := fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
	signature, err := w.Sign(data)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature

	return &tx, nil
}

// SendTransaction sends a transaction to a node in the network
func SendTransaction(tx *transaction.Transaction, selfAddress string, nodeAddress string) error {
	// Serialize the transaction
	txData := tx.Serialize()

	// Create and serialize the mwssage
	message := message.Message{
		Type:    message.NEWTRANSACTION,
		Sender:  selfAddress,
		Payload: txData,
	}
	messageData := message.Serialize()

	// Send the message to the node
	err := network.SendMessageData(nodeAddress, messageData)
	if err != nil {
		return err
	}

	return nil
}
