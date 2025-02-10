package blockchain

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/block"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/blockchain/transaction"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/mining/mempool"
)

type Blockchain struct {
	BaseReward    float64
	Blocks        []*block.Block   `json:"blocks"` // Blocks in the blockchain
	mutex         *sync.RWMutex    // Mutex to protect the blockchain
	CumulativePoW int              `json:"cumulativePoW"` // Tracks total proof-of-work (sum of difficulties)
	Mempool       *mempool.Mempool // Reference to the mempool
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain(mempool *mempool.Mempool) *Blockchain {
	genesisBlock := block.NewGenesisBlock()
	return &Blockchain{
		BaseReward:    1000.0,
		Blocks:        []*block.Block{genesisBlock},
		mutex:         &sync.RWMutex{},
		CumulativePoW: genesisBlock.Difficulty,
		Mempool:       mempool,
	}
}

// NewBlock creates a new block with the given transactions
func (bc *Blockchain) NewBlock(transactions []*transaction.Transaction, miner string) *block.Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	prevHash := bc.GetLatestBlock().BlockID
	reward := bc.CalculateReward()
	difficulty := bc.CalculateDifficulty()
	return block.NewBlock(prevHash, transactions, miner, reward, difficulty)
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *block.Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Validate the block
	if err := bc.ValidateBlock(block); err != nil {
		return err
	}

	bc.Blocks = append(bc.Blocks, block)
	bc.CumulativePoW += block.Difficulty

	// Remove transactions in the block from the mempool
	for _, tx := range block.Transactions {
		bc.Mempool.RemoveTransaction(tx.TransactionID)
	}

	return nil
}

// GetLatestBlock returns the latest block in the blockchain
func (bc *Blockchain) GetLatestBlock() *block.Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// CalculateReward calculates the reward for the miner
func (bc *Blockchain) CalculateReward() float64 {
	return bc.BaseReward / float64(len(bc.Blocks))
}

// CalculateDifficulty calculates the difficulty for the miner
func (bc *Blockchain) CalculateDifficulty() int {
	return 7
}

// CalculateCumulativePoW calculates the cumulative proof-of-work
func (bc *Blockchain) CalculateCumulativePoW() int {
	cumulativePoW := 0
	for _, b := range bc.Blocks {
		cumulativePoW += b.Difficulty
	}
	return cumulativePoW
}

// calculateUTXOs calculates the UTXOs for an address
func (bc *Blockchain) calculateUTXOs(address string) float64 {
	utxos := 0.0
	for _, b := range bc.Blocks {
		for _, tx := range b.Transactions {
			if tx.Sender == address {
				utxos -= tx.Amount + tx.Fee
			}
			if tx.Recipient == address {
				utxos += tx.Amount
			}
		}
	}
	return utxos
}

// Serialize serializes the blockchain to a JSON string
func (bc *Blockchain) Serialize() (string, error) {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	data, err := json.Marshal(bc)
	if err != nil {
		return "", fmt.Errorf("failed to serialize block: %v", err)
	}
	return string(data), nil
}

// DeserializeBlockchain deserializes a JSON string to a blockchain
func DeserializeBlockchain(data string) (*Blockchain, error) {
	var bc Blockchain
	err := json.Unmarshal([]byte(data), &bc)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize blockchain: %v", err)
	}
	return &bc, nil
}
