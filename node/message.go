package node

import (
	"log"
	"math"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
)

// HandleIncomingMessage processes incoming messages
func (node *Node) HandleIncomingMessage() {
	for {
		// Receive a message
		msg := node.Tranceiver.Receive()

		// Process the message based on its type
		switch msg.Type {
		case message.JOINREQ:
			node.MembershipManager.HandleJoinRequest(msg)
		case message.JOINRESP:
			node.MembershipManager.HandleJoinResponse(msg)
		case message.HEARTBEAT:
			node.MembershipManager.HandleHeartbeat(msg)
		case message.NEWTRANSACTION:
			go node.GossipMessage(msg)
			node.TransactionManager.HandleNewTransaction(msg)
		// case message.NEWBLOCK:
		// 	node.HandleNewBlock(msg)
		// case "GETBLOCKCHAIN":
		// 	node.HandleGetBlockchain(msg)
		default:
			log.Printf("Unknown message type: %s\n", msg.Type)
		}
	}
}

// GossipMessage sends a message to N random members
func (node *Node) GossipMessage(msg *message.Message) error {

	// Select N random members to send the message to
	n_members := len(node.MembershipManager.MemberList.Members)
	n_targetMember := int(math.Sqrt(float64(n_members)))
	selectedMembers := node.MembershipManager.SelectNMembers(n_targetMember)

	// Send the message to the selected members
	for _, member := range selectedMembers {
		msg.Receipient = member.Address
		node.Tranceiver.Transmit(msg)
	}

	return nil
}
