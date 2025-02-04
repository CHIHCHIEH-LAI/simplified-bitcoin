package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// VerifySignature checks if the given signature is valid for the given data
func VerifySignature(publicKey, data, signature string) error {
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil || len(pubKeyBytes) < 64 {
		return fmt.Errorf("invalid public key")
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil || len(signatureBytes) < 64 {
		return fmt.Errorf("invalid signature")
	}

	// Extract X and Y coordinates for public key
	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])

	// Extract R and S values for the signature
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Hash the message
	hash := sha256.Sum256([]byte(data))

	// Verify the signature
	pubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	if ecdsa.Verify(&pubKey, hash[:], r, s) {
		return nil
	} else {
		return fmt.Errorf("invalid signature")
	}
}
