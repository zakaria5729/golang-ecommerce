package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
)

type UserUseCase struct {
	userRepo *UserRepository
	roleRepo *role.RoleRepository
}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{
		userRepo: NewUserRepository(),
		roleRepo: role.NewRoleRepository(),
	}
}

func (uc *UserUseCase) UpdateProfile(userID uint, req *UpdateProfileRequest) error {
	user := &User{
		Name:    req.Name,
		PathKey: req.PathKey,
	}

	user.Sanitize()
	err := uc.userRepo.UpdateNameAndPathKey(userID, user.Name, user.PathKey)
	if err != nil {
		logger.Logger.Error("Failed to update user profile", "method", "UpdateUserProfile", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

func (uc *UserUseCase) GetAllUsersPaginated(includeStr string, showDeletedStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)

	users, total, err := uc.userRepo.GetAllUsersPaginated(include, showDeleted, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch users paginated", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return utils.BuildPaginatedResponse(uc.getUserResponses(users), total, page, pageSize), nil
}

// **REQUIRED
func (uc *UserUseCase) GetUserByID(userID uint, includeStr string, showDeletedStr string) (*UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	showDeleted := utils.ParseBoolPtr(showDeletedStr)
	user, err := uc.userRepo.GetUserByID(userID, include, showDeleted)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserByID", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

func (uc *UserUseCase) GetAuthUserByID(userID uint, includeRoles bool, includePermissions bool) (*User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := uc.userRepo.GetAuthUserByID(userID, includeRoles, includePermissions)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserByID", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *UserUseCase) GetAuthUserStatusByID(userID uint) (bool, bool, error) {
	if userID <= 0 {
		return false, false, errors.New("invalid user ID")
	}

	banned, verified, err := uc.userRepo.GetAuthUserStatusByID(userID)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetAuthUserStatusByID", "error", err, "userID", userID)
		return false, false, errors.New("user not found")
	}
	return banned, verified, nil
}

func (uc *UserUseCase) DeleteUser(userID uint) error {
	exists, err := uc.userRepo.UserExists(userID, nil)
	if err != nil || !exists {
		logger.Logger.Error("User not found", "method", "DeleteUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	if err := uc.userRepo.DeleteUser(userID); err != nil {
		logger.Logger.Error("Failed to delete user", "method", "DeleteUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logger.Logger.Info("User deleted successfully", "method", "DeleteUser", "userID", userID)
	return nil
}

func (uc *UserUseCase) UndoDeletedUser(userID uint) error {
	showDeleted := true
	exists, err := uc.userRepo.UserExists(userID, &showDeleted)
	if err != nil || !exists {
		logger.Logger.Error("User not found", "method", "DeleteUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	if err := uc.userRepo.UndoDeletedUser(userID); err != nil {
		logger.Logger.Error("Failed to delete user", "method", "DeleteUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logger.Logger.Info("User deleted successfully", "method", "DeleteUser", "userID", userID)
	return nil
}

func (uc *UserUseCase) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.IsUserExists(req.Email)
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "CreateUser", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		logger.Logger.Error("User already exists", "method", "CreateUser", "email", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	user := &User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Verified: req.Verified,
		Banned:   req.Banned,
	}

	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash password", "method", "CreateUser", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	fetchedRole, err := uc.roleRepo.GetRoleByID(req.RoleID, []string{constants.RolePermissions}, nil)
	if err != nil {
		logger.Logger.Error("Failed to get default role", "method", "CreateUser", "error", err, "email", req.Email)
		return nil, fmt.Errorf("no role found with this roleID: %w", err)
	}

	user.Roles = []role.Role{*fetchedRole}
	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to create user", "method", "CreateUser", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ToResponse(), nil
}

func (uc *UserUseCase) UpdateUser(id uint, req *UpdateUserRequest) error {
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.UserExists(id, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "UpdateUser", "error", err, "id", id)
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		logger.Logger.Error("User already exists", "method", "UpdateUser", "id", id)
		return errors.New("user with this email already exists")
	}

	user := &User{
		Password: req.Password,
		Name:     req.Name,
		Verified: req.Verified,
		Banned:   req.Banned,
	}

	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash password", "method", "Register", "error", err, "id", id)
		return fmt.Errorf("failed to process password: %w", err)
	}

	err = uc.userRepo.UpdateUserInfo(id, user)
	if err != nil {
		logger.Logger.Error("Failed to update user", "method", "Register", "error", err, "id", id)
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) UpdateUserPurchaseCountAndTotalSpent(userID uint, purchaseCount uint, totalSpent uint) error {
	err := uc.userRepo.UpdateUserPurchaseCountAndTotalSpent(userID, purchaseCount, totalSpent)
	if err != nil {
		logger.Logger.Error("Failed to update user purchase count and total spent", "method", "UpdateUserPurchaseCountAndTotalSpent", "error", err, "userID", userID, "purchaseCount", purchaseCount, "totalSpent", totalSpent)
		return fmt.Errorf("failed to update user purchase count and total spent: %w", err)
	}
	return nil
}

func (uc *UserUseCase) getUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}
	return responses
}
