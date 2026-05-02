package service

// ContextKey constants for request-scoped data.
const (
	ContextKeyUser = "user"
)

// AuthService defines the authentication service interface.
type AuthService interface {
	// ValidateToken validates a bearer token and returns the authenticated user.
	ValidateToken(token string) (User, error)
}

// User represents an authenticated user.
// This interface is defined here to avoid circular dependencies.
type User interface {
	GetID() string
	GetName() string
	GetEmail() string
	GetRole() string
	HasRole(role string) bool
	IsAdmin() bool
}
