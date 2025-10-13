package route

import (
	ao "github.com/easy-comerce/backend/internal/attribute_option"
	at "github.com/easy-comerce/backend/internal/attribute_type"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAttributeOptionRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	typeRepo := at.NewAttributeTypeRepository()
	optionRepo := ao.NewAttributeOptionRepository()
	service := ao.NewAttributeOptionService(optionRepo, typeRepo)
	h := ao.NewAttributeOptionHandler(service)

	r.GET("/attribute-options/public", h.GetAllAttributeOptionsPublic).Register()

	r.GET("/attribute-options/{id}/public", h.GetAttributeOptionByIDPublic).Register()

	r.GET("/attribute-options", h.GetAllAttributeOptions).Use(
		pm.RequirePermission(c.PermissionAttributeOptionRead),
	).Register()

	r.GET("/attribute-options/{id}", h.GetAttributeOptionByID).Use(
		pm.RequirePermission(c.PermissionAttributeOptionRead),
	).Register()

	r.POST("/attribute-options", h.CreateAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionCreate),
	).Register()

	r.PUT("/attribute-options/{id}", h.UpdateAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionUpdate),
	).Register()

	r.DELETE("/attribute-options/{id}", h.DeleteAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionDelete),
	).Register()

	r.POST("/attribute-options/{id}/undo", h.UndoDeletedAttributeOption).Use(
		pm.RequirePermission(c.PermissionAttributeOptionUndoDelete),
	).Register()
}
