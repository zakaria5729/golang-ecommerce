package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/role"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type UserService struct {
	userRepo *UserRepository
	roleRepo *role.RoleRepository
}

func NewUserService(userRepo *UserRepository, roleRepo *role.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (uc *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) error {
	user := &User{
		Name: req.Name,
	}

	if req.PathKey != nil && *req.PathKey != "" {
		user.PathKey = req.PathKey
	}

	user.Sanitize()
	err := uc.userRepo.UpdateNameAndPathKey(userID, user.Name, user.PathKey)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (uc *UserService) ChangePassword(userID uint, password string, req *ChangePasswordRequest) error {
	if err := uc.userRepo.UpdateUserPassword(userID, password, req); err != nil {
		return err
	}

	return nil
}

func (uc *UserService) GetAllUsersPaginated(includeStr string, showDeletedStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	users, total, err := uc.userRepo.GetAllUsersPaginated(include, showDeleted, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return utils.BuildPaginatedResponse(getUserResponses(users), total, page, pageSize), nil
}

func (uc *UserService) GetUserByID(userID uint, includeStr string, showDeletedStr string) (*UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)
	user, err := uc.userRepo.GetUserByID(userID, include, showDeleted)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

func (uc *UserService) GetAuthUserByID(userID uint, includeRoles bool, includePermissions bool) (*User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := uc.userRepo.GetAuthUserByID(userID, includeRoles, includePermissions)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *UserService) GetAuthUserStatusByID(userID uint) (bool, bool, *string, error) {
	if userID <= 0 {
		return false, false, nil, errors.New("invalid user ID")
	}

	banned, verified, refreshToken, err := uc.userRepo.GetAuthUserStatusByID(userID)
	if err != nil {
		return false, false, nil, errors.New("user not found")
	}

	return banned, verified, refreshToken, nil
}

func (uc *UserService) DeleteUser(ctx context.Context, userID uint) error {
	exists, err := uc.userRepo.UserExists(userID, nil)
	if err != nil || !exists {
		return errors.New("user not found")
	}

	if err := uc.userRepo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (uc *UserService) UndoDeletedUser(ctx context.Context, userID uint) error {
	showDeleted := true
	exists, err := uc.userRepo.UserExists(userID, &showDeleted)
	if err != nil || !exists {
		return errors.New("user not found")
	}

	if err := uc.userRepo.UndoDeletedUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (uc *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.IsUserExists(req.Email)
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	user := &User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Verified: req.Verified,
		Banned:   req.Banned,
	}

	userID, ok := ctx.Value(c.UserIDContextKey).(uint)
	if ok {
		user.CreatedBy = &userID
	}

	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	fetchedRole, err := uc.roleRepo.GetRoleByID(req.RoleID, []string{c.RolePermissions}, nil)
	if err != nil {
		return nil, fmt.Errorf("no role found with this roleID: %w", err)
	}

	user.Roles = []role.Role{*fetchedRole}
	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ToResponse(), nil
}

func (uc *UserService) UpdateUser(authUserID uint, updateUserID uint, req *UpdateUserRequest) error {
	req.Name = utils.Trim(req.Name)
	superAdminEmail := os.Getenv(c.EnvSuperAdminEmail)

	email, err := uc.userRepo.GetUserEmail(updateUserID, nil)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if email == nil {
		return errors.New("user not found with this id")
	}

	if (authUserID == updateUserID && !req.Verified) || (authUserID == updateUserID && req.Banned) {
		return errors.New("cannot update self banned or unverified status")
	}

	if superAdminEmail != "" && *email == superAdminEmail {
		return errors.New("cannot update super admin user")
	}

	user := &User{
		Name:     req.Name,
		Verified: req.Verified,
		Banned:   req.Banned,
	}
	user.UpdatedBy = &authUserID

	err = uc.userRepo.UpdateUserInfo(updateUserID, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (uc *UserService) UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error {
	err := uc.userRepo.UpdateUserPurchaseCountAndTotalSpent(userID, purchaseCount, totalSpent)
	if err != nil {
		return fmt.Errorf("failed to update user purchase count and total spent: %w", err)
	}
	return nil
}

func getUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}
	return responses
}
