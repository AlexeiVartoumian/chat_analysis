package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	// KeyLength defines the number of random bytes to generate
	// 32 bytes = 256 bits of entropy, sufficient for security
	KeyLength = 32

	// KeyPrefix helps identify the key type and version
	// Format: {prefix}_{version}_{randomPart}
	KeyPrefix  = "ouk" // OneUptime Key
	KeyVersion = "v1"
)

// APIKeyGenerator handles secure API key generation
type APIKeyGenerator struct {
	prefix  string
	version string
}

// NewAPIKeyGenerator creates a new generator with default settings
func NewAPIKeyGenerator() *APIKeyGenerator {
	return &APIKeyGenerator{
		prefix:  KeyPrefix,
		version: KeyVersion,
	}
}

// Generate creates a new cryptographically secure API key
// Returns the full key (for the user) and the key ID (for lookups)
func (g *APIKeyGenerator) Generate() (fullKey string, keyID string, err error) {
	// Generate random bytes using crypto/rand for cryptographic security
	randomBytes := make([]byte, KeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to URL-safe base64 for easy transmission
	randomPart := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Generate a short key ID for lookups (first 8 chars of random part)
	keyID = randomPart[:8]

	// Construct the full key with prefix and version
	fullKey = fmt.Sprintf("%s_%s_%s", g.prefix, g.version, randomPart)

	return fullKey, keyID, nil
}

func (g *APIKeyGenerator) ParseKey(fullKey string) (prefix, version, randomPart string, err error) {
	// Only split on first 2 underscores, random part may contain underscores
	parts := strings.SplitN(fullKey, "_", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid key format: expected 3 parts, got %d", len(parts))
	}

	return parts[0], parts[1], parts[2], nil
}

// ValidateFormat checks if the key has the correct format
func (g *APIKeyGenerator) ValidateFormat(fullKey string) bool {
	prefix, version, randomPart, err := g.ParseKey(fullKey)
	if err != nil {
		return false
	}

	// Verify prefix matches
	if prefix != g.prefix {
		return false
	}

	// Verify version is recognized
	if version != g.version {
		return false
	}

	// Verify random part has expected length
	expectedLen := base64.RawURLEncoding.EncodedLen(KeyLength)
	if len(randomPart) != expectedLen {
		return false
	}

	return true
}
