package route

import (
	at "github.com/easy-comerce/backend/internal/attribute_type"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAttributeTypeRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	repo := at.NewAttributeTypeRepository()
	service := at.NewAttributeTypeService(repo)
	h := at.NewAttributeTypeHandler(service)

	r.GET("/attribute-types/public", h.GetAllAttributeTypesPublic).Register()

	r.GET("/attribute-types/{id}/public", h.GetAttributeTypeByIDPublic).Register()

	r.GET("/attribute-types", h.GetAllAttributeTypes).Use(
		pm.RequirePermission(c.PermissionAttributeTypeRead),
	).Register()

	r.GET("/attribute-types/{id}", h.GetAttributeTypeByID).Use(
		pm.RequirePermission(c.PermissionAttributeTypeRead),
	).Register()

	r.POST("/attribute-types", h.CreateAttributeType).Use(
		pm.RequirePermission(c.PermissionAttributeTypeCreate),
	).Register()

	r.PUT("/attribute-types/{id}", h.UpdateAttributeType).Use(
		pm.RequirePermission(c.PermissionAttributeTypeUpdate),
	).Register()

	r.DELETE("/attribute-types/{id}", h.DeleteAttributeType).Use(
		pm.RequirePermission(c.PermissionAttributeTypeDelete),
	).Register()

	r.POST("/attribute-types/{id}/undo", h.UndoDeletedAttributeType).Use(
		pm.RequirePermission(c.PermissionAttributeTypeUndoDelete),
	).Register()
}
