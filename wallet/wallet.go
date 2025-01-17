package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/transaction"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{privateKey, publicKey}
}

// GetAddress generates a public key hash (address) for the wallet
func (w *Wallet) GetAddress() string {
	pubHash := sha256.Sum256(w.PublicKey)
	return hex.EncodeToString(pubHash[:])
}

// Sign creates a signature for the given data using the wallet's private key
func (w *Wallet) Sign(data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}

// CreateTransaction creates a new transaction
func (w *Wallet) CreateTransaction(recipient string, amount float64, fee float64) (*transaction.Transaction, error) {
	// Create the transaction
	tx := transaction.Transaction{
		Sender:    w.GetAddress(),
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
	}

	// Generate the transaction ID
	tx.TransactionID = tx.GenerateTransactionID()

	// Sign the transaction
	data := fmt.Sprintf("%s%s%f%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Fee, tx.Timestamp)
	signature, err := w.Sign(data)
	if err != nil {
		return nil, err
	}
	tx.Signature = signature

	return &tx, nil
}

// SendTransaction sends a transaction to a node in the network
func SendTransaction(tx *transaction.Transaction, nodeAddress string) error {
	// Send the transaction to the node
	return nil
}
