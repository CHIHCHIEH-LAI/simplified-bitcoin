package mining

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/gossip"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/mempool"
)

type Miner struct {
	NTransactions int                    // Number of transactions per block
	Address       string                 // Wallet Address of the Miner
	Blockchain    *blockchain.Blockchain // Blockchain reference
	GossipManager *gossip.GossipManager  // Gossip manager reference
	Mempool       *mempool.Mempool       // Mempool reference
	StopMining    chan bool              // Channel to stop mining
	MiningActive  bool                   // Indicates if mining is active
	Mutex         *sync.Mutex            // Ensures safe concurrent mining operations
}

// NewMiner creates a new miner
func NewMiner(
	address string,
	blockchain *blockchain.Blockchain,
	gossipManager *gossip.GossipManager,
	mempool *mempool.Mempool,
) *Miner {
	return &Miner{
		NTransactions: 10,
		Address:       address,
		Blockchain:    blockchain,
		GossipManager: gossipManager,
		Mempool:       mempool,
		StopMining:    make(chan bool, 1),
		Mutex:         &sync.Mutex{},
	}
}

// Run starts the mining loop
func (miner *Miner) Run() {
	miner.Mutex.Lock()
	defer miner.Mutex.Unlock()

	// Prevent duplicate mining sessions
	if miner.MiningActive {
		log.Println("Mining is already running.")
		return
	}
	miner.MiningActive = true
	log.Println("Starting mining process...")

	for {
		select {
		case <-miner.StopMining:
			log.Println("Mining stopped.")
			miner.MiningActive = false
			return
		default:
			// Get the top N rewarding transactions from the mempool
			transactions := miner.Mempool.GetTopNRewardingTransactions(miner.NTransactions)
			if len(transactions) == 0 {
				log.Println("No transactions available. Pausing mining...")
				time.Sleep(2 * time.Second) // Prevents high CPU usage when waiting for transactions
				continue
			}

			// Create a new block
			newBlock := miner.Blockchain.NewBlock(transactions, miner.Address)

			// Perform Proof of Work
			minedBlock := miner.PerformProofOfWork(newBlock)
			if minedBlock == nil {
				log.Println("Mining was interrupted.")
				return
			}

			// Add Mined Block to Blockchain
			miner.Blockchain.AddBlock(minedBlock)

			// Broadcast the Mined Block
			miner.BroadcastBlock(minedBlock)

			// Pause to allow network sync before restarting
			time.Sleep(2 * time.Second)
		}
	}
}

// PerformProofOfWork executes the proof of work algorithm
func (miner *Miner) PerformProofOfWork(block *blockchain.Block) *blockchain.Block {
	log.Printf("Mining block %s with difficulty %d...\n", block.BlockID, block.Difficulty)
	prefix := strings.Repeat("0", block.Difficulty)

	block.Nonce = 0
	for {
		select {
		case <-miner.StopMining:
			log.Println("Mining interrupted due to a new block.")
			return nil
		default:
			blockHash := block.Hash()
			if strings.HasPrefix(blockHash, prefix) {
				block.BlockID = blockHash
				log.Printf("Block mined: %s\n", block.BlockID)
				return block
			}
			block.Nonce++
		}
	}
}

// BroadcastBlock sends the newly mined block to the network
func (miner *Miner) BroadcastBlock(block *blockchain.Block) {
	msg := blockchain.NewMinedBlockMessage(block)
	miner.GossipManager.Gossip(msg)
	log.Printf("Broadcasted new block: %s", block.BlockID)
}

// StopMining stops the current mining process
func (miner *Miner) Stop() {
	miner.Mutex.Lock()
	defer miner.Mutex.Unlock()

	if !miner.MiningActive {
		log.Println("Mining is already stopped.")
		return
	}

	log.Println("Stopping mining process...")
	miner.StopMining <- true
	miner.MiningActive = false
}

// RestartMining stops mining and starts it again after a short delay
func (miner *Miner) Restart() {
	miner.Stop()
	time.Sleep(2 * time.Second) // Allow other nodes to sync before restarting
	go miner.Run()
}

// Close closes the miner
func (miner *Miner) Close() {
	miner.Stop()
	close(miner.StopMining)
}
