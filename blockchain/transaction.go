package blockchain

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	Timestamp string
	Signature string
}

func (tx *Transaction) IsValid(publicKey string) bool {
	return true
}
