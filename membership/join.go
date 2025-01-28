package membership

import (
	"fmt"
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

// JoinGroup joins the P2P group via the bootstrap node
func (mgr *MembershipManager) JoinGroup(bootstrapNodeAddress string) error {
	// Introduce self to the group if bootstrap node is self
	if bootstrapNodeAddress == "" || bootstrapNodeAddress == mgr.Address {
		mgr.IntroduceSelfToGroup()
		return nil
	}

	// Create a JOINREQ message and serialize it
	message := NewJOINREQMessage(mgr.Address)
	messageData, err := message.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize JOINREQ message: %v", err)
	}

	// Send JOINREQ message to bootstrap node
	err = network.SendMessageData(bootstrapNodeAddress, messageData)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}

// IntroduceSelfToGroup sends a JOINREQ message to the bootstrap node
func (mgr *MembershipManager) IntroduceSelfToGroup() {
	// Add self to the member list
	member := Member{
		Address:   mgr.Address,
		Heartbeat: 0,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
	mgr.MemberList = append(mgr.MemberList, member)
}

// NewJOINREQMessage creates a new JOINREQ message
func NewJOINREQMessage(sender string) message.Message {
	return message.Message{
		Type:      message.JOINREQ,
		Sender:    sender,
		Payload:   "",
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// NewJOINRESPMessage creates a new JOINRESP message
func NewJOINRESPMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:      message.JOINRESP,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// HandleJoinRequest processes a JOINREQ message
func (mgr *MembershipManager) HandleJoinRequest(msg *message.Message) {
	// Check if the sender is already in the member list
	if index := mgr.FindMemberInList(msg.Sender); index == -1 {
		// Add the sender to the member list
		member := Member{
			Address:   msg.Sender,
			Heartbeat: 0,
			Timestamp: utils.GetCurrentTimeInUnix(),
		}
		mgr.MemberList = append(mgr.MemberList, member)
	} else {
		// Update the sender's heartbeat
		mgr.MemberList[index].Heartbeat = 0
		mgr.MemberList[index].Timestamp = utils.GetCurrentTimeInUnix()
	}

	// Send JOINREP message to the sender with the current member list
	payload := SerializeMemberList(mgr.MemberList)
	message := NewJOINRESPMessage(mgr.Address, payload)
	messageData, err := message.Serialize()
	if err != nil {
		log.Printf("Failed to serialize JOINRESP message: %v\n", err)
		return
	}

	err = network.SendMessageData(msg.Sender, messageData)
	if err != nil {
		log.Printf("Failed to send JOINRESP message: %v\n", err)
	}
}

// HandleJoinResponse processes a JOINRESP message
func (mgr *MembershipManager) HandleJoinResponse(msg *message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	mgr.MemberList = memberList
}
