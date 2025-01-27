package mining

import (
	"log"
	"strings"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Miner struct {
	Transactions *[]*transaction.Transaction // Reference to the transaction pool
	Blockchain   *blockchain.Blockchain      // Reference to the blockchain
	Difficulty   int                         // Difficulty of the mining process
	StopMining   chan bool                   // Channel to stop the mining process
	MutexPool    *sync.Mutex                 // Mutex to lock the transaction pool
}

// NewMiner creates a new miner with the given transaction pool and difficulty
func NewMiner(transactions *[]*transaction.Transaction, difficulty int) *Miner {
	return &Miner{
		Transactions: transactions,
		Difficulty:   difficulty,
		StopMining:   make(chan bool),
		MutexPool:    &sync.Mutex{},
	}
}

// StartMining starts the mining process with the given miner address
func (miner *Miner) StartMining(minerAddress string) {
	log.Println("Starting mining process...")

	// Get transactions
	transactions := miner.GetTransactions()

	if len(transactions) == 0 {
		log.Println("Transaction pool is empty. Skipping mining process...")
		return
	}

	// Create a coinbase transaction to reward the miner
	coinbaseTx := transaction.NewCoinbaseTransaction(minerAddress, REWARD)
	transactions = append([]*transaction.Transaction{coinbaseTx}, transactions...)

	// Get the lastest block from the blockchain
	latestBlock := miner.Blockchain.GetLatestBlock()

	newBlock := blockchain.NewBlock(latestBlock.BlockID, transactions)

	// Perform the proof of work algorithm
	miner.PerformProofOfWork(newBlock)

}

// GetTransactionsFromPool gets transactions
func (miner *Miner) GetTransactions() []*transaction.Transaction {
	miner.MutexPool.Lock()
	defer miner.MutexPool.Unlock()

	transactions := miner.Transactions

	return *transactions
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
		}
		block.Nonce++
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
