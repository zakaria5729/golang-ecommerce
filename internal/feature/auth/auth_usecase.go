package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	userRepo       *user.UserRepository
	roleRepo       *role.RoleRepository
	permissionRepo *permission.PermissionRepository
	jwtSecret      string
}

func NewAuthUseCase(jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       user.NewUserRepository(),
		roleRepo:       role.NewRoleRepository(),
		permissionRepo: permission.NewPermissionRepository(),
		jwtSecret:      jwtSecret,
	}
}

func (uc *AuthUseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	user, err := uc.userRepo.GetUserByEmail(req.Email, []string{user.UserRoles, user.UserPermissions, user.UserPassword})
	if err != nil {
		logger.Logger.Error("User not found", "method", "Login", "error", err, "email", req.Email)
		return nil, errors.New("invalid email or password")
	}

	if user.Banned {
		logger.Logger.Error("User is banned", "method", "Login", "userID", user.ID, "email", req.Email)
		return nil, errors.New("account is banned")
	}

	if !user.CheckPassword(req.Password) {
		logger.Logger.Error("Invalid password", "method", "Login", "userID", user.ID, "email", req.Email)
		return nil, errors.New("invalid email or password")
	}

	accessToken, expiresAt, err := uc.generateJWT(user)
	if err != nil {
		logger.Logger.Error("Failed to generate JWT", "method", "Login", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := uc.generateRefreshToken()
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", "method", "Login", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := uc.userRepo.SetRefreshToken(user.ID, refreshToken, refreshExpiresAt); err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "Login", "error", err, "userID", user.ID)
	}

	if lastLoginAt, err := uc.userRepo.UpdateLastLogin(user.ID); err != nil {
		logger.Logger.Error("Failed to update last login", "method", "Login", "error", err, "userID", user.ID)
	} else {
		user.LastLoginAt = &lastLoginAt
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresAt:    expiresAt,
	}, nil
}

func (uc *AuthUseCase) Register(req *RegisterRequest) (*user.UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.UserExistsByEmail(req.Email, nil)
	if err != nil {
		logger.Logger.Error("Failed to check if user exists", "method", "Register", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		logger.Logger.Error("User already exists", "method", "Register", "email", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	user := &user.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Verified: false,
		Banned:   false,
	}

	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash password", "method", "Register", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	userRole, err := uc.roleRepo.GetRoleByType(constants.RoleTypeUser, nil)
	if err != nil {
		logger.Logger.Error("Failed to get default role", "method", "Register", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	user.Roles = []role.Role{*userRole}
	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to create user", "method", "Register", "error", err, "email", req.Email)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ToResponse(), nil
}

func (uc *AuthUseCase) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := uc.userRepo.GetUserByID(userID, []string{user.UserPassword})
	if err != nil {
		logger.Logger.Error("User not found", "method", "ChangePassword", "error", err, "userID", userID)
		return errors.New("user not found")
	}

	if !user.CheckPassword(req.CurrentPassword) {
		logger.Logger.Error("Invalid current password", "method", "ChangePassword", "userID", userID)
		return errors.New("invalid current password")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash new password", "method", "ChangePassword", "error", err, "userID", userID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	_, err = uc.userRepo.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to update password", "method", "ChangePassword", "error", err, "userID", userID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) ForgotPassword(req *ForgotPasswordRequest) error {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	user, err := uc.userRepo.GetUserByEmail(req.Email, nil)
	if err != nil {
		logger.Logger.Error("User not found for password reset", "method", "ForgotPassword", "error", err, "email", req.Email)
		return nil
	}

	token, err := uc.generatePasswordResetToken()
	if err != nil {
		logger.Logger.Error("Failed to generate password reset token", "method", "ForgotPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := time.Now().Add(constants.PasswordResetTokenExpiry)
	if err := uc.userRepo.SetPasswordResetToken(user.ID, token, expiresAt); err != nil {
		logger.Logger.Error("Failed to set password reset token", "method", "ForgotPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) ResetPassword(req *ResetPasswordRequest) error {
	user, err := uc.userRepo.GetUserByPasswordResetToken(req.Token, nil)
	if err != nil {
		logger.Logger.Error("Invalid or expired password reset token", "method", "ResetPassword", "error", err, "token", req.Token)
		return errors.New("invalid or expired reset token")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash new password", "method", "ResetPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	_, err = uc.userRepo.UpdateUser(user)
	if err != nil {
		logger.Logger.Error("Failed to update password", "method", "ResetPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := uc.userRepo.ClearPasswordResetToken(user.ID); err != nil {
		logger.Logger.Error("Failed to clear password reset token", "method", "ResetPassword", "error", err, "userID", user.ID)
	}

	return nil
}

func (uc *AuthUseCase) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (uc *AuthUseCase) generateJWT(user *user.User) (string, int64, error) {
	expirationTime := time.Now().Add(constants.AccessTokenExpiry)
	expiresAt := expirationTime.Unix()

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.RoleType)
	}

	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Roles:    roleNames,
		Username: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    constants.ProjectName,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

func (uc *AuthUseCase) generatePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (uc *AuthUseCase) HasPermission(userID uint, permissionName string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if permissionName == "" {
		return false, errors.New("permission cannot be empty")
	}

	return uc.permissionRepo.HasPermission(userID, permissionName)
}

func (uc *AuthUseCase) HasRole(userID uint, roleType string) (bool, error) {
	if userID == 0 {
		return false, errors.New("invalid user ID")
	}

	if roleType == "" {
		return false, errors.New("role type cannot be empty")
	}

	if permissions, err := uc.permissionRepo.GetUserPermissionsByRole(userID, roleType); err != nil {
		return false, err
	} else {
		return len(permissions) > 0, nil
	}
}

func (uc *AuthUseCase) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	user, err := uc.userRepo.GetUserByRefreshToken(req.RefreshToken, []string{user.UserRoles, user.UserPermissions})
	if err != nil {
		logger.Logger.Error("Invalid or expired refresh token", "method", "RefreshToken", "error", err, "token", req.RefreshToken)
		return nil, errors.New("invalid or expired refresh token")
	}

	if user.Banned {
		logger.Logger.Error("User is banned", "method", "RefreshToken", "userID", user.ID, "email", user.Email)
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := uc.generateJWT(user)
	if err != nil {
		logger.Logger.Error("Failed to generate JWT", "method", "RefreshToken", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := uc.generateRefreshToken()
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", "method", "RefreshToken", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := uc.userRepo.SetRefreshToken(user.ID, refreshToken, refreshExpiresAt); err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "RefreshToken", "error", err, "userID", user.ID)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresAt:    expiresAt,
	}, nil
}

func (uc *AuthUseCase) Logout(userID uint) error {
	if err := uc.userRepo.ClearRefreshToken(userID); err != nil {
		logger.Logger.Error("Failed to clear refresh token", "method", "Logout", "error", err, "userID", userID)
		return err
	}
	return nil
}

func (uc *AuthUseCase) generateRefreshToken() (string, time.Time, error) {
	token, err := uc.generatePasswordResetToken()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(constants.RefreshTokenExpiry)
	return token, expiresAt, nil
}
