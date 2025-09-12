package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterBrowsingHistoryRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	handler := handler.NewBrowsingHistoryHandler()

	mux.Handle(constants.GET+" /v1/browsing-history", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
	)(http.HandlerFunc(handler.GetAllBrowsingHistory)))

	mux.Handle(constants.GET+" /v1/browsing-history/paginated", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
	)(http.HandlerFunc(handler.GetAllBrowsingHistoryPaginated)))

	mux.Handle(constants.GET+" /v1/browsing-history/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
	)(http.HandlerFunc(handler.GetBrowsingHistoryByID)))

	mux.Handle(constants.POST+" /v1/browsing-history", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryCreate),
	)(http.HandlerFunc(handler.CreateBrowsingHistory)))

	mux.Handle(constants.DELETE+" /v1/browsing-history/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryDelete),
	)(http.HandlerFunc(handler.DeleteBrowsingHistory)))

	mux.Handle(constants.DELETE+" /v1/browsing-history/clear", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryDelete),
	)(http.HandlerFunc(handler.ClearBrowsingHistory)))

	mux.Handle(constants.GET+" /v1/browsing-history/recent", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
	)(http.HandlerFunc(handler.GetRecentBrowsingHistory)))

	mux.Handle(constants.GET+" /v1/browsing-history/most-viewed", middleware.ChainMiddleware(
		permissionMiddleware.RequirePermission(constants.PermissionBrowsingHistoryRead),
	)(http.HandlerFunc(handler.GetMostViewedProducts)))
}
