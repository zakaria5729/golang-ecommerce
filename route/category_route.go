package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/category"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterCategoryRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := category.NewCategoryRepository(db.GetDB())
	service := category.NewCategoryService(repo)
	h := category.NewCategoryHandler(service)

	r.GET("/categories/public", h.GetAllCategoriesPublic).Register()

	r.GET("/categories-subcategories/public", h.GetAllCategoriesWithSubcategoriesPublic).Register()

	r.GET("/categories/id/{id}/public", h.GetCategoryByIdPublic).Register()

	r.GET("/categories/paginated/public", h.GetAllCategoriesPaginatedPublic).Register()

	r.GET("/categories-subcategories", h.GetAllCategoriesWithSubcategories).Use(
		pm.RequirePermission(c.PermissionCategoryRead),
	).Register()

	r.GET("/categories", h.GetAllCategories).Use(
		pm.RequirePermission(c.PermissionCategoryRead),
	).Register()

	r.GET("/categories/id/{id}", h.GetCategoryByID).Use(
		pm.RequirePermission(c.PermissionCategoryRead),
	).Register()

	r.GET("/categories/paginated", h.GetAllCategoriesPaginated).Use(
		pm.RequirePermission(c.PermissionCategoryRead),
	).Register()

	r.POST("/categories", h.CreateCategory).Use(
		pm.RequirePermission(c.PermissionCategoryCreate),
	).Register()

	r.PUT("/categories/id/{id}", h.UpdateCategory).Use(
		pm.RequirePermission(c.PermissionCategoryUpdate),
	).Register()

	r.DELETE("/categories/id/{id}", h.DeleteCategory).Use(
		pm.RequirePermission(c.PermissionCategoryDelete),
	).Register()

	r.POST("/categories/id/{id}/undo", h.UndoDeleteCategory).Use(
		pm.RequirePermission(c.PermissionCategoryUndoDelete),
	).Register()
}
