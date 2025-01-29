package node

import (
	"log"
	"math"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
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
		MembershipManager:  membership.NewMembershipManager(address, transceiver),
		TransactionManager: transaction.NewTransactionManager(),
	}, nil
}

// Run starts the P2P node
func (node *Node) Run(bootstrapNodeAddress string) error {
	// Run the tranceiver
	go node.Transceiver.Run()

	// Handle incoming messages
	go node.handleIncomingMessage()

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

// HandleIncomingMessage processes incoming messages
func (node *Node) handleIncomingMessage() {
	for {
		msg, ok := node.Transceiver.Receive()
		if !ok {
			continue // Skip iteration if no message
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
			node.gossipMessage(msg)
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
func (node *Node) gossipMessage(msg *message.Message) error {

	// Select N random members to send the message to
	n_members := len(node.MembershipManager.MemberList.Members)
	n_targetMember := int(math.Sqrt(float64(n_members)))
	selectedMembers := node.MembershipManager.SelectNMembers(n_targetMember)

	// Send the message to the selected members
	for _, member := range selectedMembers {
		msg.Receipient = member.Address
		node.Transceiver.Transmit(msg)
	}

	return nil
}

// Close closes the P2P node
func (node *Node) Close() {
	// Close the tranceiver
	node.Transceiver.Close()
}
