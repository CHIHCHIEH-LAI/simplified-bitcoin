package blockchain

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	Signature string
}
