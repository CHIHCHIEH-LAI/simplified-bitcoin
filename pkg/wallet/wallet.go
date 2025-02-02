package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"

	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/message"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/network"
	"github.com/CHIHCHIEH-LAI/simplified-bitcoin/pkg/utils"
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
	return utils.Hash(string(w.PublicKey))
}

// Sign creates a signature for the given data using the wallet's private key
func (w *Wallet) Sign(hash string) (string, error) {
	// Decode the hash
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hashBytes)
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return hex.EncodeToString(signature), nil
}
