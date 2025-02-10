package node

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
)

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
			node.handleNewTransactionMsg(msg)
		case message.NEWBLOCK:
			node.handleNewBlockMsg(msg)
		case message.BLOCKCHAINREQ:
			node.handleBlockChainRequest(msg)
		case message.BLOCKCHAINRESP:
			node.handleBlockchainResponse(msg)
		default:
			log.Printf("Unknown message type: %s\n", msg.Type)
		}
	}
}

func (node *Node) handleNewTransactionMsg(msg *message.Message) {
	// Gossip the transaction
	node.GossipManager.Gossip(msg)

	// Deserialize the transaction
	tx, err := transaction.DeserializeTransaction(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize transaction: %v\n", err)
		return
	}

	// Validate the transaction
	if err := tx.Validate(); err != nil {
		log.Printf("Invalid transaction: %v\n", err)
		return
	}

	// Add the transaction to the pool
	node.Mempool.AddTransaction(tx)
}

// handleNewBlockMsg handles a new block message
func (node *Node) handleNewBlockMsg(msg *message.Message) {
	// Gossip the block
	node.GossipManager.Gossip(msg)

	// Deserialize the block
	block, err := block.DeserializeBlock(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize block: %v\n", err)
		return
	}

	if err := node.Blockchain.ValidateBlock(block); err != nil {
		log.Printf("Invalid block: %s\n", err)
		msg := message.NewMessage(message.BLOCKCHAINREQ, node.IPAddress, msg.Sender, "")
		node.Transceiver.Transmit(msg)
	} else {
		node.Miner.Stop()
		node.Blockchain.AddBlock(block)
		go node.Miner.Run()
	}
}

// handleBlockChainRequest handles a blockchain request message
func (node *Node) handleBlockChainRequest(msg *message.Message) {
	payload, _ := node.Blockchain.Serialize()
	newMsg := message.NewMessage(message.BLOCKCHAINRESP, node.IPAddress, msg.Sender, payload)
	node.Transceiver.Transmit(newMsg)
}

// handleBlockchainResponse handles a blockchain response message
func (node *Node) handleBlockchainResponse(msg *message.Message) {
	blockchain, _ := blockchain.DeserializeBlockchain(msg.Payload)
	if err := node.Blockchain.ShouldSwitchChain(blockchain); err != nil {
		log.Printf("Invalid blockchain: %s\n", err)
	} else {
		node.Miner.Stop()
		node.Blockchain.SwitchChain(blockchain)
		go node.Miner.Run()
	}
}
