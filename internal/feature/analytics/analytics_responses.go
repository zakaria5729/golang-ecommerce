package analytics

import "time"

// AnalyticsOverviewResponse represents the analytics overview response
type AnalyticsOverviewResponse struct {
	Summary struct {
		TotalPageViews     int64  `json:"total_page_views"`
		TotalVisitors      int64  `json:"total_visitors"`
		TotalSessions      int64  `json:"total_sessions"`
		BounceRate         string `json:"bounce_rate"`
		AvgSessionDuration string `json:"average_session_duration"`
	} `json:"summary"`
	TopPages  []map[string]interface{} `json:"top_pages"`
	DateRange struct {
		Start string `json:"start"`
		End   string `json:"end"`
		Days  int    `json:"days"`
	} `json:"date_range"`
}

// PageViewResponse represents a page view response
type PageViewResponse struct {
	ID           uint      `json:"id"`
	VisitorID    string    `json:"visitor_id"`
	SessionID    string    `json:"session_id"`
	UserID       *uint     `json:"user_id,omitempty"`
	Path         string    `json:"path"`
	PageTitle    *string   `json:"page_title,omitempty"`
	HTTPMethod   string    `json:"http_method"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int       `json:"response_time"`
	IPAddress    *string   `json:"ip_address,omitempty"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	Referrer     *string   `json:"referrer,omitempty"`
	Country      *string   `json:"country,omitempty"`
	City         *string   `json:"city,omitempty"`
	Region       *string   `json:"region,omitempty"`
	Timezone     *string   `json:"timezone,omitempty"`
	IsBot        bool      `json:"is_bot"`
	BotName      *string   `json:"bot_name,omitempty"`
	VisitedAt    time.Time `json:"visited_at"`
}

// VisitorResponse represents a visitor response
type VisitorResponse struct {
	ID             uint      `json:"id"`
	VisitorID      string    `json:"visitor_id"`
	FirstVisit     time.Time `json:"first_visit"`
	LastVisit      time.Time `json:"last_visit"`
	TotalVisits    int       `json:"total_visits"`
	TotalPageViews int       `json:"total_page_views"`
	UniquePages    int       `json:"unique_pages"`
	Country        *string   `json:"country,omitempty"`
	City           *string   `json:"city,omitempty"`
	Region         *string   `json:"region,omitempty"`
	Timezone       *string   `json:"timezone,omitempty"`
	IsBot          bool      `json:"is_bot"`
	BotName        *string   `json:"bot_name,omitempty"`
}

// SessionResponse represents a session response
type SessionResponse struct {
	ID        uint       `json:"id"`
	SessionID string     `json:"session_id"`
	VisitorID string     `json:"visitor_id"`
	UserID    *uint      `json:"user_id,omitempty"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Duration  int        `json:"duration"`
	PageViews int        `json:"page_views"`
	IsBounce  bool       `json:"is_bounce"`
	Referrer  *string    `json:"referrer,omitempty"`
	EntryPage *string    `json:"entry_page,omitempty"`
	ExitPage  *string    `json:"exit_page,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	Total       int64       `json:"total"`
	LastPage    int         `json:"last_page"`
	From        int         `json:"from"`
	To          int         `json:"to"`
}

// ConvertPageViewToResponse converts AnalyticsPageView to PageViewResponse
func ConvertPageViewToResponse(pv *AnalyticsPageView) *PageViewResponse {
	return &PageViewResponse{
		ID:           pv.ID,
		VisitorID:    pv.VisitorID,
		SessionID:    pv.SessionID,
		UserID:       pv.UserID,
		Path:         pv.Path,
		PageTitle:    pv.PageTitle,
		HTTPMethod:   pv.HTTPMethod,
		StatusCode:   pv.StatusCode,
		ResponseTime: pv.ResponseTime,
		IPAddress:    pv.IPAddress,
		UserAgent:    pv.UserAgent,
		Referrer:     pv.Referrer,
		Country:      pv.Country,
		City:         pv.City,
		Region:       pv.Region,
		Timezone:     pv.Timezone,
		IsBot:        pv.IsBot,
		BotName:      pv.BotName,
		VisitedAt:    pv.VisitedAt,
	}
}

// ConvertVisitorToResponse converts AnalyticsVisitor to VisitorResponse
func ConvertVisitorToResponse(v *AnalyticsVisitor) *VisitorResponse {
	return &VisitorResponse{
		ID:             v.ID,
		VisitorID:      v.VisitorID,
		FirstVisit:     v.FirstVisit,
		LastVisit:      v.LastVisit,
		TotalVisits:    v.TotalVisits,
		TotalPageViews: v.TotalPageViews,
		UniquePages:    v.UniquePages,
		Country:        v.Country,
		City:           v.City,
		Region:         v.Region,
		Timezone:       v.Timezone,
		IsBot:          v.IsBot,
		BotName:        v.BotName,
	}
}

// ConvertSessionToResponse converts AnalyticsSession to SessionResponse
func ConvertSessionToResponse(s *AnalyticsSession) *SessionResponse {
	return &SessionResponse{
		ID:        s.ID,
		SessionID: s.SessionID,
		VisitorID: s.VisitorID,
		UserID:    s.UserID,
		StartedAt: s.StartedAt,
		EndedAt:   s.EndedAt,
		Duration:  s.Duration,
		PageViews: s.PageViews,
		IsBounce:  s.IsBounce,
		Referrer:  s.Referrer,
		EntryPage: s.EntryPage,
		ExitPage:  s.ExitPage,
	}
}
