package user

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
}

type UpdateUserRequest struct {
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
}

type UpdateProfileRequest struct {
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
}
