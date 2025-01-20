package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Node struct {
	Address           string                        // IP address of the node
	MembershipManager *membership.MembershipManager // Membership manager
	MessageChannel    chan string                   // Channel to send and receive messages
	TransactionPool   *transaction.TransactionPool  // Pool of transactions
}

// NewNode creates a new P2P node
func NewNode(address string) *Node {
	return &Node{
		Address:           address,
		MembershipManager: membership.NewMembershipManager(address),
		MessageChannel:    make(chan string, 100),
		TransactionPool:   transaction.NewTransactionPool(),
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

	// Start processing messages
	go node.HandleMessage()

	// Join the p2p network
	err := node.MembershipManager.JoinGroup(bootstrapNodeAddress)
	if err != nil {
		log.Printf("Failed to join network via bootstrap node %s: %v", bootstrapNodeAddress, err)
		return err
	}

	// Start maintaining membership
	go node.MembershipManager.MaintainMembership()

	return nil
}
