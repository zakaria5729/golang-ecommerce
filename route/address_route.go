package route

import (
	"github.com/easy-comerce/backend/db"
	a "github.com/easy-comerce/backend/internal/address"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAddressRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	repo := a.NewAddressRepository(db.GetDB())
	service := a.NewAddressService(repo)
	h := a.NewAddressHandler(service)

	r.GET("/addresses/me", h.GetAllAddressesByUser).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.GET("/addresses/paginated", h.GetAllAddressesPaginated).Use(
		pm.RequirePermission(c.PermissionAddressRead),
	).Register()

	r.GET("/addresses/id/{id}", h.GetAddressByID).Use(
		pm.RequirePermission(c.PermissionAddressRead),
	).Register()

	r.POST("/addresses", h.CreateAddress).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PUT("/addresses/id/{id}", h.UpdateAddress).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.POST("/addresses/id/{id}/remove", h.RemoveAddress).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PATCH("/addresses/id/{id}/set-default", h.SetDefaultAddress).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.GET("/addresses/default", h.GetDefaultAddress).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.DELETE("/addresses/id/{id}", h.DeleteAddress).Use(
		pm.RequirePermission(c.PermissionAddressDelete),
	).Register()

	r.POST("/addresses/id/{id}/undo", h.UndoDeleteAddress).Use(
		pm.RequirePermission(c.PermissionAddressUndoDelete),
	).Register()
}
