package p2p

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type Node struct {
	Address        string      // IP address of the node
	MemberList     []Member    // List of members in the network
	MessageChannel chan string // Channel to send and receive messages
}

// NewNode creates a new P2P node
func NewNode(address string) *Node {
	return &Node{
		Address:        address,
		MemberList:     []Member{},
		MessageChannel: make(chan string, 100),
	}
}

// Run starts the P2P node
func (node *Node) Run(port string, bootstrapNodeAddress string) error {
	// Run the communication listener
	go func() {
		err := network.RunListener(port, node.MessageChannel)
		if err != nil {
			log.Fatalf("Failed to start listener on port %s: %v", port, err)
		}
	}()
	log.Printf("Node started on port %s\n", port)

	// Introduce self to the p2p group via the bootstrap node
	err := node.IntroduceSelfToGroup(bootstrapNodeAddress)
	if err != nil {
		log.Printf("Failed to join network via bootstrap node %s: %v", bootstrapNodeAddress, err)
	}

	// Start processing messages
	go node.HandleMessage()

	// Start maintaining membership
	go node.MaintainMembership()

	return nil
}

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
		switch {
		case msg.Type == "JOINREQ":
			node.HandleJoinRequest(msg)
		case msg.Type == "JOINRESP":
			node.HandleJoinResponse(msg)
		case msg.Type == "HEARTBEAT":
			node.HandleHeartbeat(msg)
		// case msg.Type == "NEWBLOCK":
		// 	node.HandleNewBlock(msg)
		// case msg.Type == "NEWTRANSACTION":
		// 	node.HandleNewTransaction(msg)
		// case msg.Type == "GETBLOCKCHAIN":
		// 	node.HandleGetBlockchain(msg)
		default:
			log.Printf("Unknown message type: %s\n", msg)
		}
	}
}
