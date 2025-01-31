package mining

import (
	"log"
	"strings"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/gossip"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/mempool"
)

type Miner struct {
	NTransactions  int                    // Number of transactions to include in the block
	IPAddress      string                 // IP Address of the Mining Manager
	Address        string                 // Address of the Mining Manager
	Blockchain     *blockchain.Blockchain // Blockchain reference
	GossipMananger *gossip.GossipManager  // Gossip manager reference
	Mempool        *mempool.Mempool       // Mempool reference
	StopMining     chan bool              // Channel to stop the mining process
}

// NewMiner creates a new miner with the given transactions, blockchain and difficulty
func NewMiner(
	IPAddress string,
	address string,
	blockchain *blockchain.Blockchain,
	gossipManager *gossip.GossipManager,
	mempool *mempool.Mempool,
) *Miner {
	return &Miner{
		NTransactions:  10,
		IPAddress:      IPAddress,
		Address:        address,
		Blockchain:     blockchain,
		GossipMananger: gossipManager,
		Mempool:        mempool,
	}
}

// Run runs the mining process
func (miner *Miner) Run() {

	// Get the top N rewarding transactions from the mempool
	transactions := miner.Mempool.GetTopNRewardingTransactions(miner.NTransactions)

	// Create a new block with the transactions
	newBlock := miner.Blockchain.NewBlock(transactions, miner.Address)

	// Perform the proof of work algorithm
	minedBlock := miner.PerformProofOfWork(newBlock)

	// If the block was mined, add it to the blockchain
	if minedBlock != nil {
		miner.Blockchain.AddBlock(minedBlock)
		miner.BroadcastBlock(minedBlock)
	}

	// Run the mining process again
	go miner.Run()
}

// PerformProofOfWork performs the proof of work algorithm
func (miner *Miner) PerformProofOfWork(block *blockchain.Block) *blockchain.Block {

	prefix := strings.Repeat("0", block.Difficulty)
	for {
		select {
		case <-miner.StopMining:
			return nil
		default:
			// Continue the mining process
			blockHash := block.Hash()
			if strings.HasPrefix(blockHash, prefix) {
				log.Printf("Block mined: %s\n", blockHash)
				block.BlockID = blockHash
				return block
			}
			block.Nonce++
		}
	}
}

// BroadcastBlock broadcasts the mined block to the network
func (miner *Miner) BroadcastBlock(block *blockchain.Block) {
	msg := blockchain.NewMinedBlockMessage(miner.IPAddress, block)
	miner.GossipMananger.Gossip(msg)
}

// Stop terminates the mining process
func (miner *Miner) Stop() {
	miner.StopMining <- true
}

// Restart restarts the mining process
func (miner *Miner) Restart() {
	miner.Stop()
	time.Sleep(10 * time.Second)
	go miner.Run()
}

// Close closes the miner
func (miner *Miner) Close() {
	close(miner.StopMining)
}
