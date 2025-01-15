package p2p

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type Member struct {
	Address   string
	Heartbeat int64
	Timestamp int64
}

// NewJOINREQMessage creates a new JOINREQ message
func NewJOINREQMessage(sender string) message.Message {
	return message.Message{
		Type:    "JOINREQ",
		Sender:  sender,
		Payload: "",
	}
}

// IntroduceSelfToGroup sends a JOINREQ message to the bootstrap node
func (node *Node) IntroduceSelfToGroup(bootstrapNodeAddress string) error {
	// Create a JOINREQ message and serialize it
	message := NewJOINREQMessage(node.Address)
	messageData := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessageData(bootstrapNodeAddress, messageData)
	if err != nil {
		return fmt.Errorf("failed to send JOINREQ message: %v", err)
	}

	return nil
}

// NewJOINRESPMessage creates a new JOINRESP message
func NewJOINRESPMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:    "JOINRESP",
		Sender:  sender,
		Payload: payload,
	}
}

// SerializeMemberList serializes a list of members into a string
func SerializeMemberList(memberList []Member) string {
	str := ""
	for _, member := range memberList {
		str += fmt.Sprintf("%s:%d:%d,", member.Address, member.Heartbeat, member.Timestamp)
	}
	return str
}

// HandleJoinRequest processes a JOINREQ message
func (node *Node) HandleJoinRequest(msg message.Message) {
	// Add the sender to the member list
	member := Member{
		Address:   msg.Sender,
		Heartbeat: 0,
		Timestamp: time.Now().Unix(),
	}
	node.MemberList = append(node.MemberList, member)

	// Send JOINREP message to the sender with the current member list
	payload := SerializeMemberList(node.MemberList)
	message := NewJOINRESPMessage(node.Address, payload)
	messageData := message.Serialize()
	err := network.SendMessageData(msg.Sender, messageData)
	if err != nil {
		log.Printf("Failed to send JOINRESP message: %v\n", err)
	}
}

// DeserializeMemberList deserializes a string into a list of members
func DeserializeMemberList(str string) ([]Member, error) {
	memberList := []Member{}
	members := strings.Split(str, ",")
	for _, memberStr := range members {
		if memberStr == "" {
			continue
		}
		var member Member
		_, err := fmt.Sscanf(memberStr, "%s:%d:%d", &member.Address, &member.Heartbeat, &member.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize member: %v", err)
		}
		memberList = append(memberList, member)
	}
	return memberList, nil
}

// HandleJoinResponse processes a JOINRESP message
func (node *Node) HandleJoinResponse(msg message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	node.MemberList = memberList
}

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:    "HEARTBEAT",
		Sender:  sender,
		Payload: payload,
	}
}

// SendHeartbeat sends a heartbeat message to all members in the network
func (node *Node) SendHeartbeat() {
	// Create a HEARTBEAT message and serialize it
	payload := SerializeMemberList(node.MemberList)
	message := NewJOINRESPMessage(node.Address, payload)
	messageData := message.Serialize()

	// Send HEARTBEAT message to all members in the network
	for _, member := range node.MemberList {
		if member.Address == node.Address {
			continue
		}
		err := network.SendMessageData(member.Address, messageData)
		if err != nil {
			log.Printf("Failed to send HEARTBEAT message: %v\n", err)
		}
	}
}

// HandleHeartbeat processes a HEARTBEAT message
func (node *Node) HandleHeartbeat(msg message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	node.UpdateMemberList(memberList)
}

// UpdateMemberList updates the member list with the new list of members
func (node *Node) UpdateMemberList(newMemberList []Member) {
}
