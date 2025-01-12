package blockchain

type Mempool struct {
	Transactions []Transaction
}

// NewMempool creates a new mempool
func NewMempool() Mempool {
	return Mempool{
		Transactions: []Transaction{},
	}
}

// AddTransaction adds a transaction to the mempool
func (m *Mempool) AddTransaction(tx Transaction) {
	m.Transactions = append(m.Transactions, tx)
}

// GetTransactions returns all transactions in the mempool
func (m *Mempool) GetTransactions() []Transaction {
	return m.Transactions
}

// Clear removes all transactions from the mempool
func (m *Mempool) Clear() {
	m.Transactions = []Transaction{}
}
