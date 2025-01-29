package membership

import (
	"time"

	"golang.org/x/exp/rand"
)

type MembershipManager struct {
	Address    string
	MemberList *MemberList
}

// NewMembershipManager creates a new membership manager
func NewMembershipManager(address string) *MembershipManager {
	return &MembershipManager{
		Address:    address,
		MemberList: NewMemberList(),
	}
}

// MaintainMembership maintains the membership list by sending heartbeats
func (mgr *MembershipManager) MaintainMembership() {
	for {
		mgr.MemberList.UpdateSelfInMemberList(mgr.Address)
		mgr.MemberList.RemoveFailedMembers()
		mgr.GossipHeartbeat()
		time.Sleep(TIMEHEARTBEAT * time.Second)
	}
}

// SelectMembers selects n_member random members from the member list
func (mgr *MembershipManager) SelectNMembers(n_target int) []string {
	selectedMembers := make(map[int]bool)
	limit := min(n_target, len(mgr.MemberList.Members)-1)

	for len(selectedMembers) < limit {
		index := rand.Intn(len(mgr.MemberList.Members))

		// Skip self
		if mgr.MemberList.Members[index].Address == mgr.Address {
			continue
		}

		// Skip if the member is already selected
		if _, ok := selectedMembers[index]; ok {
			continue
		}

		selectedMembers[index] = true
	}

	// Convert the selected members to a slice
	selectedMembersSlice := make([]string, 0)
	for index := range selectedMembers {
		selectedMembersSlice = append(selectedMembersSlice, mgr.MemberList.Members[index].Address)
	}

	return selectedMembersSlice
}
