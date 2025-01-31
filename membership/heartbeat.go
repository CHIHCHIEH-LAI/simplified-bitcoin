package membership

import (
	"log"
	"math"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) *message.Message {
	return &message.Message{
		Type:      message.HEARTBEAT,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// HandleHeartbeat processes a HEARTBEAT message
func (mgr *MembershipManager) HandleHeartbeat(msg *message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	mgr.MemberList.UpdateMemberList(memberList, mgr.IPAddress)
}

// GossipHeartbeat sends a HEARTBEAT message to some random members in the network
func (mgr *MembershipManager) GossipHeartbeat() {
	// Skip if there is only one member in the network
	if len(mgr.MemberList.Members) == 1 {
		return
	}

	// Create a HEARTBEAT message and serialize it
	payload, err := mgr.MemberList.Serialize()
	if err != nil {
		log.Printf("Failed to serialize member list: %v\n", err)
		return
	}

	message := NewHEARTBEATMessage(mgr.IPAddress, payload)

	// Select some random members to send the HEARTBEAT message
	n_target := int(math.Sqrt(float64(len(mgr.MemberList.Members))))
	selectedMembers := mgr.SelectNMembers(n_target)

	// Send HEARTBEAT message to some random members in the network
	for _, member := range selectedMembers {
		// Send HEARTBEAT message to the member
		message.Receipient = member.IPAddress
		mgr.Transceiver.Transmit(message)
	}
}
