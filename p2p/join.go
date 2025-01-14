package p2p

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

func (node *Node) IntroduceSelfToGroup(bootstrapNodeAddress string) error {
	// Construct and serialize JOINREQ message
	sender := node.Address
	payload := node.Address
	message := message.NewMessage("JOINREQ", sender, payload)
	messageData := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessageData(bootstrapNodeAddress, messageData)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}
