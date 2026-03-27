package auth

import (
	"api/models"
	"api/repository/sqlconnect"
	"context"
	"net"
	"net/http"
	"strings"
)

// Context keys for storing API key information
type contextKey string

const (
	APIKeyContextKey contextKey = "api_key"
)

// AuthMiddleware handles API key authentication
type AuthMiddleware struct {
	generator *APIKeyGenerator
	hasher    *KeyHasher
	store     *sqlconnect.PostgresStore
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(
	generator *APIKeyGenerator,
	hasher *KeyHasher,
	store *sqlconnect.PostgresStore,

) *AuthMiddleware {
	return &AuthMiddleware{
		generator: generator,
		hasher:    hasher,
		store:     store,
	}
}

// Authenticate returns an HTTP middleware that validates API keys
func (m *AuthMiddleware) Authenticate(requiredScopes ...models.Scope) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract API key from request
			apiKey := m.extractAPIKey(r)
			if apiKey == "" {

				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			// Validate key format
			if !m.generator.ValidateFormat(apiKey) {

				http.Error(w, "Invalid API key format", http.StatusUnauthorized)
				return
			}

			// Extract key ID for lookup
			_, _, randomPart, _ := m.generator.ParseKey(apiKey)
			keyID := randomPart[:8]

			// Fetch key from store
			key, err := m.store.GetByKeyID(ctx, keyID)
			if err != nil {

				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Verify the key hash
			valid, err := m.hasher.Verify(apiKey, key.HashedKey)
			if err != nil || !valid {

				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Check if key is active and not expired
			// if !key.IsValid() {
			// 	reason := "key_inactive"
			// 	if key.RevokedAt != nil {
			// 		reason = "key_revoked"
			// 	} else if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
			// 		reason = "key_expired"
			// 	}

			// 	http.Error(w, "API key is no longer valid", http.StatusUnauthorized)
			// 	return
			// }

			// Check IP restrictions
			if len(key.AllowedIPs) > 0 && !m.isIPAllowed(r, key.AllowedIPs) {

				http.Error(w, "IP address not allowed", http.StatusForbidden)
				return
			}

			// Check referer restrictions
			if len(key.AllowedReferers) > 0 && !m.isRefererAllowed(r, key.AllowedReferers) {

				http.Error(w, "Referer not allowed", http.StatusForbidden)
				return
			}

			// Check required scopes
			for _, scope := range requiredScopes {
				if !key.HasScope(scope) {

					http.Error(w, "Insufficient permissions", http.StatusForbidden)
					return
				}
			}

			// Check rate limit
			// result, err := m.rateLimiter.Check(ctx, keyID, int64(key.RateLimit))
			// if err != nil {
			// 	http.Error(w, "Rate limit check failed", http.StatusInternalServerError)
			// 	return
			// }

			// Set rate limit headers
			// w.Header().Set("X-RateLimit-Limit", string(rune(key.RateLimit)))
			// w.Header().Set("X-RateLimit-Remaining", string(rune(result.Remaining)))
			// w.Header().Set("X-RateLimit-Reset", result.ResetAt.Format(time.RFC3339))

			// if !result.Allowed {
			// 	w.Header().Set("Retry-After", string(rune(int(result.RetryAfter.Seconds()))))
			// 	m.auditLogger.LogRateLimited(ctx, keyID, r.RemoteAddr)
			// 	http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			// 	return
			// }

			// Update last used timestamp asynchronously
			go func() {
				_ = m.store.UpdateLastUsed(context.Background(), keyID)
			}()

			// Log successful authentication
			//m.auditLogger.LogAuthSuccess(ctx, keyID, r.RemoteAddr, r.URL.Path)

			// Add key to context for handlers
			ctx = context.WithValue(ctx, APIKeyContextKey, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractAPIKey extracts the API key from the request
// Supports: Authorization header, X-API-Key header, and query parameter
func (m *AuthMiddleware) extractAPIKey(r *http.Request) string {
	// Check Authorization header (Bearer token)
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// Check X-API-Key header
	if key := r.Header.Get("X-API-Key"); key != "" {
		return key
	}

	// Check query parameter (not recommended for production)
	if key := r.URL.Query().Get("api_key"); key != "" {
		return key
	}

	return ""
}

// isIPAllowed checks if the request IP is in the allowed list
func (m *AuthMiddleware) isIPAllowed(r *http.Request, allowedIPs []string) bool {
	clientIP := m.getClientIP(r)

	for _, allowed := range allowedIPs {
		// Check for CIDR notation
		if strings.Contains(allowed, "/") {
			_, network, err := net.ParseCIDR(allowed)
			if err == nil && network.Contains(net.ParseIP(clientIP)) {
				return true
			}
		} else if clientIP == allowed {
			return true
		}
	}

	return false
}

// getClientIP extracts the client IP considering proxies
func (m *AuthMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// isRefererAllowed checks if the request referer matches allowed patterns
func (m *AuthMiddleware) isRefererAllowed(r *http.Request, allowedReferers []string) bool {
	referer := r.Header.Get("Referer")
	if referer == "" {
		return false
	}

	for _, allowed := range allowedReferers {
		// Support wildcard matching
		if strings.HasSuffix(allowed, "*") {
			prefix := strings.TrimSuffix(allowed, "*")
			if strings.HasPrefix(referer, prefix) {
				return true
			}
		} else if referer == allowed {
			return true
		}
	}

	return false
}

// GetAPIKeyFromContext retrieves the API key from the request context
func GetAPIKeyFromContext(ctx context.Context) *models.APIKey {
	key, ok := ctx.Value(APIKeyContextKey).(*models.APIKey)
	if !ok {
		return nil
	}
	return key
}
