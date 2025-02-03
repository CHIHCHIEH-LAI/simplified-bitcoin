package transaction

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
)

func NewMessage(sender, receipient string, tx *Transaction) (*message.Message, error) {
	txData, err := tx.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	return message.NewMessage(
		message.NEWTRANSACTION,
		sender,
		receipient,
		txData,
	), nil
}
