package blockchain

import "errors"

type Blockchain struct {
	Blocks              []Block
	PendingTransactions []Transaction
	Difficulty          int
	MiningReward        float64
}

func NewBlockchain() Blockchain {
	return Blockchain{
		Blocks:              []Block{NewGenesisBlock()},
		PendingTransactions: []Transaction{},
		Difficulty:          2,
		MiningReward:        50.0,
	}
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, amount float64, timestamp string, signature string, publicKey string) error {
	tx := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Timestamp: timestamp,
		Signature: signature,
	}

	if !tx.IsValid(publicKey) {
		return errors.New("invalid transaction")
	}

	bc.PendingTransactions = append(bc.PendingTransactions, tx)
	return nil
}
