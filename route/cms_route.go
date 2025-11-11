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

	r.GET("/cms/{tag}", h.GetPageByTag).Register()

	r.GET("/cms/page/{id}", h.GetPageByID).Use(
		pm.RequirePermission(c.PermissionCmsRead),
	).Register()

	r.GET("/cms/page/paginated", h.GetPagesPaginated).Use(
		pm.RequirePermission(c.PermissionCmsRead),
	).Register()

	r.POST("/cms/page", h.CreatePageSection).Use(
		pm.RequirePermission(c.PermissionCmsCreate),
	).Register()

	r.PUT("/cms/page/{id}", h.UpdatePageSection).Use(
		pm.RequirePermission(c.PermissionCmsUpdate),
	).Register()

	r.DELETE("/cms/page/{id}", h.DeletePageSection).Use(
		pm.RequirePermission(c.PermissionCmsDelete),
	).Register()

	r.POST("/cms/page/{id}/undo", h.UndoDeletePageSection).Use(
		pm.RequirePermission(c.PermissionCmsDelete),
	).Register()
}
