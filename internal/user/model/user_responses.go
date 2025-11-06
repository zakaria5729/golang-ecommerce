package model

import (
	"time"

	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/pkg/base"
)

type UserResponse struct {
	AuditEntity      base.AuditEntity
	LastLoginAt      *time.Time        `json:"last_login_at"`
	ImageURL         *string           `json:"image_url"`
	Roles            []role.RoleEntity `json:"roles"`
	Email            string            `json:"email"`
	Name             string            `json:"name"`
	Verified         bool              `json:"verified"`
	Banned           bool              `json:"banned"`
	PurchaseCount    uint              `json:"purchase_count"`
	TotalSpent       float64           `json:"total_spent"`
	VerificationLink *string           `json:"verification_link,omitempty"`
}
