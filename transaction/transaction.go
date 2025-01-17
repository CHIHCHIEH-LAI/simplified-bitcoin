package transaction

type Transaction struct {
	TransactionID string
	Sender        string
	Recipient     string
	Amount        float64
	Fee           float64
	Timestamp     int64
	Signature     string
}
