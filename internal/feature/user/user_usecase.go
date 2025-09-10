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

// GetUserProfile retrieves a user profile by ID
func (uc *UserUseCase) GetUserProfile(userID uint, includeStr string) (*User, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	user, err := uc.userRepo.GetUserByID(userID, include)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetUserProfile", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}

// GetBasicUserByID retrieves basic user information by ID (without roles/permissions)
func (uc *UserUseCase) GetBasicUserByID(userID uint) (*User, error) {
	user, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "GetBasicUserByID", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}

// IsUserBanned checks if a user is banned (lightweight query - only checks banned status)
func (uc *UserUseCase) IsUserBanned(userID uint) (bool, error) {
	return uc.userRepo.IsUserBanned(userID)
}

// UpdateUserProfile updates a user's profile information
func (uc *UserUseCase) UpdateUserProfile(userID uint, req *User) (*User, error) {
	req.Sanitize()

	existingUser, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "UpdateUserProfile", "error", err, "userID", userID)
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != existingUser.Email {
		exists, err := uc.userRepo.UserExistsByEmail(req.Email, &userID)
		if err != nil {
			logger.Logger.Error("Failed to check if email exists", "method", "UpdateUserProfile", "error", err, "userID", userID, "email", req.Email)
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			logger.Logger.Error("Email already exists", "method", "UpdateUserProfile", "userID", userID, "email", req.Email)
			return nil, errors.New("email already exists")
		}
		existingUser.Email = req.Email
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}

	if err := uc.userRepo.UpdateUser(existingUser); err != nil {
		logger.Logger.Error("Failed to update user profile", "method", "UpdateUserProfile", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	existingUser.Password = ""
	return existingUser, nil
}

// GetAllUsers retrieves all users with optional filtering and sorting
func (uc *UserUseCase) GetAllUsers(includeStr, sortBy, sortOrder string) ([]User, error) {
	include := utils.ParseCommaSeparatedString(includeStr)
	users, err := uc.userRepo.GetAllUsers(include, nil, nil, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get all users", "method", "GetAllUsers", "error", err)
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
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
	if err := uc.userRepo.UpdateUser(user); err != nil {
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
	if err := uc.userRepo.UpdateUser(user); err != nil {
		logger.Logger.Error("Failed to unban user", "method", "UnbanUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to unban user: %w", err)
	}

	logger.Logger.Info("User unbanned successfully", "method", "UnbanUser", "userID", userID)
	return nil
}

// VerifyUser verifies a user's email
func (uc *UserUseCase) VerifyUser(userID uint) error {
	user, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "VerifyUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	user.Verified = true
	if err := uc.userRepo.UpdateUser(user); err != nil {
		logger.Logger.Error("Failed to verify user", "method", "VerifyUser", "error", err, "userID", userID)
		return fmt.Errorf("failed to verify user: %w", err)
	}

	logger.Logger.Info("User verified successfully", "method", "VerifyUser", "userID", userID)
	return nil
}

// AssignRoleToUser assigns a role to a user
func (uc *UserUseCase) AssignRoleToUser(userID uint, roleID uint) error {
	// Check if user exists
	_, err := uc.userRepo.GetUserByID(userID, nil)
	if err != nil {
		logger.Logger.Error("User not found", "method", "AssignRoleToUser", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	// Check if role exists
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

// GetUserStatus gets user status (banned, verified) from users table
func (uc *UserUseCase) GetUserStatus(userID uint) (banned bool, verified bool, err error) {
	if userID == 0 {
		return false, false, errors.New("invalid user ID")
	}

	return uc.userRepo.GetUserStatus(userID)
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(user *User) error {
	user.Sanitize()

	if user.Email == "" {
		logger.Logger.Error("Email is required", "method", "CreateUser")
		return errors.New("email is required")
	}

	if user.Password == "" {
		logger.Logger.Error("Password is required", "method", "CreateUser")
		return errors.New("password is required")
	}

	if user.Name == "" {
		logger.Logger.Error("Name is required", "method", "CreateUser")
		return errors.New("name is required")
	}

	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash password", "method", "CreateUser", "error", err)
		return errors.New("failed to process password")
	}

	if err := uc.userRepo.CreateUser(user); err != nil {
		logger.Logger.Error("Failed to create user", "method", "CreateUser", "error", err)
		return errors.New("failed to create user")
	}

	return nil
}

// UserExistsByEmail checks if a user exists with the given email
func (uc *UserUseCase) UserExistsByEmail(email string) (bool, error) {
	exists, err := uc.userRepo.UserExistsByEmail(email, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if user exists by email", "method", "UserExistsByEmail", "error", err, "email", email)
		return false, err
	}
	return exists, nil
}
