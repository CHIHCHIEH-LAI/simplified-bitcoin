package p2p

import (
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type JOINREQMessagePayload struct {
	BitcoinAddress string
	Address        string
}

// NewJOINREQMessagePayload creates a new JOINREQ message payload
func NewJOINREQMessagePayload(bitcoinAddress, address string) *JOINREQMessagePayload {
	return &JOINREQMessagePayload{
		BitcoinAddress: bitcoinAddress,
		Address:        address,
	}
}

func (payload *JOINREQMessagePayload) Serialize() string {
	return fmt.Sprintf("%s|%s", payload.BitcoinAddress, payload.Address)
}

func (node *P2PNode) IntroduceSelfToGroup(bootstrapNodeAddress string) error {
	// Construct and serialize JOINREQ payload
	payload := NewJOINREQMessagePayload(node.BitcoinAddress, node.Address)
	payloadSerialized := payload.Serialize()

	// Construct and serialize JOINREQ message
	message := network.NewMessage("JOINREQ", node.Address, payloadSerialized)
	messageSerialized := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessage(bootstrapNodeAddress, messageSerialized)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}
