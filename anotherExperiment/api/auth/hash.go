package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters - tuned for API key hashing
// These values balance security with performance
const (
	Argon2Time    = 1         // Number of iterations
	Argon2Memory  = 64 * 1024 // Memory in KB (64 MB)
	Argon2Threads = 4         // Number of threads
	Argon2KeyLen  = 32        // Output hash length
	SaltLength    = 16        // Salt length in bytes
)

// KeyHasher handles secure hashing of API keys
type KeyHasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func (h *KeyHasher) Verify(apiKey string, hashedKey string) (bool, error) {

	parts := strings.Split(hashedKey, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Parse parameters
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, fmt.Errorf("failed to parse hash params: %w", err)
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Decode stored hash to get key length
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Recompute hash with same params
	computedHash := argon2.IDKey([]byte(apiKey), salt, time, memory, threads, uint32(len(storedHash)))

	return subtle.ConstantTimeCompare(computedHash, storedHash) == 1, nil
}

// NewKeyHasher creates a hasher with recommended parameters
func NewKeyHasher() *KeyHasher {
	return &KeyHasher{
		time:    Argon2Time,
		memory:  Argon2Memory,
		threads: Argon2Threads,
		keyLen:  Argon2KeyLen,
	}
}

// Hash creates a secure hash of the API key
// Returns a string in format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
func (h *KeyHasher) Hash(key string) (string, error) {
	// Generate a random salt
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Compute the Argon2id hash
	hash := argon2.IDKey(
		[]byte(key),
		salt,
		h.time,
		h.memory,
		h.threads,
		h.keyLen,
	)

	// Encode salt and hash to base64
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=time,p=threads$salt$hash
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.memory,
		h.time,
		h.threads,
		saltB64,
		hashB64,
	)

	return encoded, nil
}
