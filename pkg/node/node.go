package node

import (
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/mining"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/mining/mempool"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/p2p/gossip"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/p2p/membership"
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

	// Create a Mempool
	mempool := mempool.NewMempool()

	// Create a Blockchain
	blockchain := blockchain.NewBlockchain(mempool)

	// Create a Gossip Manager
	gossipManager := gossip.NewGossipManager(IPAddress, transceiver, membershipManager)

	// Create a Miner
	miner := mining.NewMiner(address, blockchain, gossipManager, mempool)

	return &Node{
		IPAddress:         IPAddress,
		Port:              port,
		Address:           address,
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

	// Run the blockchain
	go node.Blockchain.Run()

	return nil
}

// Close closes the P2P node
func (node *Node) Close() {
	// Close the tranceiver
	node.Transceiver.Close()

	// Close the gossip manager
	node.GossipManager.Close()

	// Close the miner
	node.Miner.Close()

	// Close the blockchain
	node.Blockchain.Close()
}
