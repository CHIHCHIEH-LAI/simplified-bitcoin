package membership

type Member struct {
	Address   string `json:"address"`
	Heartbeat int64  `json:"heartbeat"`
	Timestamp int64  `json:"timestamp"`
}

type MemberList struct {
	Members []Member `json:"members"`
}
