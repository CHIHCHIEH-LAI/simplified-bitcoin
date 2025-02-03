package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

// SaveToFile saves the wallet to a JSON file
func (w *Wallet) SaveToFile(filename string) error {
	// Convert private key to bytes
	privateKeyBytes, err := x509.MarshalECPrivateKey(w.PrivateKey)
	if err != nil {
		return err
	}

	// Create a struct to save both keys
	walletData := struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	}{
		PrivateKey: hex.EncodeToString(privateKeyBytes),
		PublicKey:  hex.EncodeToString(w.PublicKey),
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(walletData, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0600) // Secure file permissions
}

// LoadFromFile loads the wallet from a JSON file
func LoadFromFile(filename string) (*Wallet, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("wallet file not found")
	}

	// Read file content
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Deserialize JSON
	var walletData struct {
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	}
	if err := json.Unmarshal(data, &walletData); err != nil {
		return nil, err
	}

	// Decode private key
	privateKeyBytes, err := hex.DecodeString(walletData.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Decode public key
	publicKeyBytes, err := hex.DecodeString(walletData.PublicKey)
	if err != nil {
		return nil, err
	}

	// Return the wallet
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKeyBytes,
	}, nil
}
