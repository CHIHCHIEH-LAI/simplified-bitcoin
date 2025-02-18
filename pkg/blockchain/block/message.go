package block

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"

func NewMinedBlockMessage(minedBlock *Block, sender string) *message.Message {
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
