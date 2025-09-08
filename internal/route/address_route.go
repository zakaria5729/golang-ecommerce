package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func RegisterAddressRoute(mux *http.ServeMux, permissionMiddleware *middleware.AuthPermissionMiddleware) {
	handler := handler.NewAddressHandler()

	// All address endpoints require authentication (user-specific data)
	mux.Handle("GET /v1/addresses", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllAddresses)))

	mux.Handle("GET /v1/addresses/paginated", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAllAddressesPaginated)))

	mux.Handle("GET /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetAddressByID)))

	mux.Handle("POST /v1/addresses", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.CreateAddress)))

	mux.Handle("PUT /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.UpdateAddress)))

	mux.Handle("DELETE /v1/addresses/id/{id}", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.DeleteAddress)))

	mux.Handle("PATCH /v1/addresses/id/{id}/set-default", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.SetDefaultAddress)))

	mux.Handle("GET /v1/addresses/default", middleware.ChainMiddleware(
		permissionMiddleware.RequireAuth(),
	)(http.HandlerFunc(handler.GetDefaultAddress)))
}
