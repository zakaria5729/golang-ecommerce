package handler

import (
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/internal/feature/analytics"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
)

type AnalyticsHandler struct {
	analyticsUseCase *analytics.AnalyticsUseCase
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsUseCase: analytics.NewAnalyticsUseCase(),
	}
}

// GetAnalyticsOverview returns analytics overview data
func (h *AnalyticsHandler) GetAnalyticsOverview(w http.ResponseWriter, r *http.Request) {
	// Parse date range
	dateRange := r.URL.Query().Get("date_range")
	startDate, endDate, err := h.analyticsUseCase.ParseDateRange(dateRange)
	if err != nil {
		response.SendErrorJSON(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	// Get overview data
	overviewData, err := h.analyticsUseCase.GetAnalyticsOverview(r.Context(), startDate, endDate)
	if err != nil {
		logger.Logger.Error("Failed to get analytics overview", "method", "GetAnalyticsOverview", "error", err)
		response.SendErrorJSON(w, "Failed to get analytics overview", http.StatusInternalServerError)
		return
	}

	// Build response
	overview := &analytics.AnalyticsOverviewResponse{
		Summary: struct {
			TotalPageViews     int64  `json:"total_page_views"`
			TotalVisitors      int64  `json:"total_visitors"`
			TotalSessions      int64  `json:"total_sessions"`
			BounceRate         string `json:"bounce_rate"`
			AvgSessionDuration string `json:"average_session_duration"`
		}{
			TotalPageViews:     overviewData["total_page_views"].(int64),
			TotalVisitors:      overviewData["total_visitors"].(int64),
			TotalSessions:      overviewData["total_sessions"].(int64),
			BounceRate:         overviewData["bounce_rate"].(string),
			AvgSessionDuration: overviewData["average_session_duration"].(string),
		},
		TopPages: overviewData["top_pages"].([]map[string]interface{}),
		DateRange: struct {
			Start string `json:"start"`
			End   string `json:"end"`
			Days  int    `json:"days"`
		}{
			Start: startDate.Format("2006-01-02T15:04:05Z07:00"),
			End:   endDate.Format("2006-01-02T15:04:05Z07:00"),
			Days:  int(endDate.Sub(startDate).Hours() / 24),
		},
	}

	response.SendSuccessJSON(w, overview)
}

// GetPageViews returns paginated page views
func (h *AnalyticsHandler) GetPageViews(w http.ResponseWriter, r *http.Request) {
	// Parse date range
	dateRange := r.URL.Query().Get("date_range")
	startDate, endDate, err := h.analyticsUseCase.ParseDateRange(dateRange)
	if err != nil {
		response.SendErrorJSON(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page, perPage, err := h.parsePaginationParams(r)
	if err != nil {
		response.SendErrorJSON(w, "Invalid pagination parameters", http.StatusBadRequest)
		return
	}

	// Get path filter
	path := r.URL.Query().Get("path")

	// Get page views
	pageViews, total, err := h.analyticsUseCase.GetPageViews(r.Context(), startDate, endDate, path, page, perPage)
	if err != nil {
		logger.Logger.Error("Failed to get page views", "method", "GetPageViews", "error", err)
		response.SendErrorJSON(w, "Failed to get page views", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var pageViewResponses []*analytics.PageViewResponse
	for _, pv := range pageViews {
		pageViewResponses = append(pageViewResponses, analytics.ConvertPageViewToResponse(&pv))
	}

	// Build paginated response
	paginatedResponse := &analytics.PaginatedResponse{
		Data:        pageViewResponses,
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
		From:        (page-1)*perPage + 1,
		To:          int(total),
	}

	if int64(page*perPage) < total {
		paginatedResponse.To = page * perPage
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// GetVisitors returns paginated visitors
func (h *AnalyticsHandler) GetVisitors(w http.ResponseWriter, r *http.Request) {
	// Parse date range
	dateRange := r.URL.Query().Get("date_range")
	startDate, endDate, err := h.analyticsUseCase.ParseDateRange(dateRange)
	if err != nil {
		response.SendErrorJSON(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page, perPage, err := h.parsePaginationParams(r)
	if err != nil {
		response.SendErrorJSON(w, "Invalid pagination parameters", http.StatusBadRequest)
		return
	}

	// Get visitors
	visitors, total, err := h.analyticsUseCase.GetVisitors(r.Context(), startDate, endDate, page, perPage)
	if err != nil {
		logger.Logger.Error("Failed to get visitors", "method", "GetVisitors", "error", err)
		response.SendErrorJSON(w, "Failed to get visitors", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var visitorResponses []*analytics.VisitorResponse
	for _, v := range visitors {
		visitorResponses = append(visitorResponses, analytics.ConvertVisitorToResponse(&v))
	}

	// Build paginated response
	paginatedResponse := &analytics.PaginatedResponse{
		Data:        visitorResponses,
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
		From:        (page-1)*perPage + 1,
		To:          int(total),
	}

	if int64(page*perPage) < total {
		paginatedResponse.To = page * perPage
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// GetSessions returns paginated sessions
func (h *AnalyticsHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	// Parse date range
	dateRange := r.URL.Query().Get("date_range")
	startDate, endDate, err := h.analyticsUseCase.ParseDateRange(dateRange)
	if err != nil {
		response.SendErrorJSON(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page, perPage, err := h.parsePaginationParams(r)
	if err != nil {
		response.SendErrorJSON(w, "Invalid pagination parameters", http.StatusBadRequest)
		return
	}

	// Get sessions
	sessions, total, err := h.analyticsUseCase.GetSessions(r.Context(), startDate, endDate, page, perPage)
	if err != nil {
		logger.Logger.Error("Failed to get sessions", "method", "GetSessions", "error", err)
		response.SendErrorJSON(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var sessionResponses []*analytics.SessionResponse
	for _, s := range sessions {
		sessionResponses = append(sessionResponses, analytics.ConvertSessionToResponse(&s))
	}

	// Build paginated response
	paginatedResponse := &analytics.PaginatedResponse{
		Data:        sessionResponses,
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
		From:        (page-1)*perPage + 1,
		To:          int(total),
	}

	if int64(page*perPage) < total {
		paginatedResponse.To = page * perPage
	}

	response.SendSuccessJSON(w, paginatedResponse)
}

// parsePaginationParams parses pagination parameters from request
func (h *AnalyticsHandler) parsePaginationParams(r *http.Request) (int, int, error) {
	page := 1
	perPage := 50

	// Parse page parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return 0, 0, err
		}
		page = p
	}

	// Parse per_page parameter
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 || pp > 100 {
			return 0, 0, err
		}
		perPage = pp
	}

	return page, perPage, nil
}
