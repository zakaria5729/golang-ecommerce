package cms

import (
	"context"
	"errors"

	en "github.com/easy-comerce/backend/internal/cms/entity"
	m "github.com/easy-comerce/backend/internal/cms/model"
	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CmsRepository interface {
	CreatePage(page *en.CmsPageEntity) (uint, error)
	UpdatePage(pageID uint, page *en.CmsPageEntity) (*uint, error)
	GetPageByTag(tag string) (*en.CmsPageEntity, error)
	GetPageByID(id uint, showDeleted *bool) (*en.CmsPageEntity, error)
	ExistsPageByTag(tag string) (bool, error)
	ExistsPageByID(id uint, showDeleted bool) (bool, error)
	ExistsPageByTagAndIdNot(tag string, id uint) (bool, error)
	GetSectionsByPageID(pageId uint, title *string, sortBy string, sortOrder string, showDeleted *bool) (*[]m.CmsSectionResponse, error)
	GetSectionIdsByPageID(pageId uint) (*[]uint, error)
	GetSectionByIdAndPageID(id uint, pageId uint) (*en.CmsSectionEntity, error)
	DeleteSectionsByPageID(pageId uint) error
	DeletePageByID(ctx context.Context, pageId uint) error
	DeleteSectionsByIds(ctx context.Context, sectionIds []uint) error
	UndoDeletePageByID(ctx context.Context, pageId uint) error
	UpsertSections(sections *[]en.CmsSectionEntity) error
	GetPagesPaginated(page int, pageSize int, sortBy string, sortOrder string, showDeleted *bool, title *string) (*[]en.CmsPageEntity, int64, error)
}

type cmsRepository struct {
	db *gorm.DB
}

func NewCmsRepository(db *gorm.DB) CmsRepository {
	return &cmsRepository{db: db}
}

func (r *cmsRepository) CreatePage(page *en.CmsPageEntity) (uint, error) {
	err := r.db.Create(page).Error

	if err != nil {
		l.Logger.Error("Failed to create page", "error", err)
		return 0, e.WrapServerError("Failed to create page", err)
	}

	return page.ID, nil
}

func (r *cmsRepository) UpdatePage(pageID uint, page *en.CmsPageEntity) (*uint, error) {
	err := r.db.Model(&en.CmsPageEntity{}).Where(c.FieldID+" = ?", pageID).Updates(page).Error

	if err != nil {
		l.Logger.Error("Failed to update page", "error", err)
		return nil, e.WrapServerError("Failed to update page", err)
	}

	return &pageID, nil
}

func (r *cmsRepository) GetPageByTag(tag string) (*en.CmsPageEntity, error) {
	var page en.CmsPageEntity
	err := r.db.Where(c.PageTag+" = ?", tag).First(&page).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to get page by tag", "error", err)
			return nil, e.WrapServerError("Failed to get page by tag", err)
		}

		return nil, errors.New("Page not found with this tag")
	}

	return &page, nil
}

func (r *cmsRepository) GetPageByID(id uint, showDeleted *bool) (*en.CmsPageEntity, error) {
	var page en.CmsPageEntity
	query := r.db.Where(c.FieldID+" = ?", id)

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.First(&page).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to get page by id", "error", err)
			return nil, e.WrapServerError("Failed to get page by id", err)
		}

		return nil, errors.New("Page not found with this id")
	}

	return &page, nil
}

func (r *cmsRepository) GetPagesPaginated(page int, pageSize int, sortBy string, sortOrder string, showDeleted *bool, title *string) (*[]en.CmsPageEntity, int64, error) {
	var pages []en.CmsPageEntity
	var total int64

	query := r.db.Model(&en.CmsPageEntity{})
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if title != nil && *title != "" {
		query = query.Where(c.SectionTitle+" ILIKE ?", *title+"%")
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Count(&total).Error; err != nil {
		l.Logger.Error("❌ Failed to count pages", "method", "GetPagesPaginated", "error", err, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder, "title", title, "showDeleted", showDeleted)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&pages).Error
	if err != nil {
		l.Logger.Error("❌ Failed to fetch pages paginated", "method", "GetPagesPaginated", "error", err, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder, "title", title, "showDeleted", showDeleted)
		return nil, 0, err
	}

	return &pages, total, nil
}

func (r *cmsRepository) ExistsPageByTag(tag string) (bool, error) {
	var page en.CmsPageEntity
	err := r.db.Select(c.FieldID).Where(c.PageTag+" = ?", tag).Take(&page).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to check if page exists by tag", "error", err)
			return false, e.WrapServerError("Failed to check if page exists by tag", err)
		}

		return false, errors.New("Page not found with this tag")
	}

	return page.ID > 0, nil
}

func (r *cmsRepository) ExistsPageByID(id uint, showDeleted bool) (bool, error) {
	var page en.CmsPageEntity
	query := r.db.Select(c.FieldID)

	if showDeleted {
		query = query.Unscoped()
	}

	err := query.Where(c.FieldID+" = ?", id).Take(&page).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to check if page exists by id", "error", err)
			return false, e.WrapServerError("Failed to check if page exists by id", err)
		}

		return false, errors.New("Page not found with this id")
	}

	return page.ID > 0, nil
}

func (r *cmsRepository) ExistsPageByTagAndIdNot(tag string, id uint) (bool, error) {
	var page en.CmsPageEntity
	err := r.db.Select(c.FieldID).Where(c.PageTag+" = ? AND "+c.FieldID+" != ?", tag, id).Take(&page).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to check if page exists by tag and id not", "error", err)
			return false, e.WrapServerError("Failed to check if page exists by tag", err)
		}

		return false, errors.New("Page not found with this tag")
	}

	return page.ID > 0, nil
}

func (r *cmsRepository) GetSectionsByPageID(pageId uint, title *string, sortBy string, sortOrder string, showDeleted *bool) (*[]m.CmsSectionResponse, error) {
	var sections []en.CmsSectionEntity
	query := r.db.Where(c.PagePageID+" = ?", pageId)

	if title != nil && *title != "" {
		query = query.Where(c.SectionTitle+" ILIKE ?", *title+"%")
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.Model(&en.CmsSectionEntity{}).Find(&sections).Error
	if err != nil || sections == nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to get sections by page id", "error", err)
			return nil, e.WrapServerError("Failed to get sections by page id", err)
		}

		return nil, errors.New("Sections not found with this page id")
	}

	var sectionsResponse []m.CmsSectionResponse
	for _, section := range sections {
		sectionsResponse = append(sectionsResponse, *section.ToResponse())
	}

	return &sectionsResponse, nil
}

func (r *cmsRepository) GetSectionIdsByPageID(pageId uint) (*[]uint, error) {
	var sectionIds []uint
	err := r.db.Model(&en.CmsSectionEntity{}).Select(c.FieldID).Where(c.PagePageID+" = ?", pageId).Find(&sectionIds).Error

	if err != nil || sectionIds == nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to get section ids by page id", "error", err)
			return nil, e.WrapServerError("Failed to get section ids by page id", err)
		}

		return nil, errors.New("Section ids not found with this page id")
	}

	return &sectionIds, nil
}

func (r *cmsRepository) GetSectionByIdAndPageID(id uint, pageId uint) (*en.CmsSectionEntity, error) {
	var section en.CmsSectionEntity
	err := r.db.Where(c.PagePageID+" = ?", pageId).Where(c.FieldID+" = ?", id).First(&section).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to get section by id and page id", "error", err)
			return nil, e.WrapServerError("Failed to get section by id and page id", err)
		}

		return nil, errors.New("Section not found with this id and page id")
	}

	return &section, nil
}

func (r *cmsRepository) DeleteSectionsByPageID(pageId uint) error {
	err := r.db.Where(c.PagePageID+" = ?", pageId).Delete(&en.CmsSectionEntity{}).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to delete sections by page id", "error", err)
			return e.WrapServerError("Failed to delete sections by page id", err)
		}

		return errors.New("Page not found with this id")
	}

	return nil
}

func (r *cmsRepository) DeletePageByID(ctx context.Context, pageId uint) error {
	return handlePageDeletion(r, ctx, pageId, false)
}

func (r *cmsRepository) UndoDeletePageByID(ctx context.Context, pageId uint) error {
	return handlePageDeletion(r, ctx, pageId, true)
}

func (r *cmsRepository) DeleteSectionsByIds(ctx context.Context, sectionIds []uint) error {
	if len(sectionIds) == 0 {
		return nil
	}

	entity := en.CmsSectionEntity{}
	userID, _ := cu.GetUserIDFromContext(ctx)
	entity.DeletedAt = timeutil.GormNowUTC()
	entity.DeletedBy = userID

	err := r.db.Model(&en.CmsSectionEntity{}).Where(c.FieldID+" IN ?", sectionIds).Updates(entity).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to delete sections by ids", "error", err)
			return e.WrapServerError("Failed to delete sections by ids", err)
		}

		return errors.New("Sections not found with this ids")
	}

	return nil
}

func (r *cmsRepository) UpsertSections(sections *[]en.CmsSectionEntity) error {
	if sections == nil || len(*sections) == 0 {
		return nil
	}
	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: c.FieldID}},
		DoUpdates: clause.AssignmentColumns([]string{
			c.SectionTitle, c.SectionDescription,
			c.SectionContent, c.SectionLink,
			c.SectionImageUrls, c.PagePageID,
			c.FieldCreatedAt, c.FieldUpdatedBy,
		}),
	}).CreateInBatches(*sections, 100).Error

	if err != nil {
		l.Logger.Error("Failed to create sections", "error", err)
		return e.WrapServerError("Failed to create sections", err)
	}

	return nil
}

func handlePageDeletion(r *cmsRepository, ctx context.Context, pageId uint, isUndo bool) error {
	userID, _ := cu.GetUserIDFromContext(ctx)

	return r.db.Transaction(func(tx *gorm.DB) error {
		sectionEntity := en.CmsSectionEntity{}
		if isUndo {
			sectionEntity.DeletedAt = nil
			sectionEntity.DeletedBy = nil
		} else {
			sectionEntity.DeletedAt = timeutil.GormNowUTC()
			sectionEntity.DeletedBy = userID
		}

		sectionQuery := tx.Select(c.FieldDeletedAt, c.FieldDeletedBy).Where(c.PagePageID+" = ?", pageId)
		if isUndo {
			sectionQuery = sectionQuery.Unscoped()
		}

		err := sectionQuery.Updates(sectionEntity).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				l.Logger.Error("Failed to delete sections by page id", "error", err)
				return e.WrapServerError("Failed to delete sections by page id", err)
			}

			return errors.New("Page not found with this id")
		}

		pageEntity := en.CmsPageEntity{}
		if isUndo {
			pageEntity.DeletedAt = nil
			pageEntity.DeletedBy = nil
		} else {
			pageEntity.DeletedAt = timeutil.GormNowUTC()
			pageEntity.DeletedBy = userID
		}

		pageQuery := tx.Select(c.FieldDeletedAt, c.FieldDeletedBy).Where(c.FieldID+" = ?", pageId)
		if isUndo {
			pageQuery = pageQuery.Unscoped()
		}

		err = pageQuery.Updates(pageEntity).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				l.Logger.Error("Failed to delete page by id", "error", err)
				return e.WrapServerError("Failed to delete page by id", err)
			}

			return errors.New("Page not found with this id")
		}
		return nil
	})
}
