package p2p

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type Node struct {
	Address           string      // IP address of the node
	MemberList        []Member    // List of members in the network
	MemberListSelfPos int         // Position of the node in the member list
	MessageChannel    chan string // Channel to send and receive messages
}

// NewNode creates a new P2P node
func NewNode(address string) *Node {
	return &Node{
		Address:           address,
		MemberList:        []Member{},
		MemberListSelfPos: -1,
		MessageChannel:    make(chan string, 100),
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
	// node.HandleMessage()

	return nil
}
