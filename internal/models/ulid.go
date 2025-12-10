package models

import (
	"crypto/rand"
	"math/big"
	"time"
)

// ULID generates a Universally Unique Lexicographically Sortable Identifier
// Format: 26 characters (10 timestamp + 16 random)
func GenerateULID() string {
	// Crockford's Base32 alphabet (case-insensitive, excludes I, L, O, U)
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// Generate timestamp part (10 chars, milliseconds since epoch)
	timestamp := time.Now().UnixMilli()
	timestampChars := make([]byte, 10)
	for i := 9; i >= 0; i-- {
		timestampChars[i] = alphabet[timestamp%32]
		timestamp /= 32
	}

	// Generate random part (16 chars)
	randomChars := make([]byte, 16)
	for i := 0; i < 16; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(32))
		randomChars[i] = alphabet[n.Int64()]
	}

	// Combine to 25 chars (truncate first char for compatibility)
	result := make([]byte, 25)
	copy(result[0:9], timestampChars[1:])
	copy(result[9:], randomChars)

	return string(result)
}
