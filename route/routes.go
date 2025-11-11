package route

import (
	"github.com/easy-comerce/backend/db"
	p "github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/user"
	"github.com/easy-comerce/backend/pkg/config"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterAllRoutes(r *router.Router) {
	db := db.GetDB()
	cfg := config.GetConfig()
	userRepo := user.NewUserRepository(db)
	permissionRepo := p.NewPermissionRepository(db)
	pm := m.NewPermissionMiddleware(cfg.SecretConfig.JWTSecret, userRepo, permissionRepo)

	RegisterSystemRoute(r, pm)
	RegisterAuthRoute(r, pm)
	RegisterUserRoute(r, pm)
	RegisterRoleRoute(r, pm)
	RegisterPermissionRoute(r, pm)
	RegisterCategoryRoute(r, pm)
	RegisterAddressRoute(r, pm)
	RegisterProductStatsRoute(r, pm)
	RegisterReviewRoute(r, pm)
	RegisterWishlistRoute(r, pm)
	RegisterFileRoute(r, pm)
	RegisterBrandRoute(r, pm)
	RegisterColorRoute(r, pm)
	RegisterAttributeTypeRoute(r, pm)
	RegisterAttributeOptionRoute(r, pm)
	RegisterSizeCategoryRoute(r, pm)
	RegisterSizeOptionRoute(r, pm)
	RegisterNotificationRoute(r, pm)
	RegisterContactInfoRoute(r, pm)
	RegisterCmsRoute(r, pm)
}
