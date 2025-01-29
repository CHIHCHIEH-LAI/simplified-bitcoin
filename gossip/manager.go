package gossip

import (
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
