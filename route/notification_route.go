package route

import (
	"github.com/easy-comerce/backend/db"
	n "github.com/easy-comerce/backend/internal/notification"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterNotificationRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	fcmService := n.NewFCMService()
	repo := n.NewNotificationRepository(db.GetDB())
	service := n.NewNotificationService(repo, fcmService)
	h := n.NewNotificationHandler(service)

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
