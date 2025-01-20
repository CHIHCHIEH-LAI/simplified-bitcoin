package membership

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
	"golang.org/x/exp/rand"
)

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:      message.HEARTBEAT,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// HandleHeartbeat processes a HEARTBEAT message
func (mgr *MembershipManager) HandleHeartbeat(msg message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	mgr.UpdateMemberList(memberList)
}

// SendHeartbeat sends a heartbeat message to some random members in the network
func (mgr *MembershipManager) SendHeartbeat() {
	// Skip if there is only one member in the network
	if len(mgr.MemberList) == 1 {
		return
	}

	// Create a HEARTBEAT message and serialize it
	payload := SerializeMemberList(mgr.MemberList)
	message := NewHEARTBEATMessage(mgr.Address, payload)
	messageData := message.Serialize()

	// Send HEARTBEAT message to some random members in the network
	selectedMembers := make(map[int]bool)
	limit := min(NUMMEMBERSTOHEARTBEAT, len(mgr.MemberList))

	for len(selectedMembers) < limit {
		index := rand.Intn(len(mgr.MemberList))

		// Skip self
		if mgr.MemberList[index].Address == mgr.Address {
			continue
		}

		// Skip if the member is already selected
		if _, ok := selectedMembers[index]; ok {
			continue
		}

		// Send HEARTBEAT message to the member
		err := network.SendMessageData(mgr.MemberList[index].Address, messageData)
		if err != nil {
			log.Printf("Failed to send HEARTBEAT message: %v\n", err)
		}

		selectedMembers[index] = true
	}
}
