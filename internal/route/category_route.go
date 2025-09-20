package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterCategoryRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewCategoryHandler()

	r.GET("/categories", h.GetAllCategories)

	r.GET("/categories/id/{id}", h.GetCategoryByID)

	r.GET("/categories/paginated", h.GetAllCategoriesPaginated)

	r.POST("/categories", h.CreateCategory).Use(
		pm.RequirePermission(c.PermissionCategoryCreate),
	)

	r.PUT("/categories/id/{id}", h.UpdateCategory).Use(
		pm.RequirePermission(c.PermissionCategoryUpdate),
	)

	r.DELETE("/categories/id/{id}", h.DeleteCategory).Use(
		pm.RequirePermission(c.PermissionCategoryDelete),
	)

	r.PATCH("/categories/id/{id}/toggle", h.ToggleCategoryStatus).Use(
		pm.RequirePermission(c.PermissionCategoryUpdate),
	)
}

// func RegisterCategoryRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	handler := handler.NewCategoryHandler()

// 	mux.HandleFunc(constants.GET+" /v1/categories", handler.GetAllCategories)
// 	mux.HandleFunc(constants.GET+" /v1/categories/paginated", handler.GetAllCategoriesPaginated)
// 	mux.HandleFunc(constants.GET+" /v1/categories/id/{id}", handler.GetCategoryByID)

// 	mux.Handle(constants.POST+" /v1/categories", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionCategoryCreate),
// 	)(http.HandlerFunc(handler.CreateCategory)))

// 	mux.Handle(constants.PUT+" /v1/categories/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionCategoryUpdate),
// 	)(http.HandlerFunc(handler.UpdateCategory)))

// 	mux.Handle(constants.DELETE+" /v1/categories/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionCategoryDelete),
// 	)(http.HandlerFunc(handler.DeleteCategory)))

// 	mux.Handle(constants.PATCH+" /v1/categories/id/{id}/toggle", middleware.ChainMiddleware(
// 		permissionMiddleware.RequirePermission(constants.PermissionCategoryUpdate),
// 	)(http.HandlerFunc(handler.ToggleCategoryStatus)))
// }
