package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/gossip"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/mempool"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/mining"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type Node struct {
	IPAddress         string                        // IP address of the node
	Port              string                        // Port of the node
	Address           string                        // Address of the node (Hash of the public key)
	Transceiver       *network.Transceiver          // Tranceiver instance
	MembershipManager *membership.MembershipManager // Membership manager
	GossipManager     *gossip.GossipManager         // Gossip manager
	Mempool           *mempool.Mempool              // Mempool
	Blockchain        *blockchain.Blockchain        // Blockchain
	Miner             *mining.Miner                 // Miner
}

// NewNode creates a new P2P node
func NewNode(IPAddress, port, address string) (*Node, error) {
	var err error

	// Create a new tranceiver
	transceiver, err := network.NewTransceiver(port)
	if err != nil {
		return nil, err
	}

	// Create a new membership manager
	membershipManager := membership.NewMembershipManager(IPAddress, transceiver)

	// Create a Blockchain
	blockchain := blockchain.NewBlockchain()

	// Create a Gossip Manager
	gossipManager := gossip.NewGossipManager(IPAddress, transceiver, membershipManager)

	// Create a Mempool
	mempool := mempool.NewMempool()

	// Create a Miner
	miner := mining.NewMiner(address, blockchain, gossipManager, mempool)

	return &Node{
		Port:              port,
		Transceiver:       transceiver,
		MembershipManager: membershipManager,
		GossipManager:     gossipManager,
		Mempool:           mempool,
		Blockchain:        blockchain,
		Miner:             miner,
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

	// Run the gossip manager
	go node.GossipManager.Run(60)

	// Run the miner
	go node.Miner.Run()

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
			requester := msg.Sender
			node.MembershipManager.HandleJoinRequest(requester)
		case message.JOINRESP:
			node.MembershipManager.HandleJoinResponse(msg)
		case message.HEARTBEAT:
			node.MembershipManager.HandleHeartbeat(msg)
		case message.NEWTRANSACTION:
			node.GossipManager.Gossip(msg)
			node.Mempool.HandleNewTransaction(msg)
		case message.NEWBLOCK:
			node.GossipManager.Gossip(msg)
			node.Miner.Restart()
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

	// Close the miner
	node.Miner.Close()
}
