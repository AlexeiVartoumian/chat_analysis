package sqlconnect

import (
	"api/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

// PostgresStore implements API key storage with PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new store instance
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// CreateTable initializes the API keys table with proper indexes
func (s *PostgresStore) CreateTable(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS api_keys (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        key_id VARCHAR(16) UNIQUE NOT NULL,
        hashed_key TEXT NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        user_id UUID NOT NULL,
        project_id UUID NOT NULL,
        scopes TEXT[] NOT NULL DEFAULT '{}',
        allowed_ips TEXT[] DEFAULT '{}',
        allowed_referers TEXT[] DEFAULT '{}',
        rate_limit INTEGER DEFAULT 1000,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        expires_at TIMESTAMP WITH TIME ZONE,
        last_used_at TIMESTAMP WITH TIME ZONE,
        revoked_at TIMESTAMP WITH TIME ZONE,
        rotated_from_id UUID REFERENCES api_keys(id),
        rotation_due TIMESTAMP WITH TIME ZONE,
        is_active BOOLEAN DEFAULT true
    );

    -- Index for fast key lookups by key_id
    CREATE INDEX IF NOT EXISTS idx_api_keys_key_id ON api_keys(key_id) WHERE is_active = true;

    -- Index for finding keys by user
    CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);

    -- Index for finding keys by project
    CREATE INDEX IF NOT EXISTS idx_api_keys_project_id ON api_keys(project_id);

    -- Index for finding expired keys
    CREATE INDEX IF NOT EXISTS idx_api_keys_expires_at ON api_keys(expires_at) WHERE expires_at IS NOT NULL;

    -- Index for finding keys due for rotation
    CREATE INDEX IF NOT EXISTS idx_api_keys_rotation_due ON api_keys(rotation_due) WHERE rotation_due IS NOT NULL;
    `

	_, err := s.db.ExecContext(ctx, query)
	return err
}

// Create stores a new API key
func (s *PostgresStore) Create(ctx context.Context, key *models.APIKey) error {
	query := `
    INSERT INTO api_keys (
        key_id, hashed_key, name, description, user_id, project_id,
        scopes, allowed_ips, allowed_referers, rate_limit, expires_at, rotation_due
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    RETURNING id, created_at
    `

	scopes := make([]string, len(key.Scopes))
	for i, s := range key.Scopes {
		scopes[i] = string(s)
	}

	return s.db.QueryRowContext(ctx, query,
		key.KeyID,
		key.HashedKey,
		key.Name,
		key.Description,
		key.UserID,
		key.ProjectID,
		pq.Array(scopes),
		pq.Array(key.AllowedIPs),
		pq.Array(key.AllowedReferers),
		key.RateLimit,
		key.ExpiresAt,
		key.RotationDue,
	).Scan(&key.ID, &key.CreatedAt)
}

// GetByKeyID retrieves an API key by its short identifier
func (s *PostgresStore) GetByKeyID(ctx context.Context, keyID string) (*models.APIKey, error) {
	query := `
    SELECT id, key_id, hashed_key, name, description, user_id, project_id,
           scopes, allowed_ips, allowed_referers, rate_limit, created_at,
           expires_at, last_used_at, revoked_at, rotated_from_id, rotation_due, is_active
    FROM api_keys
    WHERE key_id = $1
    `

	var key models.APIKey
	var scopes, allowedIPs, allowedReferers []string

	err := s.db.QueryRowContext(ctx, query, keyID).Scan(
		&key.ID, &key.KeyID, &key.HashedKey, &key.Name, &key.Description,
		&key.UserID, &key.ProjectID, pq.Array(&scopes), pq.Array(&allowedIPs),
		pq.Array(&allowedReferers), &key.RateLimit, &key.CreatedAt,
		&key.ExpiresAt, &key.LastUsedAt, &key.RevokedAt, &key.RotatedFromID,
		&key.RotationDue, &key.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found")
	}
	if err != nil {
		return nil, err
	}

	// Convert string slices to typed slices
	key.Scopes = make([]models.Scope, len(scopes))
	for i, s := range scopes {
		key.Scopes[i] = models.Scope(s)
	}
	key.AllowedIPs = allowedIPs
	key.AllowedReferers = allowedReferers

	return &key, nil
}

// UpdateLastUsed updates the last used timestamp
func (s *PostgresStore) UpdateLastUsed(ctx context.Context, keyID string) error {
	query := `UPDATE api_keys SET last_used_at = NOW() WHERE key_id = $1`
	_, err := s.db.ExecContext(ctx, query, keyID)
	return err
}

// Revoke marks an API key as revoked
func (s *PostgresStore) Revoke(ctx context.Context, keyID string) error {
	query := `UPDATE api_keys SET revoked_at = NOW(), is_active = false WHERE key_id = $1`
	_, err := s.db.ExecContext(ctx, query, keyID)
	return err
}

// ListByUser retrieves all API keys for a user
func (s *PostgresStore) ListByUser(ctx context.Context, userID string) ([]*models.APIKey, error) {
	query := `
    SELECT id, key_id, name, description, scopes, rate_limit, created_at,
           expires_at, last_used_at, is_active
    FROM api_keys
    WHERE user_id = $1
    ORDER BY created_at DESC
    `

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*models.APIKey
	for rows.Next() {
		var key models.APIKey
		var scopes []string

		err := rows.Scan(
			&key.ID, &key.KeyID, &key.Name, &key.Description,
			pq.Array(&scopes), &key.RateLimit, &key.CreatedAt,
			&key.ExpiresAt, &key.LastUsedAt, &key.IsActive,
		)
		if err != nil {
			return nil, err
		}

		key.Scopes = make([]models.Scope, len(scopes))
		for i, s := range scopes {
			key.Scopes[i] = models.Scope(s)
		}

		keys = append(keys, &key)
	}

	return keys, rows.Err()
}
