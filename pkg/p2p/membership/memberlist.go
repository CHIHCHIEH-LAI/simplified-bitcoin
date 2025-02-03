package membership

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
)

type Member struct {
	IPAddress string `json:"address"`   // IP:Port
	Heartbeat int64  `json:"heartbeat"` // Number of heartbeats
	Timestamp int64  `json:"timestamp"` // Timestamp of the last heartbeat
}

type MemberList struct {
	Members []*Member   `json:"members"` // List of members
	Mutex   *sync.Mutex // Mutex to protect the member list
}

// NewMemberList creates a new member list
func NewMemberList() *MemberList {
	return &MemberList{
		Members: []*Member{},
		Mutex:   &sync.Mutex{},
	}
}

// Serialize serializes the member list into a string
func (ml *MemberList) Serialize() (string, error) {
	data, err := json.Marshal(ml)
	if err != nil {
		return "", fmt.Errorf("failed to serialize memberlist: %v", err)
	}
	return string(data), nil
}

// DeserializeMemberList deserializes the member list from a string
func DeserializeMemberList(data string) (*MemberList, error) {
	var ml MemberList
	err := json.Unmarshal([]byte(data), &ml)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %v", err)
	}
	return &ml, nil
}

// UpdateMemberList updates the member list with the new list of members
func (ml *MemberList) UpdateMemberList(newMemberList *MemberList, selfAddr string) {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	for _, newMember := range newMemberList.Members {
		// Check if the member is self
		if newMember.IPAddress == selfAddr {
			continue
		}

		// Find the member in the list
		index := ml.FindMemberInList(newMember.IPAddress)
		if index == -1 { // Add the member if it does not exist
			ml.AddMemberToList(newMember)
		} else if newMember.Heartbeat > ml.Members[index].Heartbeat { // Update the member if the heartbeat is greater
			ml.UpdateMemberInList(index, newMember)
		}
	}
}

// FindMemberInList finds a member in the list by address
func (ml *MemberList) FindMemberInList(address string) int {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	for i, member := range ml.Members {
		if member.IPAddress == address {
			return i
		}
	}
	return -1
}

// AddMemberToList adds a member to the list
func (ml *MemberList) AddMemberToList(member *Member) {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	ml.Members = append(ml.Members, member)
}

// UpdateMemberInList updates a member in the list
func (ml *MemberList) UpdateMemberInList(index int, newMember *Member) {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	// Check if the newMember is failed
	if utils.GetCurrentTimeInUnix()-newMember.Timestamp > TIMENODEFAIL {
		return
	}

	// Update the member in the list
	ml.Members[index].Heartbeat = newMember.Heartbeat
	ml.Members[index].Timestamp = utils.GetCurrentTimeInUnix()
}

// UpdateSelfInMemberList updates the self member in the member list
func (ml *MemberList) UpdateSelfInMemberList(selfAddr string) {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	index := ml.FindMemberInList(selfAddr)
	if index == -1 {
		member := &Member{
			IPAddress: selfAddr,
			Heartbeat: 0,
			Timestamp: utils.GetCurrentTimeInUnix(),
		}
		ml.AddMemberToList(member)
	} else {
		ml.Members[index].Heartbeat++
		ml.Members[index].Timestamp = utils.GetCurrentTimeInUnix()
	}
}

// RemoveFailedNodes removes failed nodes from the member list
func (ml *MemberList) RemoveFailedMembers() {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()

	for i, member := range ml.Members {
		if utils.GetCurrentTimeInUnix()-member.Timestamp > TIMENODEREMOVE {
			ml.Members = append(ml.Members[:i], ml.Members[i+1:]...)
		}
	}
}
