package mining

import (
	"log"
	"strings"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Miner struct {
	Address      string                     // Address of the miner
	Transactions []*transaction.Transaction // Reference to the transaction pool
	Blockchain   *blockchain.Blockchain     // Reference to the blockchain
	Difficulty   int                        // Difficulty of the mining process
	StopMining   chan bool                  // Channel to stop the mining process
}

// NewMiner creates a new miner with the given transactions, blockchain and difficulty
func NewMiner(address string, transactions []*transaction.Transaction, blockchain *blockchain.Blockchain, difficulty int) *Miner {
	return &Miner{
		Address:      address,
		Transactions: transactions,
		Blockchain:   blockchain,
		Difficulty:   difficulty,
		StopMining:   make(chan bool),
	}
}

// StartMining starts the mining process with the given miner address
func (miner *Miner) Start() {
	log.Println("Starting mining process...")

	// Skip mining process if the transaction pool is empty
	if len(miner.Transactions) == 0 {
		log.Println("Transactions is empty. Skipping mining process...")
		return
	}

	// Create a new block with the miner's address and reward
	newBlock := miner.Blockchain.NewBlock(miner.Transactions, miner.Address, REWARD)

	// Perform the proof of work algorithm
	miner.PerformProofOfWork(newBlock)
}

// PerformProofOfWork performs the proof of work algorithm
func (miner *Miner) PerformProofOfWork(block *blockchain.Block) {
	log.Printf("Mining block %s with difficulty %d...\n", block.BlockID, miner.Difficulty)

	prefix := strings.Repeat("0", miner.Difficulty)
	for {
		select {
		case <-miner.StopMining:
			// Stop the mining process
			log.Println("Mining process terminated.")
			return
		default:
			// Continue the mining process
			blockHash := block.Hash()
			if strings.HasPrefix(blockHash, prefix) {
				log.Printf("Block mined: %s\n", blockHash)
				block.BlockID = blockHash
				miner.Blockchain.AddBlock(block)
				miner.BroadcastBlock(block)
				return
			}
			block.Nonce++
		}
	}
}

// TODO: Implement the BroadcastBlock function
// BroadcastBlock broadcasts a block to the network
func (miner *Miner) BroadcastBlock(block *blockchain.Block) {
}

// Stop terminates the mining process
func (miner *Miner) Stop() {
	log.Println("Stopping mining process...")
	miner.StopMining <- true
	close(miner.StopMining)
}
