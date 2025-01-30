package mempool

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

func NewMessage(sender, receipient string, tx *transaction.Transaction) (*message.Message, error) {
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
