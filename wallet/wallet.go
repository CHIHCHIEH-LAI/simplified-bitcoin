package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/network"
)

type Wallet struct {
	PrivateKey  *ecdsa.PrivateKey `json:"-"`
	PublicKey   []byte            `json:"public_key"`
	Transmitter *network.Transmitter
}

// NewWallet creates and returns a Wallet
func NewWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	messageChannel := make(chan *message.Message)
	transmitter := network.NewTransmitter(messageChannel)
	return &Wallet{privateKey, publicKey, transmitter}
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
