package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/utils"
)

type Transaction struct {
	Sender    string  // Public key of the sender
	Recipient string  // Public key of the recipient
	Amount    float64 // Amount of coins transferred
	Timestamp int64   // Unix timestamp
	Signature string  // Sender's signature
}

// CalculateHash calculates the hash of a transaction
func (tx *Transaction) CalculateHash() string {
	// Concatenate the transaction data
	data := fmt.Sprintf("%s%s%f%d", tx.Sender, tx.Recipient, tx.Amount, tx.Timestamp)

	// Calculate the hash
	hash := utils.Hash(data)

	return hash
}

// Sign signs the transaction with the sender's private key
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	// Calculate the hash of the transaction
	hash := tx.CalculateHash()

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte(hash))
	if err != nil {
		return err
	}

	// Encode the signature
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = hex.EncodeToString(signature)
	return nil
}

// Verify verifies the signature of a transaction
func (tx *Transaction) Verify(publicKey *ecdsa.PublicKey) bool {
	// Decode the signature
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil || len(signature) != 64 {
		return false
	}

	// Split the signature into r and s
	r := big.Int{}
	s := big.Int{}
	r.SetBytes(signature[:32])
	s.SetBytes(signature[32:])

	// Calculate the hash of the transaction
	hash := tx.CalculateHash()

	// Verify the signature
	return ecdsa.Verify(publicKey, []byte(hash), &r, &s)
}
