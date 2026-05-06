package memory

import (
	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"
)

// UserRepository is an in-memory implementation of repository.UserRepository
// backed by ligo-memory.Store.
type UserRepository struct {
	store *ligomemory.Store[int, *entity.User]
	log   ligo.Logger
}

func NewUserRepository(store *ligomemory.Store[int, *entity.User], log ligo.Logger) repository.UserRepository {
	for _, u := range []*entity.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "user"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "user"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Role: "admin"},
	} {
		store.Set(u.ID, u)
	}
	return &UserRepository{store: store, log: log}
}

// SeedDatabase initializes the repository with seed data.
// This method is registered as a lifecycle hook via Register().
func (r *UserRepository) SeedDatabase() error {
	r.log.Info("Seeding user database with initial data")
	// Seed data is already added in NewUserRepository,
	// but this demonstrates the hook pattern.
	return nil
}

// CleanupDatabase performs cleanup when the module is destroyed.
// This method is registered as a lifecycle hook via Register().
func (r *UserRepository) CleanupDatabase() error {
	r.log.Info("Cleaning up user database")
	// Clear all users on shutdown (optional)
	// r.store.Clear()
	return nil
}

// Register implements ligo.Registerable interface for compile-time safe hook registration.
// This method is called automatically when using ligo.HookedFactory.
//
// Benefits of this pattern:
// 1. Compile-time safety: typos in method names are caught at compile time
// 2. Explicit registration: clear what hooks are registered
// 3. Type safety: method signatures are checked at compile time
func (r *UserRepository) Register(registry *ligo.HookRegistry) {
	registry.OnInit(r.SeedDatabase)
	registry.OnDestroy(r.CleanupDatabase)
}

func (r *UserRepository) FindByID(id int) (*entity.User, bool) {
	return r.store.Get(id)
}

func (r *UserRepository) FindAll() []*entity.User {
	return r.store.All()
}

func (r *UserRepository) Create(name, email string) *entity.User {
	id := nextID()
	user := &entity.User{ID: id, Name: name, Email: email, Role: "user"}
	r.store.Set(id, user)
	return user
}

func (r *UserRepository) Update(id int, name, email string) (*entity.User, bool) {
	existing, found := r.store.Get(id)
	if !found {
		return nil, false
	}
	updated := &entity.User{ID: id, Name: name, Email: email, Role: existing.Role}
	r.store.Set(id, updated)
	return updated, true
}

func (r *UserRepository) Delete(id int) bool {
	return r.store.Delete(id)
}
