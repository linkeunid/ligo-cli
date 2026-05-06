package dto

// UpdateUserInput represents the input for updating a user.
type UpdateUserInput struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}
