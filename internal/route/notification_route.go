package route

import (
	"github.com/easy-comerce/backend/internal/feature/notification"
	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterNotificationRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := notification.NewNotificationRepository()
	fcmService := notification.NewFCMService()
	useCase := notification.NewNotificationUseCase(repo, fcmService)
	h := handler.NewNotificationHandler(useCase)

	r.GET("/notifications", h.GetNotifications).Use(pm.RequireAuthUser()).Register()
	r.PUT("/notifications/{id}/read", h.MarkAsRead).Use(pm.RequireAuthUser()).Register()
	r.PUT("/notifications/read-all", h.MarkAllAsRead).Use(pm.RequireAuthUser()).Register()
	r.GET("/notifications/unread-count", h.GetUnreadCount).Use(pm.RequireAuthUser()).Register()
	r.POST("/notifications/send-push", h.SendPush).Use(pm.RequireAuthUser()).Register()
}
