package user

import (
	"time"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/models"
)

type UserResponse struct {
	models.BaseModel
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Verified    bool        `json:"verified"`
	Banned      bool        `json:"banned"`
	LastLoginAt *time.Time  `json:"last_login_at"`
	ImageURL    *string     `json:"image_url"`
	Roles       []role.Role `json:"roles"`
}
