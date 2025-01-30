package mining

import (
	"log"
	"strings"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Miner struct {
	Address    string                 // Address of the miner
	Blockchain *blockchain.Blockchain // Reference to the blockchain
	Difficulty int                    // Difficulty of the mining process
	StopMining chan bool              // Channel to stop the mining process
}

// NewMiner creates a new miner with the given transactions, blockchain and difficulty
func NewMiner(address string, blockchain *blockchain.Blockchain) *Miner {
	return &Miner{
		Address:    address,
		Blockchain: blockchain,
		StopMining: make(chan bool),
	}
}

// StartMining starts the mining process with the given miner address
func (miner *Miner) Start(transactions []*transaction.Transaction, reward float64, difficulty int) {
	log.Println("Starting mining process...")

	// Skip mining process if the transaction pool is empty
	if len(transactions) == 0 {
		log.Println("Transactions is empty. Skipping mining process...")
		return
	}

	// Create a new block with the miner's address and reward
	newBlock := miner.Blockchain.NewBlock(transactions, miner.Address, reward, difficulty)

	// Set the StopMining channel to false
	miner.StopMining <- false

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
				return
			}
			block.Nonce++
		}
	}
}

// Stop terminates the mining process
func (miner *Miner) Stop() {
	log.Println("Stopping mining process...")
	miner.StopMining <- true
}

// Close closes the miner
func (miner *Miner) Close() {
	close(miner.StopMining)
}
