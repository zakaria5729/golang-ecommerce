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

// **REQUIRED
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

// **REQUIRED
func (uc *UserUseCase) GetAllUsersPaginated(includeStr string, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)
	include := utils.ParseCommaSeparatedString(includeStr)

	users, total, err := uc.userRepo.GetAllUsersPaginated(include, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to fetch users paginated", "method", "GetAllUsersPaginated", "error", err, "include", include, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return utils.BuildPaginatedResponse(uc.getUserResponses(users), total, page, pageSize), nil
}

// **REQUIRED
func (uc *UserUseCase) GetUserByID(userID uint, includeStr string) (*UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	user, err := uc.userRepo.GetUserByID(userID, include)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserByID", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

// **REQUIRED
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

// **REQUIRED
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

// **REQUIRED
func (uc *UserUseCase) DeleteUser(userID uint) error {
	exists, err := uc.userRepo.UserExists(userID)
	if err != nil || !exists {
		logger.Logger.Error("User not found", err, "method", "DeleteUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	if err := uc.userRepo.DeleteUser(userID); err != nil {
		logger.Logger.Error("Failed to delete user", "method", "DeleteUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logger.Logger.Info("User deleted successfully", "method", "DeleteUser", "userID", userID)
	return nil
}

// **REQUIRED
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

	fetchedRole, err := uc.roleRepo.GetRoleByID(req.RoleID, []string{constants.RolePermissions})
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

// **REQUIRED
func (uc *UserUseCase) UpdateUser(id uint, req *UpdateUserRequest) error {
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.UserExists(id)
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

func (uc *UserUseCase) getUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}
	return responses
}

// func (uc *UserUseCase) UserExistsByEmail(email string) (bool, error) {
// 	exists, err := uc.userRepo.UserExistsByEmail(email, nil)
// 	if err != nil {
// 		logger.Logger.Error("Failed to check if user exists by email", "method", "UserExistsByEmail", "error", err, "email", email)
// 		return false, err
// 	}
// 	return exists, nil
// }

// // BanUser bans a user account
// func (uc *UserUseCase) BanUser(userID uint) error {
// 	user, err := uc.userRepo.GetUserByID(userID, nil)
// 	if err != nil {
// 		logger.Logger.Error("User not found", "method", "BanUser", "error", err, "userID", userID)
// 		return errors.New("user not found")
// 	}

// 	user.Banned = true
// 	_, err = uc.userRepo.UpdateUser(user)
// 	if err != nil {
// 		logger.Logger.Error("Failed to ban user", "method", "BanUser", "error", err, "userID", userID)
// 		return fmt.Errorf("failed to ban user: %w", err)
// 	}

// 	logger.Logger.Info("User banned successfully", "method", "BanUser", "userID", userID)
// 	return nil
// }

// // UnbanUser unbans a user account
// func (uc *UserUseCase) UnbanUser(userID uint) error {
// 	user, err := uc.userRepo.GetUserByID(userID, nil)
// 	if err != nil {
// 		logger.Logger.Error("User not found", "method", "UnbanUser", "error", err, "userID", userID)
// 		return errors.New("user not found")
// 	}

// 	user.Banned = false
// 	_, err = uc.userRepo.UpdateUser(user)
// 	if err != nil {
// 		logger.Logger.Error("Failed to unban user", "method", "UnbanUser", "error", err, "userID", userID)
// 		return fmt.Errorf("failed to unban user: %w", err)
// 	}

// 	logger.Logger.Info("User unbanned successfully", "method", "UnbanUser", "userID", userID)
// 	return nil
// }

// func (uc *UserUseCase) VerifyUser(userID uint) error {
// 	user, err := uc.userRepo.GetUserByID(userID, nil)
// 	if err != nil {
// 		logger.Logger.Error("User not found", "method", "VerifyUser", "error", err, "userID", userID)
// 		return errors.New("user not found")
// 	}

// 	user.Verified = true
// 	_, err = uc.userRepo.UpdateUser(user)
// 	if err != nil {
// 		logger.Logger.Error("Failed to verify user", "method", "VerifyUser", "error", err, "userID", userID)
// 		return fmt.Errorf("failed to verify user: %w", err)
// 	}

// 	logger.Logger.Info("User verified successfully", "method", "VerifyUser", "userID", userID)
// 	return nil
// }

// func (uc *UserUseCase) AssignRoleToUser(userID uint, roleID uint) error {
// 	_, err := uc.userRepo.GetUserByID(userID, nil)
// 	if err != nil {
// 		logger.Logger.Error("User not found", "method", "AssignRoleToUser", "error", err, "userID", userID)
// 		return errors.New("user not found")
// 	}

// 	_, err = uc.roleRepo.GetRoleByID(roleID, nil)
// 	if err != nil {
// 		logger.Logger.Error("Role not found", "method", "AssignRoleToUser", "error", err, "roleID", roleID)
// 		return errors.New("role not found")
// 	}

// 	if err := uc.userRepo.AssignRolesToUser(userID, []uint{roleID}); err != nil {
// 		logger.Logger.Error("Failed to assign role to user", "method", "AssignRoleToUser", "error", err, "userID", userID, "roleID", roleID)
// 		return fmt.Errorf("failed to assign role: %w", err)
// 	}

// 	logger.Logger.Info("Role assigned to user successfully", "method", "AssignRoleToUser", "userID", userID, "roleID", roleID)
// 	return nil
// }

// func (uc *UserUseCase) GetUserStatus(userID uint) (banned bool, verified bool, err error) {
// 	if userID == 0 {
// 		return false, false, errors.New("invalid user ID")
// 	}

// 	return uc.userRepo.GetUserStatus(userID)
// }

// func (uc *UserUseCase) GetUserProfile(userID uint, includeStr string) (*UserResponse, error) {
// 	include := utils.ParseCommaSeparatedString(includeStr)
// 	user, err := uc.userRepo.GetUserByID(userID, include)
// 	if err != nil {
// 		logger.Logger.Error("User not found", "method", "GetUserProfile", "error", err, "userID", userID)
// 		return nil, errors.New("user not found")
// 	}

// 	return user.ToResponse(), nil
// }

// // IsUserBanned checks if a user is banned (lightweight query - only checks banned status)
// func (uc *UserUseCase) IsUserBanned(userID uint) (bool, error) {
// 	return uc.userRepo.IsUserBanned(userID)
// }
