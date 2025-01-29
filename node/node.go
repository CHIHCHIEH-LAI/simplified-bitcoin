package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Node struct {
	Address            string                          // IP address of the node
	Port               string                          // Port of the node
	Transceiver        *network.Transceiver            // Tranceiver instance
	MembershipManager  *membership.MembershipManager   // Membership manager
	TransactionManager *transaction.TransactionManager // Transaction manager
}

// NewNode creates a new P2P node
func NewNode(address, port string) (*Node, error) {
	var err error

	// Create a new tranceiver
	transceiver, err := network.NewTransceiver(port)
	if err != nil {
		return nil, err
	}

	return &Node{
		Address:            address,
		Port:               port,
		Transceiver:        transceiver,
		MembershipManager:  membership.NewMembershipManager(address),
		TransactionManager: transaction.NewTransactionManager(),
	}, nil
}

// Run starts the P2P node
func (node *Node) Run(bootstrapNodeAddress string) error {
	// Run the tranceiver
	go node.Transceiver.Run()

	// Start processing messages
	go node.HandleIncomingMessage()

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

// Close closes the P2P node
func (node *Node) Close() {
	// Close the tranceiver
	node.Transceiver.Close()
}
