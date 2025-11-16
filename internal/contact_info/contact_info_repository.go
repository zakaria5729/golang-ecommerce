package contact_info

import (
	"context"

	"github.com/easy-comerce/backend/pkg/base"
	c "github.com/easy-comerce/backend/pkg/constants"
	op "github.com/easy-comerce/backend/pkg/option"
	"gorm.io/gorm"
)

type ContactInfoRepository interface {
	base.BaseRepository[ContactInfoEntity]
	GetContactInfosPaginated(ctx context.Context, page int, pageSize int, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*[]ContactInfoEntity, int64, error)
}

type contactInfoRepository struct {
	base.BaseRepository[ContactInfoEntity]
	db *gorm.DB
}

func NewContactInfoRepository(db *gorm.DB) ContactInfoRepository {
	return &contactInfoRepository{
		base.NewBaseRepository[ContactInfoEntity](db),
		db,
	}
}

func (r *contactInfoRepository) GetContactInfosPaginated(ctx context.Context, page int, pageSize int, sortBy string, sortOrder string, contactType *string, showMessage *bool, showDeleted *bool) (*[]ContactInfoEntity, int64, error) {
	options := op.QueryOptions{
		ShowDeleted: showDeleted,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
		SortableFields: []string{
			c.ContactInfoType,
		},
	}

	if contactType != nil && *contactType != "" {
		options.AddFilter(c.ContactInfoType+" = ?", *contactType)
	}

	selectFields := []string{
		c.FieldID, c.FieldCreatedAt, c.FieldUpdatedAt, c.ContactInfoName,
		c.ContactInfoEmail, c.ContactInfoType, c.ContactInfoPhone,
	}

	if showMessage != nil && *showMessage {
		selectFields = append(selectFields, c.ContactInfoMessage)
	}

	contacts, total, err := r.GetAllPaginated(ctx, page, pageSize, &options, selectFields...)
	if err != nil {
		return nil, 0, err
	}

	return &contacts, total, nil
}
