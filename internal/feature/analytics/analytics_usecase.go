package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
)

type AnalyticsUseCase struct {
	analyticsRepo *AnalyticsRepository
}

func NewAnalyticsUseCase() *AnalyticsUseCase {
	return &AnalyticsUseCase{
		analyticsRepo: NewAnalyticsRepository(),
	}
}

// Track a page view
func (uc *AnalyticsUseCase) TrackPageView(ctx context.Context, pageView *AnalyticsPageView) error {
	// Check if visitor exists
	visitor, err := uc.analyticsRepo.GetVisitorByID(pageView.VisitorID)
	if err != nil {
		// Create new visitor if doesn't exist
		visitor = &AnalyticsVisitor{
			VisitorID:      pageView.VisitorID,
			FirstVisit:     pageView.VisitedAt,
			LastVisit:      pageView.VisitedAt,
			TotalVisits:    1,
			TotalPageViews: 1,
			UniquePages:    1,
			Country:        pageView.Country,
			City:           pageView.City,
			Region:         pageView.Region,
			Timezone:       pageView.Timezone,
			IsBot:          pageView.IsBot,
			BotName:        pageView.BotName,
		}
		err = uc.analyticsRepo.CreateVisitor(visitor)
		if err != nil {
			logger.Logger.Error("Failed to create visitor", "method", "TrackPageView", "error", err, "visitorID", pageView.VisitorID)
			return err
		}
	} else {
		// Update existing visitor
		visitor.LastVisit = pageView.VisitedAt
		visitor.TotalPageViews++

		// Check if this is a new page for this visitor
		existingPageViews, err := uc.analyticsRepo.GetPageViewsByVisitor(pageView.VisitorID, 1000, 0)
		if err == nil {
			uniquePages := make(map[string]bool)
			for _, pv := range existingPageViews {
				uniquePages[pv.Path] = true
			}
			if !uniquePages[pageView.Path] {
				visitor.UniquePages++
			}
		}

		err = uc.analyticsRepo.UpdateVisitor(visitor)
		if err != nil {
			logger.Logger.Error("Failed to update visitor", "method", "TrackPageView", "error", err, "visitorID", pageView.VisitorID)
			return err
		}
	}

	// Create page view
	err = uc.analyticsRepo.CreatePageView(pageView)
	if err != nil {
		logger.Logger.Error("Failed to create page view", "method", "TrackPageView", "error", err, "visitorID", pageView.VisitorID, "path", pageView.Path)
		return err
	}

	// Update session if exists
	session, err := uc.analyticsRepo.GetSessionByID(pageView.SessionID)
	if err == nil {
		session.PageViews++
		session.IsBounce = false // If we're tracking a page view, it's not a bounce
		if session.ExitPage == nil {
			session.ExitPage = &pageView.Path
		} else {
			session.ExitPage = &pageView.Path
		}
		err = uc.analyticsRepo.UpdateSession(session)
		if err != nil {
			logger.Logger.Error("Failed to update session", "method", "TrackPageView", "error", err, "sessionID", pageView.SessionID)
		}
	}

	return nil
}

// Get analytics overview
func (uc *AnalyticsUseCase) GetAnalyticsOverview(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	overview, err := uc.analyticsRepo.GetAnalyticsOverview(startDate, endDate)
	if err != nil {
		logger.Logger.Error("Failed to get analytics overview", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}

	// Get top pages
	topPages, err := uc.analyticsRepo.GetTopPages(startDate, endDate, 10)
	if err != nil {
		logger.Logger.Error("Failed to get top pages", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}
	overview["top_pages"] = topPages

	return overview, nil
}

// Get page views with pagination
func (uc *AnalyticsUseCase) GetPageViews(ctx context.Context, startDate, endDate time.Time, path string, page, perPage int) ([]AnalyticsPageView, int64, error) {
	offset := (page - 1) * perPage

	var pageViews []AnalyticsPageView
	var total int64
	var err error

	if path != "" {
		pageViews, err = uc.analyticsRepo.GetPageViewsByPath(path, startDate, endDate, perPage, offset)
		if err != nil {
			return nil, 0, err
		}
		// Count total for this path
		err = uc.analyticsRepo.db.Model(&AnalyticsPageView{}).Where("path LIKE ? AND visited_at BETWEEN ? AND ?", "%"+path+"%", startDate, endDate).Count(&total).Error
	} else {
		pageViews, err = uc.analyticsRepo.GetPageViewsByDateRange(startDate, endDate, perPage, offset)
		if err != nil {
			return nil, 0, err
		}
		// Count total
		err = uc.analyticsRepo.db.Model(&AnalyticsPageView{}).Where("visited_at BETWEEN ? AND ?", startDate, endDate).Count(&total).Error
	}

	if err != nil {
		logger.Logger.Error("Failed to count page views", "method", "GetPageViews", "error", err)
		return nil, 0, err
	}

	return pageViews, total, nil
}

// Get visitors with pagination
func (uc *AnalyticsUseCase) GetVisitors(ctx context.Context, startDate, endDate time.Time, page, perPage int) ([]AnalyticsVisitor, int64, error) {
	offset := (page - 1) * perPage

	visitors, err := uc.analyticsRepo.GetVisitorsByDateRange(startDate, endDate, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	var total int64
	err = uc.analyticsRepo.db.Model(&AnalyticsVisitor{}).Where("last_visit BETWEEN ? AND ?", startDate, endDate).Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count visitors", "method", "GetVisitors", "error", err)
		return nil, 0, err
	}

	return visitors, total, nil
}

// Get sessions with pagination
func (uc *AnalyticsUseCase) GetSessions(ctx context.Context, startDate, endDate time.Time, page, perPage int) ([]AnalyticsSession, int64, error) {
	offset := (page - 1) * perPage

	sessions, err := uc.analyticsRepo.GetSessionsByDateRange(startDate, endDate, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	var total int64
	err = uc.analyticsRepo.db.Model(&AnalyticsSession{}).Where("started_at BETWEEN ? AND ?", startDate, endDate).Count(&total).Error
	if err != nil {
		logger.Logger.Error("Failed to count sessions", "method", "GetSessions", "error", err)
		return nil, 0, err
	}

	return sessions, total, nil
}

// Parse date range from query parameters
func (uc *AnalyticsUseCase) ParseDateRange(dateRangeStr string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error

	if dateRangeStr == "" {
		// Default to last 30 days
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -30)
	} else {
		// Parse date range (e.g., "7", "30", "90")
		days, err := utils.ParseInt(dateRangeStr)
		if err != nil || days == nil || *days <= 0 || *days > 365 {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid date range: %s", dateRangeStr)
		}
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -*days)
	}

	return startDate, endDate, err
}
