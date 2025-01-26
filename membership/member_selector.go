package membership

import (
	"golang.org/x/exp/rand"
)

// SelectMembers selects n_member random members from the member list
func (mgr *MembershipManager) SelectNMembers(n_target int) []string {
	selectedMembers := make(map[int]bool)
	limit := min(n_target, len(mgr.MemberList)-1)

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

		selectedMembers[index] = true
	}

	// Convert the selected members to a slice
	selectedMembersSlice := make([]string, 0)
	for index := range selectedMembers {
		selectedMembersSlice = append(selectedMembersSlice, mgr.MemberList[index].Address)
	}

	return selectedMembersSlice
}
