package data_loader

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/db"
	p "github.com/easy-comerce/backend/internal/permission"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/constants"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

func InitRoleAndSuperAdmin() error {
	showDeleted := true
	db := db.GetDB()
	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	permissionRepo := p.NewPermissionRepository(db)
	permissionService := p.NewPermissionService(permissionRepo)

	allPermissionNames := getAllPermissionNames()
	if err := createAllPermissionsIfNotExists(permissionService, allPermissionNames, &showDeleted); err != nil {
		return err
	}

	if err := createUserRoleIfNotExists(roleRepo, permissionService, &showDeleted); err != nil {
		return err
	}

	if err := createSuperAdminRoleAndUserIfNotExists(permissionService, roleRepo, userRepo, allPermissionNames, &showDeleted); err != nil {
		return err
	}

	return nil
}

func createSuperAdminRoleAndUserIfNotExists(ps p.PermissionService, rr role.RoleRepository, ur user.UserRepository, allPermissionNames []string, showDeleted *bool) error {
	superAdminEmail, superAdminPassword, err := getSuperAdminCredentials()
	if err != nil {
		return err
	}

	if superAdminEmail == "" || superAdminPassword == "" {
		logger.Error("❌ Super admin email or password is not set", "method", "createSuperAdminRoleAndUserIfNotExists")
		return errors.New("super admin email or password is getting null/empty from env")
	}

	superAdminRole, err := createSuperAdminRoleIfNotExists(rr, ps, allPermissionNames, showDeleted)
	if err != nil {
		return err
	}

	exists, err := ur.UserExistsByEmailAndRoleId(superAdminEmail, superAdminRole.ID, showDeleted)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("❌ Failed to check if super admin exists", "method", "createSuperAdminRoleAndUserIfNotExists", "error", err, "email", superAdminEmail)
		return err
	}

	if exists {
		logger.Info("Super Admin user ALREADY EXISTS", "method", "createSuperAdminRoleAndUserIfNotExists", "email", superAdminEmail)
		logger.Info("Super Admin user ALREADY EXISTS", "method", "createSuperAdminRoleAndUserIfNotExists", "email", superAdminEmail)
		return nil
	}

	if err := createSuperAdminUser(ur, superAdminEmail, superAdminPassword, superAdminRole); err != nil {
		return err
	}

	logger.Info("Super Admin user CREATED SUCCESSFULLY", "method", "createSuperAdminUserIfNotExists", "email", superAdminEmail, "type", c.RoleTypeSuperAdmin)
	return nil
}

func getSuperAdminCredentials() (string, string, error) {
	cfg := config.GetConfig()
	superAdminEmail := cfg.AppConfig.SuperAdminEmail
	superAdminPassword := cfg.AppConfig.SuperAdminPassword

	err, hashedPassword := utils.HashPassword(superAdminPassword)
	if err != nil || hashedPassword == "" {
		logger.Error("❌ Failed to hash password", "method", "getSuperAdminCredentials", "error", err, "email", superAdminEmail)
		return "", "", err
	}

	superAdminPassword = hashedPassword
	if superAdminEmail == "" || superAdminPassword == "" {
		logger.Error("❌ Super admin email or password is not set", "method", "getSuperAdminCredentials")
		return "", "", errors.New("super admin email or password is getting null/empty from env")
	}

	return superAdminEmail, superAdminPassword, nil
}

func createAllPermissionsIfNotExists(permissionService p.PermissionService, allPermissionNames []string, showDeleted *bool) error {
	err := permissionService.CreatePermissionsIfNotExists(context.Background(), allPermissionNames, showDeleted)
	if err != nil {
		logger.Error("❌ Failed to create all permissions if not exists", "method", "createAllPermissionsIfNotExists", "error", err)
		return err
	}
	return nil
}

func createSuperAdminUser(userRepo user.UserRepository, email, password string, superAdminRole *role.RoleEntity) error {
	userModel := &user.UserEntity{
		Email:    email,
		Password: password,
		Name:     c.RoleNameSuperAdmin,
		Verified: true,
		Banned:   false,
		Roles:    []role.RoleEntity{*superAdminRole},
	}

	_, err := userRepo.CreateUser(userModel)
	if err != nil {
		logger.Error("❌ Failed to create super admin", "method", "createSuperAdminUser", "error", err)
		return err
	}
	return nil
}

func createSuperAdminRoleIfNotExists(roleRepo role.RoleRepository, permissionService p.PermissionService, allPermissionNames []string, showDeleted *bool) (*role.RoleEntity, error) {
	superAdminRole, err := roleRepo.GetRoleByType(c.RoleTypeSuperAdmin, showDeleted)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("❌ Failed to get super admin role", "method", "createSuperAdminRoleIfNotExists", "error", err)
		return nil, err
	}

	if superAdminRole == nil {
		allPermissions, err := permissionService.GetAllPermissions(context.Background(), "", "", showDeleted)
		if err != nil {
			logger.Error("❌ Failed to load all permissions to create super admin role", "method", "createSuperAdminRoleIfNotExists", "error", err)
			return nil, err
		}

		superAdminRole = &role.RoleEntity{
			RoleName:    c.RoleNameSuperAdmin,
			RoleType:    c.RoleTypeSuperAdmin,
			Description: nil,
			Permissions: allPermissions,
		}

		superAdminRole, err = roleRepo.CreateRole(superAdminRole)
		if err != nil {
			logger.Error("❌ Failed to create super admin role", "method", "createSuperAdminRoleIfNotExists", "error", err)
			return nil, err
		}

		logger.Info("Super Admin role CREATED SUCCESSFULLY", "method", "createSuperAdminRoleIfNotExists", "type", c.RoleTypeSuperAdmin)
	} else {
		err = roleRepo.AppendPermissionsToRole(superAdminRole.ID, allPermissionNames, showDeleted)
		if err != nil {
			logger.Error("❌ Failed to add permissions to super admin role", "method", "createSuperAdminRoleIfNotExists", "error", err)
			return nil, err
		}

		logger.Info("Super Admin role ALREADY EXISTS", "method", "createSuperAdminRoleIfNotExists", "type", c.RoleTypeSuperAdmin)
	}

	return superAdminRole, nil
}

func createUserRoleIfNotExists(roleRepo role.RoleRepository, permissionService p.PermissionService, showDeleted *bool) error {
	userRole, err := roleRepo.GetRoleByType(c.RoleTypeUser, showDeleted)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("❌ Failed to get user role", "method", "createUserRoleIfNotExists", "error", err)
		return err
	}

	if userRole == nil {
		permission, err := permissionService.GetPermissionByName(context.Background(), c.PermissionGeneralUser)
		if err != nil {
			logger.Error("❌ Failed to load permission to create user role", "method", "createUserRoleIfNotExists", "error", err)
			return err
		}

		userRole = &role.RoleEntity{
			RoleName:    c.RoleNameUser,
			RoleType:    c.RoleTypeUser,
			Description: nil,
			Permissions: []p.PermissionEntity{*permission},
		}

		userRole, err = roleRepo.CreateRole(userRole)
		if err != nil {
			logger.Error("❌ Failed to create user role", "method", "createUserRoleIfNotExists", "error", err)
			return err
		}

		logger.Info("User role CREATED SUCCESSFULLY", "method", "createUserRoleIfNotExists", "type", c.RoleTypeUser)
	} else {
		err = roleRepo.AppendPermissionsToRole(userRole.ID, []string{c.PermissionGeneralUser}, showDeleted)
		if err != nil {
			logger.Error("❌ Failed to add permissions to user role", "method", "createUserRoleIfNotExists", "error", err)
			return err
		}

		logger.Info("User role ALREADY EXISTS", "method", "createUserRoleIfNotExists", "type", c.RoleTypeUser)
	}

	return nil
}

func getAllPermissionNames() []string {
	allPermissionNames := []string{}

	for permissionName := range constants.PermissionMap {
		allPermissionNames = append(allPermissionNames, permissionName)
	}
	return allPermissionNames
}
