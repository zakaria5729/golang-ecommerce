package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user/model"
	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type UserService interface {
	UpdateProfile(userID uint, req *model.UpdateProfileRequest) error
	ChangePassword(userID uint, req *model.ChangePasswordRequest) error
	GetAllUsersPaginated(includeStr string, showDeletedStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error)
	GetUserByID(userID uint, includeStr string, showDeletedStr string) (*model.UserResponse, error)
	GetAuthUserByID(userID uint, includeRoles bool, includePermissions bool) (*UserEntity, error)
	GetAuthUserStatusByID(userID uint) (banned bool, verified bool, refreshToken *string, err error)
	DeleteUser(ctx context.Context, userID uint) error
	UndoDeletedUser(ctx context.Context, userID uint) error
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	UpdateUser(authUserID uint, updateUserID uint, req *model.UpdateUserRequest) error
	UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error
}

type userService struct {
	userRepo UserRepository
	roleRepo role.RoleRepository
}

func NewUserService(userRepo UserRepository, roleRepo role.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) UpdateProfile(userID uint, req *model.UpdateProfileRequest) error {
	user := &UserEntity{
		Name: req.Name,
	}

	if req.PathKey != nil && *req.PathKey != "" {
		user.PathKey = req.PathKey
	}

	user.Sanitize()
	err := s.userRepo.UpdateNameAndPathKey(userID, user.Name, user.PathKey)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (s *userService) ChangePassword(userID uint, req *model.ChangePasswordRequest) error {
	password, err := s.userRepo.getUserPasswordByID(userID)
	if err != nil {
		return errors.New("Failed to verify user current password")
	}

	if err := s.userRepo.UpdateUserPassword(userID, password, req); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAllUsersPaginated(includeStr string, showDeletedStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	users, total, err := s.userRepo.GetAllUsersPaginated(include, showDeleted, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return utils.BuildPaginatedResponse(getUserResponses(users), total, page, pageSize), nil
}

func (s *userService) GetUserByID(userID uint, includeStr string, showDeletedStr string) (*model.UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)
	user, err := s.userRepo.GetUserByID(userID, include, showDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user.ToResponse(), nil
}

func (s *userService) GetAuthUserByID(userID uint, includeRoles bool, includePermissions bool) (*UserEntity, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetAuthUserByID(userID, includeRoles, includePermissions)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) GetAuthUserStatusByID(userID uint) (bool, bool, *string, error) {
	if userID <= 0 {
		return false, false, nil, errors.New("invalid user ID")
	}

	banned, verified, refreshToken, err := s.userRepo.GetAuthUserStatusByID(userID)
	if err != nil {
		return false, false, nil, errors.New("user not found")
	}

	return banned, verified, refreshToken, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	exists, err := s.userRepo.UserExists(userID, nil)
	if err != nil || !exists {
		return errors.New("user not found")
	}

	if err := s.userRepo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) UndoDeletedUser(ctx context.Context, userID uint) error {
	showDeleted := true
	exists, err := s.userRepo.UserExists(userID, &showDeleted)
	if err != nil || !exists {
		return errors.New("user not found")
	}

	if err := s.userRepo.UndoDeletedUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}

	return nil
}

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := s.userRepo.IsUserExists(req.Email)
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	user := &UserEntity{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Verified: req.Verified,
		Banned:   req.Banned,
	}

	userID, _ := cu.GetUserIDFromContext(ctx)
	user.CreatedBy = userID

	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	fetchedRole, err := s.roleRepo.GetRoleByID(req.RoleID, []string{c.RolePermissions}, nil)
	if err != nil {
		return nil, fmt.Errorf("no role found with this roleID: %w", err)
	}

	user.Roles = []role.RoleEntity{*fetchedRole}
	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ToResponse(), nil
}

func (s *userService) UpdateUser(authUserID uint, updateUserID uint, req *model.UpdateUserRequest) error {
	req.Name = utils.Trim(req.Name)
	superAdminEmail := os.Getenv(c.EnvSuperAdminEmail)

	email, err := s.userRepo.GetUserEmail(updateUserID, nil)
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

	user := &UserEntity{
		Name:     req.Name,
		Verified: req.Verified,
		Banned:   req.Banned,
	}
	user.UpdatedBy = &authUserID

	err = s.userRepo.UpdateUserInfo(updateUserID, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (uc *userService) UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error {
	err := uc.userRepo.UpdateUserPurchaseCountAndTotalSpent(userID, purchaseCount, totalSpent)
	if err != nil {
		return fmt.Errorf("failed to update user purchase count and total spent: %w", err)
	}
	return nil
}

func getUserResponses(users []UserEntity) []model.UserResponse {
	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}
	return responses
}
