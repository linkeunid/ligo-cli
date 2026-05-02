package entity

// Role constants for user roles.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents a user entity in the domain layer.
// Contains only business data and behavior, no external dependencies.
type User struct {
	ID    string
	Name  string
	Email string
	Role  string
}

// GetID returns the user ID (implements service.User).
func (u *User) GetID() string {
	return u.ID
}

// GetName returns the user name (implements service.User).
func (u *User) GetName() string {
	return u.Name
}

// GetEmail returns the user email (implements service.User).
func (u *User) GetEmail() string {
	return u.Email
}

// GetRole returns the user role (implements service.User).
func (u *User) GetRole() string {
	return u.Role
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// IsAdmin checks if user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsRegular checks if user has regular user role.
func (u *User) IsRegular() bool {
	return u.Role == RoleUser || u.Role == ""
}
