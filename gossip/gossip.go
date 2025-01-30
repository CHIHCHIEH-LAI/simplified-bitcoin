package gossip

import (
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/membership"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type GossipManager struct {
	Transceiver       *network.Transceiver          // Tranceiver instance
	MembershipManager *membership.MembershipManager // Membership manager
	SeenMessage       map[message.Message]bool      // Seen messages
}

// NewGossipManager creates a new gossip manager
func NewGossipManager(transceiver *network.Transceiver, membershipManager *membership.MembershipManager) *GossipManager {
	return &GossipManager{
		Transceiver:       transceiver,
		MembershipManager: membershipManager,
		SeenMessage:       make(map[message.Message]bool),
	}
}

// Run starts the gossip manager
func (mgr *GossipManager) Run(staleThreshold int64) {
	mgr.removeStaleMessages(staleThreshold)

	time.Sleep(5 * time.Second)
}

// Gossip sends a message to N random members
func (mgr *GossipManager) Gossip(msg *message.Message) {
	// Check if the message has been seen before
	if _, ok := mgr.SeenMessage[*msg]; ok {
		return
	}

	// Select N random members to send the message to
	n_members := mgr.MembershipManager.GetNumberOfMembers()
	n_targetMembers := utils.ISqrt(n_members)
	selectedMembers := mgr.MembershipManager.SelectNMembers(n_targetMembers)

	// Send the message to the selected members
	for _, member := range selectedMembers {
		msg.Receipient = member.Address
		mgr.Transceiver.Transmit(msg)
	}

	// Mark the message as seen
	mgr.SeenMessage[*msg] = true
}

// removeStaleMessages removes messages that are older than the given threshold
func (mgr *GossipManager) removeStaleMessages(threshold int64) {
	for msg := range mgr.SeenMessage {
		if msg.Timestamp < threshold {
			delete(mgr.SeenMessage, msg)
		}
	}
}

// Close closes the gossip manager
func (mgr *GossipManager) Close() {
}
