package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/contact_info"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterContactInfoRoute(r *router.Router, pm middleware.PermissionMiddleware) {
	db := db.GetDB()
	repo := contact_info.NewContactInfoRepository(db)
	service := contact_info.NewContactInfoService(repo)
	h := contact_info.NewContactInfoHandler(service)

	r.POST("/subscribe-newsletter", h.PostNewsLetter).Register()

	r.POST("/contact-us", h.CreateContactUs).Register()
}
