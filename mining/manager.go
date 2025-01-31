package mining

import (
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/blockchain"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/gossip"
)

type MiningManager struct {
	Miners         []*Miner               // List of Miners
	Blockchain     *blockchain.Blockchain // Blockchain reference
	GossipMananger *gossip.GossipManager  // Gossip manager reference
	MiningActive   bool                   // Mining active flag
	Mutex          *sync.Mutex            // Mutex for mining manager
}

// NewMiningManager creates a new mining manager
func NewMiningManager(blockchain *blockchain.Blockchain, gossipManager *gossip.GossipManager) *MiningManager {
	return &MiningManager{
		Miners:         make([]*Miner, 0),
		Blockchain:     blockchain,
		GossipMananger: gossipManager,
		MiningActive:   false,
		Mutex:          &sync.Mutex{},
	}
}
