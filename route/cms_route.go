package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/cms"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterCmsRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	cmsRepo := cms.NewCmsRepository(db.GetDB())
	cmsService := cms.NewCmsService(cmsRepo)
	h := cms.NewCmsHandler(cmsService)

	r.GET("/cms/{tag}", h.GetByTag).Register()

	r.POST("/cms", h.CreatePageSection).Use(
		pm.RequirePermission(c.PermissionCmsCreate),
	).Register()

}
