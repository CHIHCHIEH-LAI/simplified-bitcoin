package transaction

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"

func NewTransactionMessage(tx *Transaction, sender string) *message.Message {
	return &message.Message{
		Type:    message.NEWTRANSACTION,
		Sender:  sender,
		Payload: tx.Serialize(),
	}
}
