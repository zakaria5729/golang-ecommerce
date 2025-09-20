package user

import (
	"time"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/models"
)

type UserResponse struct {
	models.BaseModel
	LastLoginAt    *time.Time  `json:"last_login_at"`
	ImageURL       *string     `json:"image_url"`
	Roles          []role.Role `json:"roles"`
	Email          string      `json:"email"`
	Name           string      `json:"name"`
	Verified       bool        `json:"verified"`
	Banned         bool        `json:"banned"`
	TotalPurchases uint        `json:"total_purchases"`
	TotalSpent     uint        `json:"total_spent"`
}
