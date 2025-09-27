package analytics

import (
	"time"

	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

// AnalyticsVisitor represents a unique visitor
type AnalyticsVisitor struct {
	models.BaseModel
	VisitorID      string    `gorm:"not null;unique;column:visitor_id"`
	FirstVisit     time.Time `gorm:"not null;column:first_visit"`
	LastVisit      time.Time `gorm:"not null;column:last_visit"`
	TotalVisits    int       `gorm:"default:1;column:total_visits"`
	TotalPageViews int       `gorm:"default:0;column:total_page_views"`
	UniquePages    int       `gorm:"default:0;column:unique_pages"`
	Country        *string   `gorm:"column:country"`
	City           *string   `gorm:"column:city"`
	Region         *string   `gorm:"column:region"`
	Timezone       *string   `gorm:"column:timezone"`
	IsBot          bool      `gorm:"default:false;column:is_bot"`
	BotName        *string   `gorm:"column:bot_name"`
}

func (AnalyticsVisitor) TableName() string {
	return "analytics_visitors"
}

func (v *AnalyticsVisitor) Sanitize() {
	if v.VisitorID != "" {
		v.VisitorID = utils.Trim(v.VisitorID)
	}
	if v.Country != nil && *v.Country != "" {
		sanitized := utils.Trim(*v.Country)
		v.Country = &sanitized
	}
	if v.City != nil && *v.City != "" {
		sanitized := utils.Trim(*v.City)
		v.City = &sanitized
	}
	if v.Region != nil && *v.Region != "" {
		sanitized := utils.Trim(*v.Region)
		v.Region = &sanitized
	}
	if v.Timezone != nil && *v.Timezone != "" {
		sanitized := utils.Trim(*v.Timezone)
		v.Timezone = &sanitized
	}
	if v.BotName != nil && *v.BotName != "" {
		sanitized := utils.Trim(*v.BotName)
		v.BotName = &sanitized
	}
}

// AnalyticsPageView represents a single page view
type AnalyticsPageView struct {
	models.BaseModel
	VisitorID    string    `gorm:"not null;column:visitor_id"`
	SessionID    string    `gorm:"not null;column:session_id"`
	UserID       *uint     `gorm:"column:user_id"`
	Path         string    `gorm:"not null;column:path"`
	PageTitle    *string   `gorm:"column:page_title"`
	HTTPMethod   string    `gorm:"default:'GET';column:http_method"`
	StatusCode   int       `gorm:"default:200;column:status_code"`
	ResponseTime int       `gorm:"default:0;column:response_time"`
	IPAddress    *string   `gorm:"column:ip_address"`
	UserAgent    *string   `gorm:"column:user_agent"`
	Referrer     *string   `gorm:"column:referrer"`
	Country      *string   `gorm:"column:country"`
	City         *string   `gorm:"column:city"`
	Region       *string   `gorm:"column:region"`
	Timezone     *string   `gorm:"column:timezone"`
	IsBot        bool      `gorm:"default:false;column:is_bot"`
	BotName      *string   `gorm:"column:bot_name"`
	VisitedAt    time.Time `gorm:"not null;column:visited_at"`
}

func (AnalyticsPageView) TableName() string {
	return "analytics_page_views"
}

func (pv *AnalyticsPageView) Sanitize() {
	if pv.VisitorID != "" {
		pv.VisitorID = utils.Trim(pv.VisitorID)
	}
	if pv.SessionID != "" {
		pv.SessionID = utils.Trim(pv.SessionID)
	}
	if pv.Path != "" {
		pv.Path = utils.Trim(pv.Path)
	}
	if pv.PageTitle != nil && *pv.PageTitle != "" {
		sanitized := utils.Trim(*pv.PageTitle)
		pv.PageTitle = &sanitized
	}
	if pv.HTTPMethod != "" {
		pv.HTTPMethod = utils.Trim(pv.HTTPMethod)
	}
	if pv.IPAddress != nil && *pv.IPAddress != "" {
		sanitized := utils.Trim(*pv.IPAddress)
		pv.IPAddress = &sanitized
	}
	if pv.UserAgent != nil && *pv.UserAgent != "" {
		sanitized := utils.Trim(*pv.UserAgent)
		pv.UserAgent = &sanitized
	}
	if pv.Referrer != nil && *pv.Referrer != "" {
		sanitized := utils.Trim(*pv.Referrer)
		pv.Referrer = &sanitized
	}
	if pv.Country != nil && *pv.Country != "" {
		sanitized := utils.Trim(*pv.Country)
		pv.Country = &sanitized
	}
	if pv.City != nil && *pv.City != "" {
		sanitized := utils.Trim(*pv.City)
		pv.City = &sanitized
	}
	if pv.Region != nil && *pv.Region != "" {
		sanitized := utils.Trim(*pv.Region)
		pv.Region = &sanitized
	}
	if pv.Timezone != nil && *pv.Timezone != "" {
		sanitized := utils.Trim(*pv.Timezone)
		pv.Timezone = &sanitized
	}
	if pv.BotName != nil && *pv.BotName != "" {
		sanitized := utils.Trim(*pv.BotName)
		pv.BotName = &sanitized
	}
}

// AnalyticsSession represents a user session
type AnalyticsSession struct {
	models.BaseModel
	SessionID string     `gorm:"not null;unique;column:session_id"`
	VisitorID string     `gorm:"not null;column:visitor_id"`
	UserID    *uint      `gorm:"column:user_id"`
	StartedAt time.Time  `gorm:"not null;column:started_at"`
	EndedAt   *time.Time `gorm:"column:ended_at"`
	Duration  int        `gorm:"default:0;column:duration"`
	PageViews int        `gorm:"default:0;column:page_views"`
	IsBounce  bool       `gorm:"default:true;column:is_bounce"`
	Referrer  *string    `gorm:"column:referrer"`
	EntryPage *string    `gorm:"column:entry_page"`
	ExitPage  *string    `gorm:"column:exit_page"`
}

func (AnalyticsSession) TableName() string {
	return "analytics_sessions"
}

func (s *AnalyticsSession) Sanitize() {
	if s.SessionID != "" {
		s.SessionID = utils.Trim(s.SessionID)
	}
	if s.VisitorID != "" {
		s.VisitorID = utils.Trim(s.VisitorID)
	}
	if s.Referrer != nil && *s.Referrer != "" {
		sanitized := utils.Trim(*s.Referrer)
		s.Referrer = &sanitized
	}
	if s.EntryPage != nil && *s.EntryPage != "" {
		sanitized := utils.Trim(*s.EntryPage)
		s.EntryPage = &sanitized
	}
	if s.ExitPage != nil && *s.ExitPage != "" {
		sanitized := utils.Trim(*s.ExitPage)
		s.ExitPage = &sanitized
	}
}
