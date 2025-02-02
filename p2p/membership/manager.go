package membership

import (
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
	"golang.org/x/exp/rand"
)

type MembershipManager struct {
	IPAddress   string
	MemberList  *MemberList
	Transceiver *network.Transceiver
}

// NewMembershipManager creates a new membership manager
func NewMembershipManager(IPAddress string, transceiver *network.Transceiver) *MembershipManager {
	return &MembershipManager{
		IPAddress:   IPAddress,
		MemberList:  NewMemberList(),
		Transceiver: transceiver,
	}
}

// Run starts the membership manager
func (mgr *MembershipManager) Run(bootstrapNodeAddr string) {
	// Join the p2p network
	mgr.JoinGroup(bootstrapNodeAddr)

	for {
		mgr.MemberList.UpdateSelfInMemberList(mgr.IPAddress)
		mgr.MemberList.RemoveFailedMembers()
		mgr.GossipHeartbeat()
		time.Sleep(TIMEHEARTBEAT * time.Second)
	}
}

func (mgr *MembershipManager) GetNumberOfMembers() int {
	return len(mgr.MemberList.Members)
}

// SelectMembers selects n_member random members from the member list
func (mgr *MembershipManager) SelectNMembers(n_target int) []*Member {
	selectedMembers := make(map[int]bool)
	limit := min(n_target, len(mgr.MemberList.Members)-1)

	for len(selectedMembers) < limit {
		index := rand.Intn(len(mgr.MemberList.Members))

		// Skip self
		if mgr.MemberList.Members[index].IPAddress == mgr.IPAddress {
			continue
		}

		// Skip if the member is already selected
		if _, ok := selectedMembers[index]; ok {
			continue
		}

		selectedMembers[index] = true
	}

	// Convert the selected members to a slice
	selectedMembersSlice := make([]*Member, 0)
	for index := range selectedMembers {
		selectedMembersSlice = append(selectedMembersSlice, mgr.MemberList.Members[index])
	}

	return selectedMembersSlice
}
