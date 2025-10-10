package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterNotificationRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewNotificationHandler()

	r.GET("/notifications", h.GetNotifications).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PUT("/notifications/{id}/read", h.MarkAsRead).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PUT("/notifications/read-all", h.MarkAllAsRead).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.GET("/notifications/unread-count", h.GetUnreadCount).Use(
		pm.RequireAuthUserStatus(),
	).Register()
}
