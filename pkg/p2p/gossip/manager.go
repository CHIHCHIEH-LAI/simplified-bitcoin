package gossip

import (
	"fmt"
	"sync"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/p2p/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
)

type GossipManager struct {
	IPAddress         string                        // IP address of the node
	Transceiver       *network.Transceiver          // Tranceiver instance
	MembershipManager *membership.MembershipManager // Membership manager
	SeenMessage       map[string]bool               // Seen messages
	Mutex             *sync.Mutex                   // Mutex to protect the seen messages
}

// NewGossipManager creates a new gossip manager
func NewGossipManager(IPAddress string, transceiver *network.Transceiver, membershipManager *membership.MembershipManager) *GossipManager {
	return &GossipManager{
		IPAddress:         IPAddress,
		Transceiver:       transceiver,
		MembershipManager: membershipManager,
		SeenMessage:       make(map[string]bool),
		Mutex:             &sync.Mutex{},
	}
}

// Run starts the gossip manager
func (mgr *GossipManager) Run(staleThreshold int64) {
	mgr.cleanSeenMessages()

	time.Sleep(10 * 60 * time.Second)
}

// Gossip sends a message to N random members
func (mgr *GossipManager) Gossip(msg *message.Message) {
	mgr.Mutex.Lock()
	defer mgr.Mutex.Unlock()

	// Check if the message has been seen before
	if _, ok := mgr.SeenMessage[hashMessage(msg)]; ok {
		return
	}

	// Set the sender of the message
	msg.Sender = mgr.IPAddress

	// Select N random members to send the message to
	n_members := mgr.MembershipManager.GetNumberOfMembers()
	n_targetMembers := utils.ISqrt(n_members)
	selectedMembers := mgr.MembershipManager.SelectNMembers(n_targetMembers)

	// Send the message to the selected members
	for _, member := range selectedMembers {
		msg.Receipient = member.IPAddress
		mgr.Transceiver.Transmit(msg)
	}

	// Mark the message as seen
	mgr.SeenMessage[hashMessage(msg)] = true
}

// cleanSeenMessages cleans the seen messages
func (mgr *GossipManager) cleanSeenMessages() {
	mgr.Mutex.Lock()
	defer mgr.Mutex.Unlock()

	mgr.SeenMessage = make(map[string]bool)
}

// Close closes the gossip manager
func (mgr *GossipManager) Close() {
}

// hashMessage hashes the message
func hashMessage(msg *message.Message) string {
	msgData := fmt.Sprintf("%s|%s", msg.Type, msg.Payload)
	return utils.Hash(msgData)
}
