package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterProductStatsRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewProductStatsHandler()

	r.GET("/product-stats/paginated", h.GetAllProductStatsPaginated).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	)

	r.GET("/product-stats/id/{id}", h.GetProductStatsByID).Use(
		pm.RequirePermission(c.PermissionProductStatsRead),
	)

	r.POST("/product-stats", h.IncreaseProductStats).Use(
		pm.RequireAuthUserStatus(),
	)
}

// r.GET("/browsing-history", h.GetAllBrowsingHistory).Use(
// 	pm.RequirePermission(c.PermissionBrowsingHistoryRead),
// )

// r.DELETE("/browsing-history/id/{id}", h.DeleteBrowsingHistory).Use(
// 	pm.RequirePermission(c.PermissionBrowsingHistoryDelete),
// )

// r.DELETE("/browsing-history/clear", h.ClearBrowsingHistory).Use(
// 	pm.RequirePermission(c.PermissionBrowsingHistoryDelete),
// )

// 	r.GET("/browsing-history/recent", h.GetRecentBrowsingHistory).Use(
// 		pm.RequirePermission(c.PermissionBrowsingHistoryRead),
// 	)

// 	r.GET("/browsing-history/most-viewed", h.GetMostViewedProducts).Use(
// 		pm.RequirePermission(c.PermissionBrowsingHistoryRead),
// 	)
// }

// func RegisterBrowsingHistoryRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	handler := handler.NewBrowsingHistoryHandler()

// 	mux.Handle(constants.GET+" /v1/browsing-history", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
// 	)(http.HandlerFunc(handler.GetAllBrowsingHistory)))

// 	mux.Handle(constants.GET+" /v1/browsing-history/paginated", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
// 	)(http.HandlerFunc(handler.GetAllBrowsingHistoryPaginated)))

// 	mux.Handle(constants.GET+" /v1/browsing-history/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
// 	)(http.HandlerFunc(handler.GetBrowsingHistoryByID)))

// 	mux.Handle(constants.POST+" /v1/browsing-history", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryCreate),
// 	)(http.HandlerFunc(handler.CreateBrowsingHistory)))

// 	mux.Handle(constants.DELETE+" /v1/browsing-history/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryDelete),
// 	)(http.HandlerFunc(handler.DeleteBrowsingHistory)))

// 	mux.Handle(constants.DELETE+" /v1/browsing-history/clear", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryDelete),
// 	)(http.HandlerFunc(handler.ClearBrowsingHistory)))

// 	mux.Handle(constants.GET+" /v1/browsing-history/recent", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
// 	)(http.HandlerFunc(handler.GetRecentBrowsingHistory)))

// 	mux.Handle(constants.GET+" /v1/browsing-history/most-viewed", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
// 	)(http.HandlerFunc(handler.GetMostViewedProducts)))
// }
