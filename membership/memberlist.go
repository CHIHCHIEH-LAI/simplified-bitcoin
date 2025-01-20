package membership

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

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
