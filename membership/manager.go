package membership

import "time"

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

// MaintainMembership maintains the membership list by sending heartbeats
func (mgr *MembershipManager) MaintainMembership() {
	for {
		mgr.UpdateSelfInMemberList()
		mgr.RemoveFailedNodes()
		mgr.SendHeartbeat()
		time.Sleep(TIMEHEARTBEAT * time.Second)
	}
}
