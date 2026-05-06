package repository

import (
	"{{.ModulePath}}/internal/domain/entity"
)

// UserRepository defines the interface for user data access.
// This is the repository contract - implementations can be memory, database, etc.
type UserRepository interface {
	// FindByID retrieves a user by ID.
	// Returns the user and true if found, nil and false otherwise.
	FindByID(id int) (*entity.User, bool)

	// FindAll returns all users.
	FindAll() []*entity.User

	// Create adds a new user and returns the created entity.
	Create(name, email string) *entity.User

	// Update updates an existing user.
	// Returns the updated user and true if successful, nil and false if not found.
	Update(id int, name, email string) (*entity.User, bool)

	// Delete removes a user by ID.
	// Returns true if deleted, false if not found.
	Delete(id int) bool
}
