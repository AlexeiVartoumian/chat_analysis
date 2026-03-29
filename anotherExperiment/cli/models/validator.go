package models

import (
	"time"
)

// Scope represents a permission scope for the API key
type Scope string

const (
	ScopeRead       Scope = "read"
	ScopeWrite      Scope = "write"
	ScopeDelete     Scope = "delete"
	ScopeAdmin      Scope = "admin"
	ScopeMonitoring Scope = "monitoring"
	ScopeAlerts     Scope = "alerts"
)

// APIKey represents an API key with its metadata and permissions
type APIKey struct {
	ID          string `json:"id" db:"id"`
	KeyID       string `json:"key_id" db:"key_id"` // Short identifier for lookups
	HashedKey   string `json:"-" db:"hashed_key"`  // Never expose in JSON
	Name        string `json:"name" db:"name"`     // Human-readable name
	Description string `json:"description" db:"description"`

	// Ownership
	UserID    string `json:"user_id" db:"user_id"`
	ProjectID string `json:"project_id" db:"project_id"`

	// Permissions
	Scopes []Scope `json:"scopes" db:"scopes"`

	// Restrictions
	AllowedIPs      []string `json:"allowed_ips" db:"allowed_ips"`
	AllowedReferers []string `json:"allowed_referers" db:"allowed_referers"`

	// Rate limiting
	RateLimit int `json:"rate_limit" db:"rate_limit"` // Requests per minute

	// Lifecycle
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`

	// Rotation
	RotatedFromID *string    `json:"rotated_from_id,omitempty" db:"rotated_from_id"`
	RotationDue   *time.Time `json:"rotation_due,omitempty" db:"rotation_due"`

	// Status
	IsActive bool `json:"is_active" db:"is_active"`
}

// HasScope checks if the API key has a specific scope
func (k *APIKey) HasScope(scope Scope) bool {
	for _, s := range k.Scopes {
		if s == scope || s == ScopeAdmin {
			return true
		}
	}
	return false
}

// HasAnyScope checks if the API key has any of the specified scopes
func (k *APIKey) HasAnyScope(scopes ...Scope) bool {
	for _, scope := range scopes {
		if k.HasScope(scope) {
			return true
		}
	}
	return false
}

// IsValid checks if the key is currently valid
func (k *APIKey) IsValid() bool {
	if !k.IsActive {
		return false
	}

	if k.RevokedAt != nil {
		return false
	}

	if k.ExpiresAt != nil && time.Now().After(*k.ExpiresAt) {
		return false
	}

	return true
}

// NeedsRotation checks if the key is due for rotation
func (k *APIKey) NeedsRotation() bool {
	if k.RotationDue == nil {
		return false
	}
	return time.Now().After(*k.RotationDue)
}
