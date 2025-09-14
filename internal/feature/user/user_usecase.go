package user

import (
	"errors"
	"fmt"

	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/pkg/logger"
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

func (uc *UserUseCase) GetUserProfile(userID uint, includeStr string) (*UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	user, err := uc.userRepo.GetUserByID(userID, include)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserProfile", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	return user.ToResponse(), nil
}

// IsUserBanned checks if a user is banned (lightweight query - only checks banned status)
func (uc *UserUseCase) IsUserBanned(userID uint) (bool, error) {
	return uc.userRepo.IsUserBanned(userID)
}

func (uc *UserUseCase) UpdateUserProfile(userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	existingUser, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "UpdateUserProfile", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Name != nil {
		existingUser.Name = *req.Name
	}

	existingUser.Sanitize()

	if req.Email != nil && *req.Email != existingUser.Email {
		exists, err := uc.userRepo.UserExistsByEmail(*req.Email, &userID)
		if err != nil {
			logger.Logger.Error("Failed to check if email exists", "method", "UpdateUserProfile", "error", err, "userID", userID, "email", *req.Email)
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			logger.Logger.Error("Email already exists", "method", "UpdateUserProfile", "userID", userID, "email", *req.Email)
			return nil, errors.New("email already exists")
		}
	}

	updatedUser, err := uc.userRepo.UpdateUser(existingUser)
	if err != nil {
		logger.Logger.Error("Failed to update user profile", "method", "UpdateUserProfile", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return updatedUser.ToResponse(), nil
}

func (uc *UserUseCase) GetAllUsers(includeStr, sortBy, sortOrder string) ([]UserResponse, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	users, err := uc.userRepo.GetAllUsers(include, nil, nil, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get all users", "method", "GetAllUsers", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return uc.getUserResponses(users), nil
}

// GetUserByID retrieves a specific user by ID
func (uc *UserUseCase) GetUserByID(userID uint, includeStr string) (*User, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	user, err := uc.userRepo.GetUserByID(userID, include)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserByID", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
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

// DeleteUser soft deletes a user
func (uc *UserUseCase) DeleteUser(userID uint) error {
	// Check if user exists
	_, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
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

// BanUser bans a user account
func (uc *UserUseCase) BanUser(userID uint) error {
	user, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "BanUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	user.Banned = true
	_, err = uc.userRepo.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to ban user", "method", "BanUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to ban user: %w", err)
	}

	logger.Logger.Info("User banned successfully", "method", "BanUser", "userID", userID)
	return nil
}

// UnbanUser unbans a user account
func (uc *UserUseCase) UnbanUser(userID uint) error {
	user, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "UnbanUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	user.Banned = false
	_, err = uc.userRepo.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to unban user", "method", "UnbanUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to unban user: %w", err)
	}

	logger.Logger.Info("User unbanned successfully", "method", "UnbanUser", "userID", userID)
	return nil
}

func (uc *UserUseCase) VerifyUser(userID uint) error {
	user, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "VerifyUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	user.Verified = true
	_, err = uc.userRepo.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to verify user", "method", "VerifyUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to verify user: %w", err)
	}

	logger.Logger.Info("User verified successfully", "method", "VerifyUser", "userID", userID)
	return nil
}

func (uc *UserUseCase) AssignRoleToUser(userID uint, roleID uint) error {
	_, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "AssignRoleToUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	_, err = uc.roleRepo.GetRoleByID(roleID, nil)
	if err != nil {
		logger.Logger.Error("Role not found", "method", "AssignRoleToUser", "error", err, "roleID", roleID)
		return errors.New("role not found")
	}

	if err := uc.userRepo.AssignRolesToUser(userID, []uint{roleID}); err != nil {
		logger.Logger.Error("Failed to assign role to user", "method", "AssignRoleToUser", "error", err, "userID", userID, "roleID", roleID)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	logger.Logger.Info("Role assigned to user successfully", "method", "AssignRoleToUser", "userID", userID, "roleID", roleID)
	return nil
}

func (uc *UserUseCase) GetUserStatus(userID uint) (banned bool, verified bool, err error) {
	if userID == 0 {
		return false, false, errors.New("invalid user ID")
	}

	return uc.userRepo.GetUserStatus(userID)
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	user := &User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	user.Sanitize()

	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash password", "method", "CreateUser", "error", err)
		return nil, errors.New("failed to process password")
	}

	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to create user", "method", "CreateUser", "error", err)
		return nil, errors.New("failed to create user")
	}

	return createdUser.ToResponse(), nil
}

func (uc *UserUseCase) getUserResponses(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}
	return responses
}

func (uc *UserUseCase) UserExistsByEmail(email string) (bool, error) {
	exists, err := uc.userRepo.UserExistsByEmail(email, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "UserExistsByEmail", "error", err, "email", email)
		return false, err
	}
	return exists, nil
}
