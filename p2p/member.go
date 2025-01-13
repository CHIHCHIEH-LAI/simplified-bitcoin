package p2p

type Member struct {
	BitcoinAddress string
	Address        string
	Heartbeat      int64
	Timestamp      int64
}
