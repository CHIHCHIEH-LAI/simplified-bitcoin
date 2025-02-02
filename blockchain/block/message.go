package block

import "github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"

func NewMinedBlockMessage(minedBlock *Block) *message.Message {
	blockData, err := minedBlock.Serialize()
	if err != nil {
		return nil
	}
	return message.NewMessage(
		message.NEWBLOCK,
		"",
		"",
		blockData,
	)
}
