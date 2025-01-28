package membership

import (
	"encoding/json"
	"fmt"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Member struct {
	Address   string `json:"address"`
	Heartbeat int64  `json:"heartbeat"`
	Timestamp int64  `json:"timestamp"`
}

type MemberList struct {
	Members []*Member `json:"members"`
}

func NewMemberList() *MemberList {
	return &MemberList{
		Members: []*Member{},
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
	for _, newMember := range newMemberList.Members {
		// Check if the member is self
		if newMember.Address == selfAddr {
			continue
		}

		// Find the member in the list
		index := ml.FindMemberInList(newMember.Address)
		if index == -1 { // Add the member if it does not exist
			ml.AddMemberToList(newMember)
		} else if newMember.Heartbeat > ml.Members[index].Heartbeat { // Update the member if the heartbeat is greater
			ml.UpdateMemberInList(index, newMember)
		}
	}
}

// FindMemberInList finds a member in the list by address
func (ml *MemberList) FindMemberInList(address string) int {
	for i, member := range ml.Members {
		if member.Address == address {
			return i
		}
	}
	return -1
}

// AddMemberToList adds a member to the list
func (ml *MemberList) AddMemberToList(member *Member) {
	ml.Members = append(ml.Members, member)
}

// UpdateMemberInList updates a member in the list
func (ml *MemberList) UpdateMemberInList(index int, newMember *Member) {
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
	index := ml.FindMemberInList(selfAddr)
	if index == -1 {
		member := &Member{
			Address:   selfAddr,
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
	for i, member := range ml.Members {
		if utils.GetCurrentTimeInUnix()-member.Timestamp > TIMENODEREMOVE {
			ml.Members = append(ml.Members[:i], ml.Members[i+1:]...)
		}
	}
}
