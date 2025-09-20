package auth

import (
	"os"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
)

func InitializeDefaultSuperAdmin() error {
	userRepo := user.NewUserRepository()
	roleUseCase := role.NewRoleUseCase()

	superAdminEmail := os.Getenv(constants.EnvSuperAdminEmail)
	superAdminPassword := os.Getenv(constants.EnvSuperAdminPassword)

	exists, err := userRepo.UserExistsByEmail(superAdminEmail, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if super admin exists", "method", "InitializeDefaultSuperAdmin", "error", err, "email", superAdminEmail)
		return err
	}

	if exists {
		logger.Logger.Info("Super admin already exists", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail)
		return nil
	}

	superAdminRole, err := roleUseCase.GetRoleByType(constants.RoleTypeSuperAdmin)
	if err != nil {
		logger.Logger.Error("Failed to get super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}

	superAdminName := "Super Admin"
	userModel := &user.User{
		Email:    superAdminEmail,
		Password: superAdminPassword,
		Name:     superAdminName,
		Verified: true,
		Banned:   false,
		Roles:    []role.Role{*superAdminRole},
	}

	_, err = userRepo.CreateUser(userModel)
	if err != nil {
		logger.Logger.Error("Failed to create super admin", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}

	logger.Logger.Info("Super admin created successfully", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail, "name", superAdminName)
	return nil
}
