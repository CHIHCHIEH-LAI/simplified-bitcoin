package blockchain

type Blockchain struct {
	Blocks              []Block
	PendingTransactions []Transaction
	Difficulty          int
	MiningReward        float64
}
