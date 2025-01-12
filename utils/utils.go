package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Hash returns the SHA256 hash of the input data
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SerializeHashes serializes a slice of hashes into a JSON array
func SerializeHashes(hashes []string) ([]byte, error) {
	data, err := json.Marshal(hashes)
	if err != nil {
		return nil, err
	}
	return data, nil
}
