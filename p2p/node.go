package p2p

import "fmt"

type P2PNode struct {
	BitcoinAddress    string   // Bitcoin address of the node
	Address           string   // IP address of the node
	MemberList        []Member // List of members in the network
	MemberListSelfPos int      // Position of the node in the member list
}

// NewP2PNode creates a new P2P node
func NewP2PNode(bitcoinAddress, address string) *P2PNode {
	return &P2PNode{
		BitcoinAddress:    bitcoinAddress,
		Address:           address,
		MemberList:        []Member{},
		MemberListSelfPos: -1,
	}
}

// Start starts the P2P node
func (node *P2PNode) Start(bootstrapNodeAddress string) error {
	err := node.IntroduceSelfToGroup(bootstrapNodeAddress)
	if err != nil {
		return fmt.Errorf("failed to introduce self to group: %v", err)
	}
	return nil
}
