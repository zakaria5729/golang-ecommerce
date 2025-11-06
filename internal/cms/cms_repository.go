package cms

import (
	"errors"

	en "github.com/easy-comerce/backend/internal/cms/entity"
	m "github.com/easy-comerce/backend/internal/cms/model"
	e "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CmsRepository interface {
	CreatePage(page *en.CmsPageEntity) (uint, error)
	GetPageByTag(tag string) (*en.CmsPageEntity, error)
	ExistsPageByTag(tag string) (bool, error)
	ExistsPageByTagAndIdNot(tag string, id uint) (bool, error)
	GetSectionsByPageID(pageId uint, title *string, sortBy string, sortOrder string, showDeleted *bool) ([]m.CmsSectionResponse, error)
	GetSectionByIdAndPageID(id uint, pageId uint) (*en.CmsSectionEntity, error)
	DeleteSectionsByPageID(pageId uint) error
	DeletePageByID(pageId uint, hardDelete bool) error
	UpsertSections(sections *[]en.CmsSectionEntity) error
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

func (r *cmsRepository) GetSectionsByPageID(pageId uint, title *string, sortBy string, sortOrder string, showDeleted *bool) ([]m.CmsSectionResponse, error) {
	var sections []en.CmsSectionEntity
	query := r.db.Where(c.PagePageID+" = ?", pageId)

	if title != nil && *title != "" {
		query = query.Where(c.SectionTitle+" LIKE ?", "%"+*title+"%")
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

	return sectionsResponse, nil
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

func (r *cmsRepository) DeletePageByID(pageId uint, hardDelete bool) error {
	query := r.db.Where(c.FieldID+" = ?", pageId)
	if hardDelete {
		query = query.Unscoped()
	}

	err := query.Delete(&en.CmsPageEntity{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.Error("Failed to delete page by id", "error", err)
			return e.WrapServerError("Failed to delete page by id", err)
		}

		return errors.New("Page not found with this id")
	}

	return nil
}

func (r *cmsRepository) UpsertSections(sections *[]en.CmsSectionEntity) error {
	if sections == nil || len(*sections) == 0 {
		return nil
	}
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: c.FieldID}},
		DoUpdates: clause.AssignmentColumns([]string{c.SectionTitle, c.SectionDescription, c.SectionContent, c.SectionLink, c.SectionImageUrls, c.PagePageID, c.FieldCreatedAt, c.FieldUpdatedBy}),
	}).CreateInBatches(*sections, 100).Error

	if err != nil {
		l.Logger.Error("Failed to create sections", "error", err)
		return e.WrapServerError("Failed to create sections", err)
	}

	return nil
}
