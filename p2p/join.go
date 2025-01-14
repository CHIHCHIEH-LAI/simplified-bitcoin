package p2p

import (
	"fmt"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

// IntroduceSelfToGroup sends a JOINREQ message to the bootstrap node
func (node *Node) IntroduceSelfToGroup(bootstrapNodeAddress string) error {
	message := message.NewMessage("JOINREQ", node.Address, "")
	messageData := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessageData(bootstrapNodeAddress, messageData)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}

// HandleJoinRequest processes a JOINREQ message
func (node *Node) HandleJoinRequest(msg message.Message) {
	// Add the sender to the member list
	member := Member{
		Address:   msg.Sender,
		Heartbeat: 0,
		Timestamp: time.Now().Unix(),
	}
	node.MemberList = append(node.MemberList, member)

	// Send JOINREP message to the sender

}
