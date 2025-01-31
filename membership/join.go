package membership

import (
	"log"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

// JoinGroup joins the P2P group via the bootstrap node
func (mgr *MembershipManager) JoinGroup(bootstrapNodeAddress string) {
	// Introduce self to the group if bootstrap node is self
	if bootstrapNodeAddress == "" || bootstrapNodeAddress == mgr.IPAddress {
		mgr.IntroduceSelfToGroup()
		return
	}

	// Create a JOINREQ message
	message := NewJOINREQMessage(mgr.IPAddress, bootstrapNodeAddress)

	// Send JOINREQ message
	mgr.Transceiver.Transmit(message)
}

// IntroduceSelfToGroup sends a JOINREQ message to the bootstrap node
func (mgr *MembershipManager) IntroduceSelfToGroup() {
	// Add self to the member list
	member := &Member{
		IPAddress: mgr.IPAddress,
		Heartbeat: 0,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
	mgr.MemberList.AddMemberToList(member)
}

// NewJOINREQMessage creates a new JOINREQ message
func NewJOINREQMessage(selfAddr, bootstrapAddr string) *message.Message {
	return &message.Message{
		Type:       message.JOINREQ,
		Sender:     selfAddr,
		Receipient: bootstrapAddr,
		Payload:    "",
		Timestamp:  utils.GetCurrentTimeInUnix(),
	}
}

// NewJOINRESPMessage creates a new JOINRESP message
func NewJOINRESPMessage(selfAddr, receipient, payload string) *message.Message {
	return &message.Message{
		Type:       message.JOINRESP,
		Sender:     selfAddr,
		Receipient: receipient,
		Payload:    payload,
		Timestamp:  utils.GetCurrentTimeInUnix(),
	}
}

// HandleJoinRequest processes a JOINREQ message
func (mgr *MembershipManager) HandleJoinRequest(requester string) {
	member := &Member{
		IPAddress: requester,
		Heartbeat: 0,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}

	// Check if the sender is already in the member list
	if index := mgr.MemberList.FindMemberInList(requester); index == -1 {
		mgr.MemberList.AddMemberToList(member)
	} else {
		mgr.MemberList.UpdateMemberInList(index, member)
	}

	// Send JOINREP message to the sender with the current member list
	payload, err := mgr.MemberList.Serialize()
	if err != nil {
		log.Printf("Failed to serialize member list: %v\n", err)
		return
	}

	msg := NewJOINRESPMessage(mgr.IPAddress, requester, payload)

	// Send JOINRESP message
	mgr.Transceiver.Transmit(msg)
}

// HandleJoinResponse processes a JOINRESP message
func (mgr *MembershipManager) HandleJoinResponse(msg *message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	mgr.MemberList = memberList
}
