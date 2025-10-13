package route

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
)

func RegisterUserRoute(r *router.Router, pm *middleware.PermissionMiddleware) {
	roleRepo := role.NewRoleRepository(db.GetDB())
	userRepo := user.NewUserRepository(db.GetDB())
	service := user.NewUserService(userRepo, roleRepo)
	h := user.NewUserHandler(service)

	r.GET("/users/profile", h.GetProfile).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.PUT("/users/profile", h.UpdateProfile).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.POST("/user/change-password", h.ChangePassword).Use(
		pm.RequireAuthUserStatus(),
	).Register()

	r.POST("/users", h.CreateUser).Use(
		pm.RequirePermission(c.PermissionUserCreate),
	).Register()

	r.PUT("/users/{id}", h.UpdateUser).Use(
		pm.RequirePermission(c.PermissionUserUpdate),
	).Register()

	r.GET("/users", h.GetAllUsersPaginated).Use(
		pm.RequirePermission(c.PermissionUserRead),
	).Register()

	r.GET("/users/{id}", h.GetUserByID).Use(
		pm.RequirePermission(c.PermissionUserRead),
	).Register()

	r.DELETE("/users/{id}", h.DeleteUser).Use(
		pm.RequirePermission(c.PermissionUserDelete),
	).Register()

	r.POST("/users/{id}/undo", h.UndoDeletedUser).Use(
		pm.RequirePermission(c.PermissionUserUndoDelete),
	).Register()
}
