package membership

import "time"

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
		mgr.SendHeartbeat()
		time.Sleep(TIMEHEARTBEAT * time.Second)
	}
}
