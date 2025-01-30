package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/gossip"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Node struct {
	Address           string                        // IP address of the node
	Port              string                        // Port of the node
	Transceiver       *network.Transceiver          // Tranceiver instance
	MembershipManager *membership.MembershipManager // Membership manager
	GossipManager     *gossip.GossipManager         // Gossip manager
	Mempool           *transaction.Mempool          // Transaction manager
}

// NewNode creates a new P2P node
func NewNode(address, port string) (*Node, error) {
	var err error

	// Create a new tranceiver
	transceiver, err := network.NewTransceiver(port)
	if err != nil {
		return nil, err
	}

	// Create a new membership manager
	membershipManager := membership.NewMembershipManager(address, transceiver)

	return &Node{
		Address:           address,
		Port:              port,
		Transceiver:       transceiver,
		MembershipManager: membershipManager,
		GossipManager:     gossip.NewGossipManager(transceiver, membershipManager),
		Mempool:           transaction.NewMempool(),
	}, nil
}

// Run starts the P2P node
func (node *Node) Run(bootstrapNodeAddr string) error {
	// Run the tranceiver
	go node.Transceiver.Run()

	// Handle incoming messages
	go node.handleIncomingMessage()

	// Run the membership manager
	go node.MembershipManager.Run(bootstrapNodeAddr)

	// Start the gossip manager
	go node.GossipManager.Run(60)

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
			node.GossipManager.Gossip(msg)
			node.Mempool.HandleNewTransaction(msg)
		// case message.NEWBLOCK:
		// 	node.HandleNewBlock(msg)
		// case "GETBLOCKCHAIN":
		// 	node.HandleGetBlockchain(msg)
		default:
			log.Printf("Unknown message type: %s\n", msg.Type)
		}
	}
}

// Close closes the P2P node
func (node *Node) Close() {
	// Close the tranceiver
	node.Transceiver.Close()

	// Close the gossip manager
	node.GossipManager.Close()
}
