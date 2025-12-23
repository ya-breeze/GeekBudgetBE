package bankimporters

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	// Create a new SHA-256 hash object
	hasher := sha256.New()

	// Write the input string to the hash object
	hasher.Write([]byte(input))

	// Get the resulting hash as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
