package transaction

import (
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

func NewMessage(sender string, tx *Transaction) *message.Message {
	return message.NewMessage(
		message.NEWTRANSACTION,
		sender,
		tx.Serialize(),
	)
}
