package p2p

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

const TIMENODEFAIL = 1 * 60
const TIMENODEREMOVE = 4 * TIMENODEFAIL

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

// NewJOINRESPMessage creates a new JOINRESP message
func NewJOINRESPMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:    "JOINRESP",
		Sender:  sender,
		Payload: payload,
	}
}

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:    "HEARTBEAT",
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

// DeserializeMemberList deserializes a string into a list of members
func DeserializeMemberList(str string) ([]Member, error) {
	memberList := []Member{}
	members := strings.Split(str, ",")
	for _, memberStr := range members {
		if memberStr == "" {
			continue
		}
		var member Member
		parts := strings.Split(memberStr, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid member format")
		}
		member.Address = parts[0]
		member.Heartbeat, _ = strconv.ParseInt(parts[1], 10, 64)
		member.Timestamp, _ = strconv.ParseInt(parts[2], 10, 64)
		memberList = append(memberList, member)
	}
	return memberList, nil
}

// HandleJoinRequest processes a JOINREQ message
func (node *Node) HandleJoinRequest(msg message.Message) {
	// Check if the sender is already in the member list
	if index := node.FindMemberInList(msg.Sender); index != -1 {
		return
	}

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
	for _, newMember := range newMemberList {
		// Check if the member is self
		if newMember.Address == node.Address {
			continue
		}

		// Find the member in the list
		index := node.FindMemberInList(newMember.Address)
		if index == -1 { // Add the member if it does not exist
			node.AddMemberToList(newMember)
		} else if newMember.Heartbeat > node.MemberList[index].Heartbeat { // Update the member if the heartbeat is greater
			node.UpdateMemberInList(index, newMember)
		}
	}
}

// FindMemberInList finds a member in the list by address
func (node *Node) FindMemberInList(address string) int {
	for i, member := range node.MemberList {
		if member.Address == address {
			return i
		}
	}
	return -1
}

// AddMemberToList adds a member to the list
func (node *Node) AddMemberToList(member Member) {
	node.MemberList = append(node.MemberList, member)
}

// UpdateMemberInList updates a member in the list
func (node *Node) UpdateMemberInList(index int, newMember Member) {
	// Check if the newMember is failed
	if time.Now().Unix()-newMember.Timestamp > TIMENODEFAIL {
		return
	}

	// Update the member in the list
	node.MemberList[index].Heartbeat = newMember.Heartbeat
	node.MemberList[index].Timestamp = time.Now().Unix()
}

// JoinGroup joins the P2P group via the bootstrap node
func (node *Node) JoinGroup(bootstrapNodeAddress string) error {
	// Introduce self to the group if bootstrap node is self
	if bootstrapNodeAddress == node.Address {
		node.IntroduceSelfToGroup()
		return nil
	}

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

// IntroduceSelfToGroup sends a JOINREQ message to the bootstrap node
func (node *Node) IntroduceSelfToGroup() {
	// Add self to the member list
	member := Member{
		Address:   node.Address,
		Heartbeat: 0,
		Timestamp: time.Now().Unix(),
	}
	node.MemberList = append(node.MemberList, member)
}

// MaintainMembership maintains the membership list by sending heartbeats
func (node *Node) MaintainMembership() {
	for {
		node.UpdateSelfInMemberList()
		node.RemoveFailedNodes()
		node.SendHeartbeat()
		time.Sleep(5 * time.Second)
	}
}

// UpdateSelfInMemberList updates the self member in the member list
func (node *Node) UpdateSelfInMemberList() {
	index := node.FindMemberInList(node.Address)
	if index == -1 {
		member := Member{
			Address:   node.Address,
			Heartbeat: 0,
			Timestamp: time.Now().Unix(),
		}
		node.AddMemberToList(member)
	} else {
		node.MemberList[index].Heartbeat++
		node.MemberList[index].Timestamp = time.Now().Unix()
	}
}

// RemoveFailedNodes removes failed nodes from the member list
func (node *Node) RemoveFailedNodes() {
	for i, member := range node.MemberList {
		if time.Now().Unix()-member.Timestamp > TIMENODEREMOVE {
			node.MemberList = append(node.MemberList[:i], node.MemberList[i+1:]...)
		}
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
		// Skip self
		if member.Address == node.Address {
			continue
		}
		err := network.SendMessageData(member.Address, messageData)
		if err != nil {
			log.Printf("Failed to send HEARTBEAT message: %v\n", err)
		}
	}
}
