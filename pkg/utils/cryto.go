package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// VerifySignature checks if the given signature is valid for the given data
func VerifySignature(publicKey, data, signature string) bool {
	pubKeyBytes, _ := hex.DecodeString(publicKey)
	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])
	pubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	signatureBytes, _ := hex.DecodeString(signature)
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])
	hash := sha256.Sum256([]byte(data))
	return ecdsa.Verify(&pubKey, hash[:], r, s)
}
