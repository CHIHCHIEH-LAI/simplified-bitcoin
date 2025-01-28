package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

// HandleMessage processes incoming messages
func (node *Node) HandleMessage() {
	for msgData := range node.MessageChannel {
		log.Printf("Received message: %s\n", msgData)

		// Parse the message
		msg, err := message.DeserializeMessage(msgData)
		if err != nil {
			log.Printf("Failed to deserialize message: %v\n", err)
			continue
		}

		// Process the message based on its type
		switch msg.Type {
		case message.JOINREQ:
			node.MembershipManager.HandleJoinRequest(msg)
		case message.JOINRESP:
			node.MembershipManager.HandleJoinResponse(msg)
		case message.HEARTBEAT:
			node.MembershipManager.HandleHeartbeat(msg)
		case message.NEWTRANSACTION:
			selectedMembers := node.MembershipManager.SelectNMembers(transaction.NUM_MEMBERS_TO_BROADCAST)
			node.TransactionManager.HandleNewTransaction(msg, selectedMembers)
		// case message.NEWBLOCK:
		// 	node.HandleNewBlock(msg)
		// case "GETBLOCKCHAIN":
		// 	node.HandleGetBlockchain(msg)
		default:
			log.Printf("Unknown message type: %s\n", msg.Type)
		}
	}
}
