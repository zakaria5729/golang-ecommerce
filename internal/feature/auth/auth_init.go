package auth

import (
	"os"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
)

func InitializeDefaultSuperAdmin() error {
	userUseCase := user.NewUserUseCase()
	userRepo := user.NewUserRepository()
	roleUseCase := role.NewRoleUseCase(userRepo.ToSharedInterface())

	superAdminEmail := os.Getenv("SUPER_ADMIN_EMAIL")
	if superAdminEmail == "" {
		superAdminEmail = "super123@admin.com"
	}

	superAdminPassword := os.Getenv("SUPER_ADMIN_PASSWORD")
	if superAdminPassword == "" {
		superAdminPassword = "super123@admin"
	}

	superAdminName := os.Getenv("SUPER_ADMIN_NAME")
	if superAdminName == "" {
		superAdminName = "Super Admin"
	}

	// Check if super admin already exists
	exists, err := userUseCase.UserExistsByEmail(superAdminEmail)
	if err != nil {
		logger.Logger.Error("Failed to check if super admin exists", "method", "InitializeDefaultSuperAdmin", "error", err, "email", superAdminEmail)
		return err
	}

	if exists {
		logger.Logger.Info("Super admin already exists", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail)
		return nil
	}

	superAdminRole, err := roleUseCase.GetRoleByType(constants.RoleTypeSuperAdmin, "")
	if err != nil {
		logger.Logger.Error("Failed to get super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}

	user := &user.User{
		Email:    superAdminEmail,
		Password: superAdminPassword,
		Name:     superAdminName,
		Verified: true,
		Banned:   false,
		Roles:    []role.Role{*superAdminRole},
	}

	if err := userUseCase.CreateUser(user); err != nil {
		logger.Logger.Error("Failed to create super admin", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}

	logger.Logger.Info("Super admin created successfully", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail, "name", superAdminName)
	return nil
}
