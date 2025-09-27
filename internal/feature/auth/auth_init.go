package auth

import (
	"errors"
	"os"

	p "github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
)

func InitializeDefaultSuperAdmin() error {
	showDeleted := true
	userRepo := user.NewUserRepository()
	roleRepo := role.NewRoleRepository()
	permissionRepo := p.NewPermissionRepository()

	if err := createSuperAdminUserIfNotExists(permissionRepo, roleRepo, userRepo, &showDeleted); err != nil {
		return err
	}

	// superAdminEmail, superAdminPassword, err := getSuperAdminCredentials()
	// if err != nil {
	// 	return err
	// }

	// allPermissionNames := getAllPermissionNames()
	// if err := createAllPermissionsIfNotExists(permissionRepo, allPermissionNames, &showDeleted); err != nil {
	// 	return err
	// }

	// superAdminRole, err := createSuperAdminRoleIfNotExists(roleRepo, permissionRepo, allPermissionNames, &showDeleted)
	// if err != nil {
	// 	return err
	// }

	// exists, err := userRepo.UserExistsByEmailAndRoleId(superAdminEmail, superAdminRole.ID, &showDeleted)
	// if err != nil {
	// 	logger.Logger.Error("Failed to check if super admin exists", "method", "InitializeDefaultSuperAdmin", "error", err, "email", superAdminEmail)
	// 	return err
	// }

	// if exists {
	// 	logger.Logger.Info("Super admin ALREADY EXISTS", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail)
	// 	return nil
	// }

	// if err := createSuperAdminUser(userRepo, superAdminEmail, superAdminPassword, superAdminRole); err != nil {
	// 	return err
	// }

	// logger.Logger.Info("Super admin CREATED SUCCESSFULLY", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail, "name", c.RoleNameSuperAdmin)
	// return nil
	return nil
}

func createSuperAdminUserIfNotExists(pr *p.PermissionRepository, rr *role.RoleRepository, ur *user.UserRepository, showDeleted *bool) error {
	superAdminEmail, superAdminPassword, err := getSuperAdminCredentials()
	if err != nil {
		return err
	}

	allPermissionNames := getAllPermissionNames()
	if err := createAllPermissionsIfNotExists(pr, allPermissionNames, showDeleted); err != nil {
		return err
	}

	superAdminRole, err := createSuperAdminRoleIfNotExists(rr, pr, allPermissionNames, showDeleted)
	if err != nil {
		return err
	}

	exists, err := ur.UserExistsByEmailAndRoleId(superAdminEmail, superAdminRole.ID, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to check if super admin exists", "method", "InitializeDefaultSuperAdmin", "error", err, "email", superAdminEmail)
		return err
	}

	if exists {
		logger.Logger.Info("Super admin ALREADY EXISTS", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail)
		return nil
	}

	if err := createSuperAdminUser(ur, superAdminEmail, superAdminPassword, superAdminRole); err != nil {
		return err
	}

	logger.Logger.Info("Super admin CREATED SUCCESSFULLY", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail, "name", c.RoleNameSuperAdmin)
	return nil
}

func createUserRoleIfNotExists(pr *p.PermissionRepository, rr *role.RoleRepository, ur *user.UserRepository, showDeleted *bool) error {

	userRole, err := rr.GetRoleByType(c.RoleTypeUser, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get user role", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}

	if userRole == nil {
		allPermissionNames := getAllPermissionNames()
		if err := createAllPermissionsIfNotExists(pr, allPermissionNames, showDeleted); err != nil {
			return err
		}

	}

	superAdminEmail, superAdminPassword, err := getSuperAdminCredentials()
	if err != nil {
		return err
	}

	allPermissionNames := getAllPermissionNames()
	if err := createAllPermissionsIfNotExists(pr, allPermissionNames, showDeleted); err != nil {
		return err
	}

	superAdminRole, err := createSuperAdminRoleIfNotExists(rr, pr, allPermissionNames, showDeleted)
	if err != nil {
		return err
	}

	exists, err := ur.UserExistsByEmailAndRoleId(superAdminEmail, superAdminRole.ID, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to check if super admin exists", "method", "InitializeDefaultSuperAdmin", "error", err, "email", superAdminEmail)
		return err
	}

	if exists {
		logger.Logger.Info("Super admin ALREADY EXISTS", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail)
		return nil
	}

	if err := createSuperAdminUser(ur, superAdminEmail, superAdminPassword, superAdminRole); err != nil {
		return err
	}

	logger.Logger.Info("Super admin CREATED SUCCESSFULLY", "method", "InitializeDefaultSuperAdmin", "email", superAdminEmail, "name", c.RoleNameSuperAdmin)
	return nil
}

func getSuperAdminCredentials() (string, string, error) {
	superAdminEmail := os.Getenv(c.EnvSuperAdminEmail)
	superAdminPassword := os.Getenv(c.EnvSuperAdminPassword)

	if superAdminEmail == "" || superAdminPassword == "" {
		logger.Logger.Error("Super admin email or password is not set", "method", "InitializeDefaultSuperAdmin")
		return "", "", errors.New("super admin email or password is getting null/empty from env")
	}
	return superAdminEmail, superAdminPassword, nil
}

func createAllPermissionsIfNotExists(permissionRepo *p.PermissionRepository, allPermissionNames []string, showDeleted *bool) error {
	err := permissionRepo.CreatePermissionsIfNotExists(allPermissionNames, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to create all permissions if not exists", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}
	return nil
}

func createSuperAdminRoleIfNotExists(roleRepo *role.RoleRepository, permissionRepo *p.PermissionRepository, allPermissionNames []string, showDeleted *bool) (*role.Role, error) {
	superAdminRole, err := roleRepo.GetRoleByType(c.RoleTypeSuperAdmin, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
		return nil, err
	}

	if superAdminRole == nil {
		allPermissions, err := permissionRepo.GetAllPermissions("", "", showDeleted)
		if err != nil {
			logger.Logger.Error("Failed to loadd all permissions to create super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
			return nil, err
		}

		superAdminRole = &role.Role{
			RoleName:    c.RoleNameSuperAdmin,
			RoleType:    c.RoleTypeSuperAdmin,
			Description: nil,
			Permissions: allPermissions,
		}

		superAdminRole, err = roleRepo.CreateRole(superAdminRole)
		if err != nil {
			logger.Logger.Error("Failed to create super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
			return nil, err
		}
	} else {
		err = roleRepo.AddPermissionsToRole(superAdminRole.ID, allPermissionNames, showDeleted)
		if err != nil {
			logger.Logger.Error("Failed to add permissions to super admin role", "method", "InitializeDefaultSuperAdmin", "error", err)
			return nil, err
		}
	}
	return superAdminRole, nil
}

func createSuperAdminUser(userRepo *user.UserRepository, email, password string, superAdminRole *role.Role) error {
	userModel := &user.User{
		Email:    email,
		Password: password,
		Name:     c.RoleNameSuperAdmin,
		Verified: true,
		Banned:   false,
		Roles:    []role.Role{*superAdminRole},
	}

	_, err := userRepo.CreateUser(userModel)
	if err != nil {
		logger.Logger.Error("Failed to create super admin", "method", "InitializeDefaultSuperAdmin", "error", err)
		return err
	}
	return nil
}

func getAllPermissionNames() []string {
	allPermissions := []string{}
	for permissionName := range constants.PermissionMap {
		allPermissions = append(allPermissions, permissionName)
	}
	return allPermissions
}
