package notification

import (
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterNotificationRoutes(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := NewNotificationRepository()
	fcmService := NewFCMService()
	useCase := NewNotificationUseCase(repo, fcmService)
	handler := NewNotificationHandler(useCase)

	r.GET("/notifications", handler.GetNotifications).Use(pm.RequireAuthUser()).Register()
	r.PUT("/notifications/{id}/read", handler.MarkAsRead).Use(pm.RequireAuthUser()).Register()
	r.PUT("/notifications/read-all", handler.MarkAllAsRead).Use(pm.RequireAuthUser()).Register()
	r.GET("/notifications/unread-count", handler.GetUnreadCount).Use(pm.RequireAuthUser()).Register()
	r.POST("/notifications/send-push", handler.SendPush).Use(pm.RequireAuthUser()).Register()
}
