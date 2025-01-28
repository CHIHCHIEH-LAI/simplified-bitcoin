package mining

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

func NewMessage(block *blockchain.Block, sender string) (*message.Message, error) {
	payload, err := block.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize block: %w", err)

	}

	return message.NewMessage(
		message.NEWBLOCK,
		sender,
		payload,
	), nil
}
