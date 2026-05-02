package memory

import (
	ligomemory "github.com/linkeunid/ligo-memory"
	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"
)

// UserRepository is an in-memory implementation of repository.UserRepository
// backed by ligo-memory.Store.
type UserRepository struct {
	store *ligomemory.Store[string, *entity.User]
}

func NewUserRepository(store *ligomemory.Store[string, *entity.User]) repository.UserRepository {
	for _, u := range []*entity.User{
		{ID: "550e8400-e29b-41d4-a716-446655440001", Name: "Alice", Email: "alice@example.com", Role: "user"},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Name: "Bob", Email: "bob@example.com", Role: "user"},
		{ID: "550e8400-e29b-41d4-a716-446655440003", Name: "Charlie", Email: "charlie@example.com", Role: "admin"},
	} {
		store.Set(u.ID, u)
	}
	return &UserRepository{store: store}
}

func (r *UserRepository) FindByID(id string) (*entity.User, bool) {
	return r.store.Get(id)
}

func (r *UserRepository) FindAll() []*entity.User {
	return r.store.All()
}

func (r *UserRepository) Create(name, email string) *entity.User {
	id := newUUID()
	user := &entity.User{ID: id, Name: name, Email: email, Role: "user"}
	r.store.Set(id, user)
	return user
}

func (r *UserRepository) Update(id, name, email string) (*entity.User, bool) {
	existing, found := r.store.Get(id)
	if !found {
		return nil, false
	}
	updated := &entity.User{ID: id, Name: name, Email: email, Role: existing.Role}
	r.store.Set(id, updated)
	return updated, true
}

func (r *UserRepository) Delete(id string) bool {
	return r.store.Delete(id)
}
