package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAnalyticsRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewAnalyticsHandler()

	// Analytics overview endpoint (temporarily public for testing)
	r.GET("/analytics/overview", h.GetAnalyticsOverview).Register()

	// Page views endpoint (temporarily public for testing)
	r.GET("/analytics/page-views", h.GetPageViews).Register()

	// Visitors endpoint (temporarily public for testing)
	r.GET("/analytics/visitors", h.GetVisitors).Register()

	// Sessions endpoint (temporarily public for testing)
	r.GET("/analytics/sessions", h.GetSessions).Register()
}
