package cms

import (
	"context"
	"errors"
	"fmt"

	en "github.com/easy-comerce/backend/internal/cms/entity"
	m "github.com/easy-comerce/backend/internal/cms/model"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type CmsService interface {
	CreatePageIfNotExists(tag string, title string)
	CreatePageSection(ctx context.Context, request *m.CmsPageRequest) (*m.CmsPageResponse, error)
	UpdatePageSection(ctx context.Context, pageID uint, request *m.CmsPageRequest) (*m.CmsPageResponse, error)
	DeletePageSection(ctx context.Context, pageID uint) error
	UndoDeletePageSection(ctx context.Context, pageID uint) error
	GetPageByTag(tag string, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error)
	GetPageByID(id uint, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error)
	GetSectionsDataWithPageID(ctx context.Context, pageID uint, sectionRequests *[]m.CmsSectionRequest, isFromUpdate bool) (*[]en.CmsSectionEntity, error)
	GetPagesPaginated(pageStr string, pageSizeStr string, sortBy string, sortOrder string, showDeleted *bool, title *string, includeSections *bool) (*response.PaginatedResponse, error)
}

type cmsService struct {
	repo CmsRepository
}

func NewCmsService(repo CmsRepository) CmsService {
	return &cmsService{
		repo: repo,
	}
}

func (s *cmsService) CreatePageIfNotExists(tag string, name string) {
	page, _ := s.repo.GetPageByTag(tag)

	if page == nil {
		page = &en.CmsPageEntity{
			Tag:  tag,
			Name: name,
		}
		s.repo.CreatePage(page)
	}
}

func (s *cmsService) GetPageByTag(tag string, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error) {
	page, err := s.repo.GetPageByTag(tag)
	if err != nil || page == nil {
		return nil, errors.New("Page not found with this tag")
	}

	if includeSections != nil && *includeSections {
		sections, _ := s.repo.GetSectionsByPageID(page.ID, title, sortBy, sortOrder, showDeleted)
		if sections != nil {
			pageResponse := page.ToResponse()
			pageResponse.Sections = sections
			return pageResponse, nil
		}
	}

	return page.ToResponse(), nil
}

func (s *cmsService) GetPageByID(id uint, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error) {
	page, err := s.repo.GetPageByID(id, showDeleted)
	if err != nil || page == nil {
		return nil, errors.New("Page not found with this id")
	}

	if includeSections != nil && *includeSections {
		sections, _ := s.repo.GetSectionsByPageID(page.ID, title, sortBy, sortOrder, showDeleted)
		if sections != nil {
			pageResponse := page.ToResponse()
			pageResponse.Sections = sections
			return pageResponse, nil
		}
	}

	return page.ToResponse(), nil
}

func (s *cmsService) GetPagesPaginated(pageStr string, pageSizeStr string, sortBy string, sortOrder string, showDeleted *bool, title *string, includeSections *bool) (*response.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	pages, total, err := s.repo.GetPagesPaginated(page, pageSize, sortBy, sortOrder, showDeleted, title)
	if err != nil || pages == nil {
		return nil, fmt.Errorf("failed to fetch pages: %w", err)
	}

	var pageResponses []m.CmsPageResponse
	for _, page := range *pages {
		pageResponse := page.ToResponse()

		if includeSections != nil && *includeSections {
			sections, _ := s.repo.GetSectionsByPageID(page.ID, title, sortBy, sortOrder, showDeleted)
			if sections != nil {
				pageResponse.Sections = sections
			}
		}
		pageResponses = append(pageResponses, *pageResponse)
	}

	return utils.BuildPaginatedResponse(pageResponses, total, page, pageSize), nil
}

func (s *cmsService) CreatePageSection(ctx context.Context, request *m.CmsPageRequest) (*m.CmsPageResponse, error) {
	exists, _ := s.repo.ExistsPageByTag(request.Tag)
	if exists {
		return nil, errors.New("This tag (" + request.Tag + ") already exists with other page. Please try with a different one")
	}

	page := &en.CmsPageEntity{
		Tag:         request.Tag,
		Name:        request.Name,
		Description: request.Description,
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	page.CreatedBy = userID

	pageID, err := s.repo.CreatePage(page)
	if err != nil {
		return nil, err
	}

	sections, err := s.GetSectionsDataWithPageID(ctx, pageID, request.Sections, false)
	if err != nil {
		return nil, err
	}

	return upsertSectionsByPage(ctx, s.repo, page, sections, false)
}

func (s *cmsService) UpdatePageSection(ctx context.Context, pageID uint, request *m.CmsPageRequest) (*m.CmsPageResponse, error) {
	_, err := s.repo.ExistsPageByID(pageID, false)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsPageByTagAndIdNot(request.Tag, pageID)
	if exists {
		return nil, errors.New("This tag (" + request.Tag + ") already exists with other page. Please try with a different one")
	}

	page, err := s.repo.GetPageByID(pageID, nil)
	if err != nil || page == nil {
		return nil, err
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	page.UpdatedBy = userID
	page.Tag = request.Tag
	page.Name = request.Name
	page.Description = request.Description

	_, err = s.repo.UpdatePage(pageID, page)
	if err != nil {
		return nil, err
	}

	sections, err := s.GetSectionsDataWithPageID(ctx, pageID, request.Sections, true)
	if err != nil {
		return nil, err
	}

	return upsertSectionsByPage(ctx, s.repo, page, sections, true)
}

func (s *cmsService) GetSectionsDataWithPageID(ctx context.Context, pageID uint, sectionRequests *[]m.CmsSectionRequest, isFromUpdate bool) (*[]en.CmsSectionEntity, error) {
	if sectionRequests == nil || len(*sectionRequests) == 0 {
		return nil, nil
	}
	sections := []en.CmsSectionEntity{}

	for _, secRequest := range *sectionRequests {
		var err error
		var section *en.CmsSectionEntity

		if isFromUpdate {
			if secRequest.ID == nil || *secRequest.ID == 0 {
				section, err = getSectionData(ctx, pageID, &en.CmsSectionEntity{}, &secRequest, false)
				if err != nil || section == nil {
					return nil, err
				}
				sections = append(sections, *section)
				continue
			}

			section, err = s.repo.GetSectionByIdAndPageID(*secRequest.ID, pageID)
			if err != nil || section == nil {
				return nil, errors.New("No section found with this id: " + fmt.Sprint(*secRequest.ID) + " and page id: " + fmt.Sprint(pageID))
			}

			section, err = getSectionData(ctx, pageID, section, &secRequest, true)
			if err != nil || section == nil {
				return nil, err
			}
			sections = append(sections, *section)
		} else {
			section, err = getSectionData(ctx, pageID, &en.CmsSectionEntity{}, &secRequest, false)
			if err != nil || section == nil {
				return nil, err
			}
			sections = append(sections, *section)
		}
	}

	return &sections, nil
}

func (s *cmsService) DeletePageSection(ctx context.Context, pageID uint) error {
	exists, err := s.repo.ExistsPageByID(pageID, false)
	if err != nil || !exists {
		return err
	}

	return s.repo.DeletePageByID(ctx, pageID)
}

func (s *cmsService) UndoDeletePageSection(ctx context.Context, pageID uint) error {
	exists, err := s.repo.ExistsPageByID(pageID, true)
	if err != nil || !exists {
		return err
	}

	return s.repo.UndoDeletePageByID(ctx, pageID)
}

func upsertSectionsByPage(ctx context.Context, repo CmsRepository, page *en.CmsPageEntity, sections *[]en.CmsSectionEntity, isFromUpdate bool) (*m.CmsPageResponse, error) {
	if page == nil || sections == nil || len(*sections) == 0 {
		return nil, nil
	}

	if isFromUpdate {
		existingIDs, _ := repo.GetSectionIdsByPageID(page.ID)

		if existingIDs != nil && len(*existingIDs) > 0 {
			updatableSectionIds := make(map[uint]bool)
			for _, section := range *sections {
				if section.ID > 0 {
					updatableSectionIds[section.ID] = true
				}
			}

			var deletableSectionIds []uint
			for _, id := range *existingIDs {
				if !updatableSectionIds[id] {
					deletableSectionIds = append(deletableSectionIds, id)
				}
			}
			if len(deletableSectionIds) > 0 {
				repo.DeleteSectionsByIds(ctx, deletableSectionIds)
			}
		}
	}

	if err := repo.UpsertSections(sections); err != nil {
		return nil, err
	}

	resp := make([]m.CmsSectionResponse, len(*sections))
	for i, s := range *sections {
		resp[i] = *s.ToResponse()
	}

	pageResp := page.ToResponse()
	pageResp.Sections = &resp
	return pageResp, nil
}

func getSectionData(ctx context.Context, pageId uint, section *en.CmsSectionEntity, request *m.CmsSectionRequest, isFromUpdate bool) (*en.CmsSectionEntity, error) {
	if section == nil || request == nil {
		l.Logger.Error("Invalid section request", "section", section, "request", request)
		return nil, errors.New("Invalid section request")
	}

	if request.Title == "" || request.Content == "" {
		l.Logger.Error("Invalid section request", "section", section, "request", request)
		return nil, errors.New("Invalid section request. Title and Content is required")
	}

	section.PageId = pageId
	section.Title = request.Title
	section.Content = request.Content
	section.Link = request.Link
	section.Description = request.Description
	section.SetImageUrls(request.ImageUrls)

	userID, _ := cu.GetUserIDFromContext(ctx)
	if isFromUpdate {
		section.UpdatedBy = userID
	} else {
		section.CreatedBy = userID
	}

	return section, nil
}
