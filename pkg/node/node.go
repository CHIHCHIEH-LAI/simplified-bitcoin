package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
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
		case message.NEWBLOCK:
			node.GossipManager.Gossip(msg)
			block, _ := block.DeserializeBlock(msg.Payload)
			if err := node.Blockchain.ValidateBlock(block); err != nil {
				log.Printf("Invalid block: %s\n", err)
				node.AskForBlockchain(msg.Sender)
			} else {
				node.Miner.Stop()
				node.Blockchain.AddBlock(block)
				go node.Miner.Run()
			}
		case message.BLOCKCHAINREQ:
			sender := msg.Sender
			node.ShareBlockchain(sender)
		case message.BLOCKCHAINRESP:
			blockchain, _ := blockchain.DeserializeBlockchain(msg.Payload)
			if err := node.Blockchain.ShouldSwitchChain(blockchain); err != nil {
				log.Printf("Invalid blockchain: %s\n", err)
			} else {
				node.Miner.Stop()
				node.Blockchain.SwitchChain(blockchain)
				go node.Miner.Run()
			}
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

func (node *Node) AskForBlockchain(IPAddress string) {
	msg := message.NewMessage(message.BLOCKCHAINREQ, node.IPAddress, IPAddress, "")
	node.Transceiver.Transmit(msg)
}

func (node *Node) ShareBlockchain(IPAddress string) {
	payload, _ := node.Blockchain.Serialize()
	msg := message.NewMessage(message.BLOCKCHAINRESP, node.IPAddress, IPAddress, payload)
	node.Transceiver.Transmit(msg)
}
