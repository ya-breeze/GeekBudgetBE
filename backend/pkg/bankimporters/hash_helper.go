package bankimporters

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

// ComputeStableHash computes a stable identifier based on transaction semantics:
// Date, Amounts, and Currencies.
func ComputeStableHash(t *goserver.TransactionNoId) string {
	if len(t.Movements) == 0 {
		return ""
	}

	// Use the main movement amount (usually index 0)
	amount := t.Movements[0].Amount
	currency := t.Movements[0].CurrencyId
	date := t.Date.Unix()

	data := fmt.Sprintf("%d|%.2f|%s", date, amount, currency)
	return HashString(data)
}

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
