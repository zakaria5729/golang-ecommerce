package cms

import (
	"context"
	"errors"
	"fmt"

	en "github.com/easy-comerce/backend/internal/cms/entity"
	m "github.com/easy-comerce/backend/internal/cms/model"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	l "github.com/easy-comerce/backend/pkg/logger"
)

type CmsService interface {
	CreatePageIfNotExists(tag string, title string)
	CreatePageSection(ctx context.Context, request *m.CmsPageRequest) (*m.CmsPageResponse, error)
	GetByTag(tag string, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error)
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

func (s *cmsService) GetByTag(tag string, title *string, includeSections *bool, sortBy string, sortOrder string, showDeleted *bool) (*m.CmsPageResponse, error) {
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

func (s *cmsService) CreatePageSection(ctx context.Context, request *m.CmsPageRequest) (*m.CmsPageResponse, error) {
	if request == nil {
		return nil, errors.New("Request body is required")
	}

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

	sections, err := s.upsertSectionsWithPageID(ctx, pageID, request.Sections, false)
	if err != nil {
		s.repo.DeletePageByID(pageID, true)
		return nil, err
	}

	if sections != nil && len(*sections) > 0 {
		sectionsResponses := []m.CmsSectionResponse{}
		for _, section := range *sections {
			sectionsResponses = append(sectionsResponses, *section.ToResponse())
		}

		pageResponse := page.ToResponse()
		pageResponse.Sections = sectionsResponses
		return pageResponse, nil
	}

	return page.ToResponse(), nil
}

func (s *cmsService) upsertSectionsWithPageID(ctx context.Context, pageID uint, sectionRequests *[]m.CmsSectionRequest, isFromUpdate bool) (*[]en.CmsSectionEntity, error) {
	if sectionRequests == nil || len(*sectionRequests) == 0 {
		return nil, nil
	}
	sections := []en.CmsSectionEntity{}

	for _, secRequest := range *sectionRequests {
		if isFromUpdate {
			section, err := s.repo.GetSectionByIdAndPageID(secRequest.ID, pageID)
			if err != nil || section == nil {
				return nil, errors.New("No section found with this id: " + fmt.Sprint(secRequest.ID) + " and page id: " + fmt.Sprint(pageID))
			}

			section, err = getSectionData(ctx, pageID, section, &secRequest, true)
			if err != nil {
				return nil, err
			}
			sections = append(sections, *section)
		} else {
			section, err := getSectionData(ctx, pageID, &en.CmsSectionEntity{}, &secRequest, false)
			if err != nil {
				return nil, err
			}
			sections = append(sections, *section)
		}
	}

	if len(sections) > 0 {
		err := s.repo.UpsertSections(&sections)
		if err != nil {
			return nil, err
		}
	}

	return &sections, nil
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
