package membership

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
	"golang.org/x/exp/rand"
)

const TIMEHEARTBEAT = 60
const TIMENODEFAIL = 10 * TIMEHEARTBEAT
const TIMENODEREMOVE = 5 * TIMENODEFAIL
const NUMMEMBERSTOHEARTBEAT = 10

type Member struct {
	Address   string
	Heartbeat int64
	Timestamp int64
}

type MembershipManager struct {
	Address    string
	MemberList []Member
}

// NewMembershipManager creates a new membership manager
func NewMembershipManager(address string) *MembershipManager {
	return &MembershipManager{
		Address:    address,
		MemberList: []Member{},
	}
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

// NewHEARTBEATMessage creates a new HEARTBEAT message
func NewHEARTBEATMessage(sender string, payload string) message.Message {
	return message.Message{
		Type:      message.HEARTBEAT,
		Sender:    sender,
		Payload:   payload,
		Timestamp: utils.GetCurrentTimeInUnix(),
	}
}

// SerializeMemberList serializes a list of members into a string
func SerializeMemberList(memberList []Member) string {
	str := ""
	for _, member := range memberList {
		str += fmt.Sprintf("%s;%d;%d,", member.Address, member.Heartbeat, member.Timestamp)
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
		parts := strings.Split(memberStr, ";")
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
func (mgr *MembershipManager) HandleJoinRequest(msg message.Message) {
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
	messageData := message.Serialize()
	err := network.SendMessageData(msg.Sender, messageData)
	if err != nil {
		log.Printf("Failed to send JOINRESP message: %v\n", err)
	}
}

// HandleJoinResponse processes a JOINRESP message
func (mgr *MembershipManager) HandleJoinResponse(msg message.Message) {
	// Deserialize the member list from the payload
	memberList, err := DeserializeMemberList(msg.Payload)
	if err != nil {
		log.Printf("Failed to deserialize member list: %v\n", err)
		return
	}

	// Update the member list
	mgr.MemberList = memberList
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

// UpdateMemberList updates the member list with the new list of members
func (mgr *MembershipManager) UpdateMemberList(newMemberList []Member) {
	for _, newMember := range newMemberList {
		// Check if the member is self
		if newMember.Address == mgr.Address {
			continue
		}

		// Find the member in the list
		index := mgr.FindMemberInList(newMember.Address)
		if index == -1 { // Add the member if it does not exist
			mgr.AddMemberToList(newMember)
		} else if newMember.Heartbeat > mgr.MemberList[index].Heartbeat { // Update the member if the heartbeat is greater
			mgr.UpdateMemberInList(index, newMember)
		}
	}
}

// FindMemberInList finds a member in the list by address
func (mgr *MembershipManager) FindMemberInList(address string) int {
	for i, member := range mgr.MemberList {
		if member.Address == address {
			return i
		}
	}
	return -1
}

// AddMemberToList adds a member to the list
func (mgr *MembershipManager) AddMemberToList(member Member) {
	mgr.MemberList = append(mgr.MemberList, member)
}

// UpdateMemberInList updates a member in the list
func (mgr *MembershipManager) UpdateMemberInList(index int, newMember Member) {
	// Check if the newMember is failed
	if utils.GetCurrentTimeInUnix()-newMember.Timestamp > TIMENODEFAIL {
		return
	}

	// Update the member in the list
	mgr.MemberList[index].Heartbeat = newMember.Heartbeat
	mgr.MemberList[index].Timestamp = utils.GetCurrentTimeInUnix()
}

// JoinGroup joins the P2P group via the bootstrap node
func (mgr *MembershipManager) JoinGroup(bootstrapNodeAddress string) error {
	// Introduce self to the group if bootstrap node is self
	if bootstrapNodeAddress == "" || bootstrapNodeAddress == mgr.Address {
		mgr.IntroduceSelfToGroup()
		return nil
	}

	// Create a JOINREQ message and serialize it
	message := NewJOINREQMessage(mgr.Address)
	messageData := message.Serialize()

	// Send JOINREQ message to bootstrap node
	err := network.SendMessageData(bootstrapNodeAddress, messageData)
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

// MaintainMembership maintains the membership list by sending heartbeats
func (mgr *MembershipManager) MaintainMembership() {
	for {
		mgr.UpdateSelfInMemberList()
		mgr.RemoveFailedNodes()
		mgr.SendHeartbeat()
		time.Sleep(TIMEHEARTBEAT * time.Second)
	}
}

// UpdateSelfInMemberList updates the self member in the member list
func (mgr *MembershipManager) UpdateSelfInMemberList() {
	index := mgr.FindMemberInList(mgr.Address)
	if index == -1 {
		member := Member{
			Address:   mgr.Address,
			Heartbeat: 0,
			Timestamp: utils.GetCurrentTimeInUnix(),
		}
		mgr.AddMemberToList(member)
	} else {
		mgr.MemberList[index].Heartbeat++
		mgr.MemberList[index].Timestamp = utils.GetCurrentTimeInUnix()
	}
}

// RemoveFailedNodes removes failed nodes from the member list
func (mgr *MembershipManager) RemoveFailedNodes() {
	for i, member := range mgr.MemberList {
		if utils.GetCurrentTimeInUnix()-member.Timestamp > TIMENODEREMOVE {
			mgr.MemberList = append(mgr.MemberList[:i], mgr.MemberList[i+1:]...)
		}
	}
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
