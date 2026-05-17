package presenter

import (
	"{{.ModulePath}}/internal/domain/entity"
)

// UserPresenter handles response formatting for user endpoints.
type UserPresenter struct{}

// NewUserPresenter creates a new user presenter.
func NewUserPresenter() *UserPresenter {
	return &UserPresenter{}
}

// ToResponse converts a user entity to response format.
func (p *UserPresenter) ToResponse(user *entity.User) map[string]any {
	return map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
}

// ToListResponse converts a list of users to response format.
func (p *UserPresenter) ToListResponse(users []*entity.User) map[string]any {
	list := make([]map[string]any, len(users))
	for i, user := range users {
		list[i] = p.ToResponse(user)
	}
	return map[string]any{
		"users": list,
		"count": len(list),
	}
}
