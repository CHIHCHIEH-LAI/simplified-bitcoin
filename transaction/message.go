package transaction

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

func NewMessage(sender string, tx *Transaction) (*message.Message, error) {
	txData, err := tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	return message.NewMessage(
		message.NEWTRANSACTION,
		sender,
		"",
		txData,
	), nil
}
