package p2p

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

func (node *Node) IntroduceSelfToGroup(bootstrapNodeAddress string) error {
	// Construct and serialize JOINREQ message
	sender := node.Address
	payload := node.Address
	message := NewMessage("JOINREQ", sender, payload)
	messageSerialized := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessageData(bootstrapNodeAddress, messageSerialized)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}
