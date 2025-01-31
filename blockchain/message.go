package blockchain

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"

func NewMinedBlockMessage(sender string, minedBlock *Block) *message.Message {
	blockData, err := minedBlock.Serialize()
	if err != nil {
		return nil
	}
	return message.NewMessage(
		message.NEWBLOCK,
		sender,
		"",
		blockData,
	)
}
