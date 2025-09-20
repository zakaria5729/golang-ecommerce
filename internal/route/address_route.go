package route

import (
	"github.com/easy-comerce/backend/internal/handler"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAddressRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	h := handler.NewAddressHandler()

	r.GET("/addresses/me", h.GetAllAddressesByUser).Use(
		pm.RequireAuthUserStatus(),
	)

	r.GET("/addresses/paginated", h.GetAllAddressesPaginated).Use(
		pm.RequirePermission(c.PermissionAddressRead),
	)

	r.GET("/addresses/id/{id}", h.GetAddressByID).Use(
		pm.RequirePermission(c.PermissionAddressRead),
	)

	r.POST("/addresses", h.CreateAddress).Use(
		pm.RequireAuthUserStatus(),
	)

	r.PUT("/addresses/id/{id}", h.UpdateAddress).Use(
		pm.RequireAuthUserStatus(),
	)

	r.DELETE("/addresses/id/{id}", h.DeleteAddress).Use(
		pm.RequireAuthUser(),
	)

	r.PATCH("/addresses/id/{id}/set-default", h.SetDefaultAddress).Use(
		pm.RequireAuthUserStatus(),
	)

	r.GET("/addresses/default", h.GetDefaultAddress).Use(
		pm.RequireAuthUserStatus(),
	)
}

// func RegisterAddressRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
// 	handler := handler.NewAddressHandler()

// 	mux.Handle(constants.GET+" /v1/addresses", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.GetAllAddresses)))

// 	mux.Handle(constants.GET+" /v1/addresses/paginated", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.GetAllAddressesPaginated)))

// 	mux.Handle(constants.GET+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.GetAddressByID)))

// 	mux.Handle(constants.POST+" /v1/addresses", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.CreateAddress)))

// 	mux.Handle(constants.PUT+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.UpdateAddress)))

// 	mux.Handle(constants.DELETE+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.DeleteAddress)))

// 	mux.Handle(constants.PATCH+" /v1/addresses/id/{id}/set-default", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.SetDefaultAddress)))

// 	mux.Handle(constants.GET+" /v1/addresses/default", middleware.ChainMiddleware(
// 		permissionMiddleware.RequireAuth(),
// 	)(http.HandlerFunc(handler.GetDefaultAddress)))
// }
