package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterAddressRoute(mux *http.ServeMux, permissionMiddleware *middleware.PermissionMiddleware) {
	handler := handler.NewAddressHandler()

	mux.Handle(constants.GET+" /v1/addresses", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllAddresses)))

	mux.Handle(constants.GET+" /v1/addresses/paginated", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllAddressesPaginated)))

	mux.Handle(constants.GET+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAddressByID)))

	mux.Handle(constants.POST+" /v1/addresses", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.CreateAddress)))

	mux.Handle(constants.PUT+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.UpdateAddress)))

	mux.Handle(constants.DELETE+" /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.DeleteAddress)))

	mux.Handle(constants.PATCH+" /v1/addresses/id/{id}/set-default", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.SetDefaultAddress)))

	mux.Handle(constants.GET+" /v1/addresses/default", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetDefaultAddress)))
}
