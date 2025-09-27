package analytics

import (
	"fmt"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db.GetDB(),
	}
}

// Visitor operations
func (r *AnalyticsRepository) CreateVisitor(visitor *AnalyticsVisitor) error {
	visitor.Sanitize()
	err := r.db.Create(visitor).Error
	if err != nil {
		logger.Logger.Error("Failed to create visitor", "method", "CreateVisitor", "error", err, "visitorID", visitor.VisitorID)
	}
	return err
}

func (r *AnalyticsRepository) GetVisitorByID(visitorID string) (*AnalyticsVisitor, error) {
	var visitor AnalyticsVisitor
	err := r.db.Where("visitor_id = ?", visitorID).First(&visitor).Error
	if err != nil {
		logger.Logger.Error("Failed to get visitor", "method", "GetVisitorByID", "error", err, "visitorID", visitorID)
		return nil, err
	}
	return &visitor, nil
}

func (r *AnalyticsRepository) UpdateVisitor(visitor *AnalyticsVisitor) error {
	visitor.Sanitize()
	err := r.db.Save(visitor).Error
	if err != nil {
		logger.Logger.Error("Failed to update visitor", "method", "UpdateVisitor", "error", err, "visitorID", visitor.VisitorID)
	}
	return err
}

// Page view operations
func (r *AnalyticsRepository) CreatePageView(pageView *AnalyticsPageView) error {
	pageView.Sanitize()
	err := r.db.Create(pageView).Error
	if err != nil {
		logger.Logger.Error("Failed to create page view", "method", "CreatePageView", "error", err, "visitorID", pageView.VisitorID, "path", pageView.Path)
	}
	return err
}

func (r *AnalyticsRepository) GetPageViewsByDateRange(startDate, endDate time.Time, limit, offset int) ([]AnalyticsPageView, error) {
	var pageViews []AnalyticsPageView
	query := r.db.Where("visited_at BETWEEN ? AND ?", startDate, endDate).
		Order("visited_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&pageViews).Error
	if err != nil {
		logger.Logger.Error("Failed to get page views", "method", "GetPageViewsByDateRange", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, err
	}
	return pageViews, nil
}

func (r *AnalyticsRepository) GetPageViewsByPath(path string, startDate, endDate time.Time, limit, offset int) ([]AnalyticsPageView, error) {
	var pageViews []AnalyticsPageView
	query := r.db.Where("path LIKE ? AND visited_at BETWEEN ? AND ?", "%"+path+"%", startDate, endDate).
		Order("visited_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&pageViews).Error
	if err != nil {
		logger.Logger.Error("Failed to get page views by path", "method", "GetPageViewsByPath", "error", err, "path", path, "startDate", startDate, "endDate", endDate)
		return nil, err
	}
	return pageViews, nil
}

func (r *AnalyticsRepository) GetPageViewsByVisitor(visitorID string, limit, offset int) ([]AnalyticsPageView, error) {
	var pageViews []AnalyticsPageView
	query := r.db.Where("visitor_id = ?", visitorID).
		Order("visited_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&pageViews).Error
	if err != nil {
		logger.Logger.Error("Failed to get page views by visitor", "method", "GetPageViewsByVisitor", "error", err, "visitorID", visitorID)
		return nil, err
	}
	return pageViews, nil
}

// Session operations
func (r *AnalyticsRepository) CreateSession(session *AnalyticsSession) error {
	session.Sanitize()
	err := r.db.Create(session).Error
	if err != nil {
		logger.Logger.Error("Failed to create session", "method", "CreateSession", "error", err, "sessionID", session.SessionID)
	}
	return err
}

func (r *AnalyticsRepository) GetSessionByID(sessionID string) (*AnalyticsSession, error) {
	var session AnalyticsSession
	err := r.db.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		logger.Logger.Error("Failed to get session", "method", "GetSessionByID", "error", err, "sessionID", sessionID)
		return nil, err
	}
	return &session, nil
}

func (r *AnalyticsRepository) UpdateSession(session *AnalyticsSession) error {
	session.Sanitize()
	err := r.db.Save(session).Error
	if err != nil {
		logger.Logger.Error("Failed to update session", "method", "UpdateSession", "error", err, "sessionID", session.SessionID)
	}
	return err
}

// Analytics summary operations
func (r *AnalyticsRepository) GetAnalyticsOverview(startDate, endDate time.Time) (map[string]interface{}, error) {
	var result map[string]interface{} = make(map[string]interface{})

	// Total page views
	var totalPageViews int64
	err := r.db.Model(&AnalyticsPageView{}).Where("visited_at BETWEEN ? AND ?", startDate, endDate).Count(&totalPageViews).Error
	if err != nil {
		logger.Logger.Error("Failed to get total page views", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}
	result["total_page_views"] = totalPageViews

	// Total unique visitors
	var totalVisitors int64
	err = r.db.Model(&AnalyticsVisitor{}).Where("last_visit BETWEEN ? AND ?", startDate, endDate).Count(&totalVisitors).Error
	if err != nil {
		logger.Logger.Error("Failed to get total visitors", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}
	result["total_visitors"] = totalVisitors

	// Total sessions
	var totalSessions int64
	err = r.db.Model(&AnalyticsSession{}).Where("started_at BETWEEN ? AND ?", startDate, endDate).Count(&totalSessions).Error
	if err != nil {
		logger.Logger.Error("Failed to get total sessions", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}
	result["total_sessions"] = totalSessions

	// Bounce rate
	var bounceSessions int64
	err = r.db.Model(&AnalyticsSession{}).Where("started_at BETWEEN ? AND ? AND is_bounce = ?", startDate, endDate, true).Count(&bounceSessions).Error
	if err != nil {
		logger.Logger.Error("Failed to get bounce sessions", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}

	var bounceRate float64
	if totalSessions > 0 {
		bounceRate = float64(bounceSessions) / float64(totalSessions) * 100
	}
	result["bounce_rate"] = fmt.Sprintf("%.1f%%", bounceRate)

	// Average session duration
	var avgDuration float64
	if totalSessions > 0 {
		err = r.db.Raw("SELECT COALESCE(AVG(duration), 0) FROM analytics_sessions WHERE started_at BETWEEN ? AND ? AND duration > 0", startDate, endDate).Scan(&avgDuration).Error
		if err != nil {
			logger.Logger.Error("Failed to get average session duration", "method", "GetAnalyticsOverview", "error", err)
			return nil, err
		}
		result["average_session_duration"] = fmt.Sprintf("%.0fs", avgDuration)
	} else {
		result["average_session_duration"] = "0s"
	}

	// Top pages
	topPages, err := r.GetTopPages(startDate, endDate, 10)
	if err != nil {
		logger.Logger.Error("Failed to get top pages", "method", "GetAnalyticsOverview", "error", err)
		return nil, err
	}
	result["top_pages"] = topPages

	return result, nil
}

func (r *AnalyticsRepository) GetTopPages(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Model(&AnalyticsPageView{}).
		Select("path, COUNT(*) as views").
		Where("visited_at BETWEEN ? AND ?", startDate, endDate).
		Group("path").
		Order("views DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		logger.Logger.Error("Failed to get top pages", "method", "GetTopPages", "error", err)
		return nil, err
	}

	return results, nil
}

func (r *AnalyticsRepository) GetVisitorsByDateRange(startDate, endDate time.Time, limit, offset int) ([]AnalyticsVisitor, error) {
	var visitors []AnalyticsVisitor
	query := r.db.Where("last_visit BETWEEN ? AND ?", startDate, endDate).
		Order("last_visit DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&visitors).Error
	if err != nil {
		logger.Logger.Error("Failed to get visitors by date range", "method", "GetVisitorsByDateRange", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, err
	}
	return visitors, nil
}

func (r *AnalyticsRepository) GetSessionsByDateRange(startDate, endDate time.Time, limit, offset int) ([]AnalyticsSession, error) {
	var sessions []AnalyticsSession
	query := r.db.Where("started_at BETWEEN ? AND ?", startDate, endDate).
		Order("started_at DESC").
		Limit(limit).
		Offset(offset)

	err := query.Find(&sessions).Error
	if err != nil {
		logger.Logger.Error("Failed to get sessions by date range", "method", "GetSessionsByDateRange", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, err
	}
	return sessions, nil
}
