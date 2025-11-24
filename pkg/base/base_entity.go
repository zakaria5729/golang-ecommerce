package base

import (
	"time"

	c "github.com/easy-comerce/backend/pkg/constants"
	tu "github.com/easy-comerce/backend/pkg/tokenutil"
	u "github.com/easy-comerce/backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint       `json:"id" gorm:"primarykey; column:id;"`
	UUID      uuid.UUID  `json:"uuid" gorm:"unique; not null; column:uuid;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at; autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at; autoUpdateTime"`
}

func (be *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid, err := u.RetryWithDelay(tu.GenerateNewUUID, c.RetryLimit, 5)
	if err != nil {
		return err
	}

	be.UUID = uuid
	return nil
}
