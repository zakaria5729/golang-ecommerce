package route

import (
	"net/http"

	"github.com/easy-comerce/backend/internal/handler"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
)

const (
	versionV1 = "/v1"
)

func RegisterRoutes(mux *http.ServeMux) {
	cfg := config.Load()
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	registerAuthRoutes(mux, authMiddleware)
	registerCategoryRoutes(mux, authMiddleware)
	registerAddressRoutes(mux, authMiddleware)
	registerBrowsingHistoryRoutes(mux, authMiddleware)
	registerReviewRoutes(mux, authMiddleware)
	registerWishlistRoutes(mux, authMiddleware)
}

func registerAuthRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	cfg := config.Load()
	handler := handler.NewAuthHandler(cfg.JWTSecret)

	mux.HandleFunc("POST "+versionV1+"/auth/login", handler.Login)
	mux.HandleFunc("POST "+versionV1+"/auth/register", handler.Register)
	mux.HandleFunc("POST "+versionV1+"/auth/refresh-token", handler.RefreshToken)
	mux.HandleFunc("POST "+versionV1+"/auth/logout", authMiddleware.RequireAuth(handler.Logout))
	mux.HandleFunc("POST "+versionV1+"/auth/forgot-password", handler.ForgotPassword)
	mux.HandleFunc("POST "+versionV1+"/auth/reset-password", handler.ResetPassword)

	mux.HandleFunc("GET "+versionV1+"/auth/me", authMiddleware.RequireAuth(handler.GetMe))
	mux.HandleFunc("GET "+versionV1+"/auth/profile", authMiddleware.RequireAuth(handler.GetProfile))
	mux.HandleFunc("PUT "+versionV1+"/auth/profile", authMiddleware.RequireAuth(handler.UpdateProfile))
	mux.HandleFunc("POST "+versionV1+"/auth/change-password", authMiddleware.RequireAuth(handler.ChangePassword))

	mux.HandleFunc("GET "+versionV1+"/auth/roles",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionRoleRead),
		)(handler.GetAllRoles))

	mux.HandleFunc("GET "+versionV1+"/auth/roles/id",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionRoleRead),
		)(handler.GetRoleByID))

	mux.HandleFunc("POST "+versionV1+"/auth/roles",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionRoleCreate),
		)(handler.CreateRole))

	mux.HandleFunc("PUT "+versionV1+"/auth/roles/id",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionRoleUpdate),
		)(handler.UpdateRole))

	mux.HandleFunc("DELETE "+versionV1+"/auth/roles/id",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionRoleDelete),
		)(handler.DeleteRole))

	mux.HandleFunc("POST "+versionV1+"/auth/roles/assign",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionUserUpdate),
		)(handler.AssignRoleToUser))

	mux.HandleFunc("GET "+versionV1+"/auth/permissions",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionPermissionRead),
		)(handler.GetAllPermissions))

	mux.HandleFunc("GET "+versionV1+"/auth/permissions/id",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionPermissionRead),
		)(handler.GetPermissionByID))
}

func registerCategoryRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	handler := handler.NewCategoryHandler()

	// Public endpoints - no authentication required
	mux.HandleFunc("GET "+versionV1+"/categories", handler.GetAllCategories)
	mux.HandleFunc("GET "+versionV1+"/categories/paginated", handler.GetAllCategoriesPaginated)
	mux.HandleFunc("GET "+versionV1+"/categories/id/{id}", handler.GetCategoryByID)

	// Protected endpoints - require authentication and permissions
	mux.HandleFunc("POST "+versionV1+"/categories",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionCategoryCreate),
		)(handler.CreateCategory))

	mux.HandleFunc("PUT "+versionV1+"/categories/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionCategoryUpdate),
		)(handler.UpdateCategory))

	mux.HandleFunc("DELETE "+versionV1+"/categories/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionCategoryDelete),
		)(handler.DeleteCategory))

	mux.HandleFunc("PATCH "+versionV1+"/categories/id/{id}/toggle",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionCategoryUpdate),
		)(handler.ToggleCategoryStatus))
}

func registerAddressRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	handler := handler.NewAddressHandler()

	// All address endpoints require authentication (user-specific data)
	mux.HandleFunc("GET "+versionV1+"/addresses",
		authMiddleware.RequireAuth(handler.GetAllAddresses))

	mux.HandleFunc("GET "+versionV1+"/addresses/paginated",
		authMiddleware.RequireAuth(handler.GetAllAddressesPaginated))

	mux.HandleFunc("GET "+versionV1+"/addresses/id/{id}",
		authMiddleware.RequireAuth(handler.GetAddressByID))

	mux.HandleFunc("POST "+versionV1+"/addresses",
		authMiddleware.RequireAuth(handler.CreateAddress))

	mux.HandleFunc("PUT "+versionV1+"/addresses/id/{id}",
		authMiddleware.RequireAuth(handler.UpdateAddress))

	mux.HandleFunc("DELETE "+versionV1+"/addresses/id/{id}",
		authMiddleware.RequireAuth(handler.DeleteAddress))

	mux.HandleFunc("PATCH "+versionV1+"/addresses/id/{id}/set-default",
		authMiddleware.RequireAuth(handler.SetDefaultAddress))

	mux.HandleFunc("GET "+versionV1+"/addresses/default",
		authMiddleware.RequireAuth(handler.GetDefaultAddress))
}

func registerBrowsingHistoryRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	handler := handler.NewBrowsingHistoryHandler()

	mux.HandleFunc("GET "+versionV1+"/browsing-history",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryRead),
		)(handler.GetAllBrowsingHistory))

	mux.HandleFunc("GET "+versionV1+"/browsing-history/paginated",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryRead),
		)(handler.GetAllBrowsingHistoryPaginated))

	mux.HandleFunc("GET "+versionV1+"/browsing-history/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryRead),
		)(handler.GetBrowsingHistoryByID))

	mux.HandleFunc("POST "+versionV1+"/browsing-history",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryCreate),
		)(handler.CreateBrowsingHistory))

	mux.HandleFunc("DELETE "+versionV1+"/browsing-history/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryDelete),
		)(handler.DeleteBrowsingHistory))

	mux.HandleFunc("DELETE "+versionV1+"/browsing-history/clear",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryDelete),
		)(handler.ClearBrowsingHistory))

	mux.HandleFunc("GET "+versionV1+"/browsing-history/recent",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryRead),
		)(handler.GetRecentBrowsingHistory))

	mux.HandleFunc("GET "+versionV1+"/browsing-history/most-viewed",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionBrowsingHistoryRead),
		)(handler.GetMostViewedProducts))
}

func registerReviewRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	handler := handler.NewReviewHandler()

	// Public endpoints - anyone can read reviews
	mux.HandleFunc("GET "+versionV1+"/products/{id}/reviews", handler.GetReviewsByProduct)
	mux.HandleFunc("GET "+versionV1+"/products/{id}/rating-stats", handler.GetProductRatingStats)

	// Protected endpoints - require authentication
	mux.HandleFunc("GET "+versionV1+"/reviews",
		authMiddleware.RequireAuth(handler.GetAllReviews))

	mux.HandleFunc("GET "+versionV1+"/reviews/paginated",
		authMiddleware.RequireAuth(handler.GetAllReviewsPaginated))

	mux.HandleFunc("GET "+versionV1+"/reviews/id/{id}",
		authMiddleware.RequireAuth(handler.GetReviewByID))

	mux.HandleFunc("GET "+versionV1+"/users/{id}/reviews",
		authMiddleware.RequireAuth(handler.GetReviewsByUser))

	// Write operations - require authentication
	mux.HandleFunc("POST "+versionV1+"/reviews",
		authMiddleware.RequireAuth(handler.CreateReview))

	mux.HandleFunc("PUT "+versionV1+"/reviews/id/{id}",
		authMiddleware.RequireAuth(handler.UpdateReview))

	mux.HandleFunc("DELETE "+versionV1+"/reviews/id/{id}",
		authMiddleware.RequireAuth(handler.DeleteReview))
}

func registerWishlistRoutes(mux *http.ServeMux, authMiddleware *middleware.AuthMiddleware) {
	handler := handler.NewWishlistHandler()

	mux.HandleFunc("GET "+versionV1+"/wishlists",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.GetAllWishlists))

	mux.HandleFunc("GET "+versionV1+"/wishlists/paginated",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.GetAllWishlistsPaginated))

	mux.HandleFunc("GET "+versionV1+"/wishlists/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.GetWishlistByID))

	mux.HandleFunc("POST "+versionV1+"/wishlists",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistCreate),
		)(handler.CreateWishlist))

	mux.HandleFunc("DELETE "+versionV1+"/wishlists/id/{id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistDelete),
		)(handler.DeleteWishlist))

	mux.HandleFunc("DELETE "+versionV1+"/wishlists/clear",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistDelete),
		)(handler.ClearWishlist))

	mux.HandleFunc("GET "+versionV1+"/wishlists/count",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.GetWishlistCount))

	mux.HandleFunc("GET "+versionV1+"/wishlists/product/{product_id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.GetWishlistByProduct))

	mux.HandleFunc("DELETE "+versionV1+"/wishlists/product/{product_id}",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistDelete),
		)(handler.DeleteWishlistByProduct))

	mux.HandleFunc("GET "+versionV1+"/wishlists/product/{product_id}/check",
		middleware.ChainAuthMiddleware(
			authMiddleware.RequireAuth,
			authMiddleware.RequireAnyPermission(constants.PermissionWishlistRead),
		)(handler.IsProductInWishlist))
}
