package user

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
	RoleID   uint   `json:"role_id" validate:"required"`
	Verified bool   `json:"verified" validate:"required"`
	Banned   bool   `json:"banned" validate:"required"`
}

type UpdateUserRequest struct {
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
	Verified bool   `json:"verified" validate:"required"`
	Banned   bool   `json:"banned" validate:"required"`
}

type UpdateProfileRequest struct {
	Name    string  `json:"name" validate:"required,min=2"`
	PathKey *string `json:"path_key" validate:"required"`
}
