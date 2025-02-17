package node

import (
	"log"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/p2p/membership"
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
			node.handleJoinRequest(msg)
		case message.JOINRESP:
			node.handleJoinResponse(msg)
		case message.HEARTBEAT:
			node.handleHeartbeatMsg(msg)
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

// handleJoinRequest handles a JOINREQ message
func (node *Node) handleJoinRequest(msg *message.Message) {
	requester := msg.Sender
	node.MembershipManager.HandleJoinRequest(requester)
}

// handleJoinResponse handles a JOINRESP message
func (node *Node) handleJoinResponse(msg *message.Message) {
	// Deserialize the member list from the payload
	memberList, err := membership.DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	node.MembershipManager.HandleJoinResponse(memberList)
}

// handleHeartbeatMsg handles a heartbeat message
func (node *Node) handleHeartbeatMsg(msg *message.Message) {
	// Deserialize the member list from the payload
	memberList, err := membership.DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	node.MembershipManager.HandleHeartbeat(memberList)
}

// handleNewTransactionMsg handles a new transaction message
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
	if err := node.Blockchain.ValidateTransaction(tx); err != nil {
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

	if err := node.Blockchain.ValidateNewBlock(block); err != nil {
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
	payload, err := node.Blockchain.Serialize()
	if err != nil {
		log.Printf("Failed to serialize blockchain: %v\n", err)
		return
	}

	newMsg := message.NewMessage(message.BLOCKCHAINRESP, node.IPAddress, msg.Sender, payload)
	node.Transceiver.Transmit(newMsg)
}

// handleBlockchainResponse handles a blockchain response message
func (node *Node) handleBlockchainResponse(msg *message.Message) {
	blockchain, err := blockchain.DeserializeBlockchain(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize blockchain: %v\n", err)
		return
	}

	if err := node.Blockchain.ShouldSwitchChain(blockchain); err != nil {
		log.Printf("Invalid blockchain: %s\n", err)
	} else {
		log.Printf("Switching to a new blockchain\n")
		node.Miner.Stop()
		node.Blockchain.SwitchChain(blockchain)
		time.Sleep(5 * time.Second)
		go node.Miner.Run()
	}
}
