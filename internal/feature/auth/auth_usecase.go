package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/easy-comerce/backend/internal/feature/permission"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
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

// **REQUIRED
func (uc *AuthUseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	user, err := uc.userRepo.GetFullUserByEmail(req.Email)
	if err != nil {
		logger.Logger.Error("User not found", "method", "Login", "error", err, "email", req.Email)
		return nil, errors.New("invalid email or password")
	}

	if !user.Verified {
		logger.Logger.Error("User is not verified yet", "method", "Login", "userID", user.ID, "email", req.Email)
		return nil, errors.New("account is not verified yet")
	}

	if user.Banned {
		logger.Logger.Error("User is banned", "method", "Login", "userID", user.ID, "email", req.Email)
		return nil, errors.New("account is banned")
	}

	if !user.CheckPassword(req.Password) {
		logger.Logger.Error("Invalid password", "method", "Login", "userID", user.ID, "email", req.Email)
		return nil, errors.New("invalid email or password")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, uc.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to generate JWT", "method", "Login", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewRefreshToken()
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", "method", "Login", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if lastLoginAt, err := uc.userRepo.SetRefreshTokenAndLastLoginAt(user.ID, refreshToken, refreshExpiresAt); err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "Login", "error", err, "userID", user.ID)
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

// **REQUIRED
func (uc *AuthUseCase) Register(req *RegisterRequest) (*user.UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := uc.userRepo.IsUserExists(req.Email)
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

	userRole, err := uc.roleRepo.GetRoleWithPermissionsByType(constants.RoleTypeUser)
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

// **REQUIRED
func (uc *AuthUseCase) ChangePassword(user user.User, req *ChangePasswordRequest) error {
	if !user.CheckPassword(req.CurrentPassword) {
		logger.Logger.Error("Invalid current password", "method", "ChangePassword", "userID", user.ID)
		return errors.New("invalid current password")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash new password", "method", "ChangePassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	if err := uc.userRepo.UpdateUserPassword(user.ID, user.Password); err != nil {
		logger.Logger.Error("Failed to update password", "method", "ChangePassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *AuthUseCase) ForgotPassword(req *ForgotPasswordRequest) error {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	userID, err := uc.userRepo.GetUserIdByEmail(req.Email)
	if err != nil || userID == nil {
		logger.Logger.Error("User not found for password reset", "method", "ForgotPassword", "error", err, "email", req.Email)
		return nil
	}

	token, err := tokenutil.GenerateNewToken()
	if err != nil {
		logger.Logger.Error("Failed to generate password reset token", "method", "ForgotPassword", "error", err, "userID", userID)
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := timeutil.AddHoursUTC(constants.PasswordResetTokenExpiryHours)
	if err := uc.userRepo.SetPasswordResetToken(*userID, token, expiresAt); err != nil {
		logger.Logger.Error("Failed to set password reset token", "method", "ForgotPassword", "error", err, "userID", userID)
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *AuthUseCase) ResetPassword(req *ResetPasswordRequest) error {
	user, err := uc.userRepo.GetUserByResetPasswordToken(req.Token)
	if err != nil || user == nil || user.PasswordResetExpires.Before(timeutil.NowUTC()) {
		logger.Logger.Error("Invalid or expired password reset token", "method", "ResetPassword", "error", err, "token", req.Token)
		return errors.New("invalid or expired reset token")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		logger.Logger.Error("Failed to hash new password", "method", "ResetPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	err = uc.userRepo.UpdatePasswordAndClearResetPasswordToken(user.ID, user.Password)
	if err != nil {
		logger.Logger.Error("Failed to update password", "method", "ResetPassword", "error", err, "userID", user.ID)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// **REQUIRED
func (uc *AuthUseCase) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	user, err := uc.userRepo.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		logger.Logger.Error("Invalid or expired refresh token", "method", "RefreshToken", "error", err, "token", req.RefreshToken)
		return nil, errors.New("invalid or expired refresh token")
	}

	if !user.Verified {
		logger.Logger.Error("User is not verified", "method", "RefreshToken", "userID", user.ID, "email", user.Email)
		return nil, errors.New("account is not verified yet")
	}

	if user.Banned {
		logger.Logger.Error("User is banned", "method", "RefreshToken", "userID", user.ID, "email", user.Email)
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, uc.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to generate JWT", "method", "RefreshToken", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewRefreshToken()
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", "method", "RefreshToken", "error", err, "userID", user.ID)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := uc.userRepo.SetRefreshToken(user.ID, &refreshToken, &refreshExpiresAt); err != nil {
		logger.Logger.Error("Failed to set refresh token", "method", "RefreshToken", "error", err, "userID", user.ID)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// **REQUIRED
func (uc *AuthUseCase) Logout(userID uint) error {
	if err := uc.userRepo.SetRefreshToken(userID, nil, nil); err != nil {
		logger.Logger.Error("Failed to clear refresh token", "method", "Logout", "error", err, "userID", userID)
		return err
	}
	return nil
}

// func (uc *AuthUseCase) VerifyToken(tokenString string) (*JwtClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (any, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(uc.jwtSecret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*JwtClaims)
// 	if ok && token.Valid {
// 		return claims, nil
// 	}

// 	if claims.ExpiresAt.Before(timeutil.NowUTC()) {
// 		logger.Logger.Error("Token expired", "method", "VerifyToken", "error", err)
// 		return nil, errors.New("token expired")
// 	}

// 	return nil, errors.New("invalid token")
// }

// func (uc *AuthUseCase) HasPermission(userID uint, permissionName string) (bool, error) {
// 	if userID == 0 {
// 		return false, errors.New("invalid user ID")
// 	}

// 	if permissionName == "" {
// 		return false, errors.New("permission cannot be empty")
// 	}

// 	return uc.permissionRepo.HasPermission(userID, permissionName)
// }

// func (uc *AuthUseCase) HasRole(userID uint, roleType string) (bool, error) {
// 	if userID == 0 {
// 		return false, errors.New("invalid user ID")
// 	}

// 	if roleType == "" {
// 		return false, errors.New("role type cannot be empty")
// 	}

// 	if permissions, err := uc.permissionRepo.GetUserPermissionsByRole(userID, roleType); err != nil {
// 		return false, err
// 	} else {
// 		return len(permissions) > 0, nil
// 	}
// }

// // **REQUIRED
// func (uc *AuthUseCase) generateRefreshToken() (string, time.Time, error) {
// 	token, err := tokenutil.GenerateNewToken()
// 	if err != nil {
// 		return "", time.Time{}, err
// 	}

// 	expiresAt := timeutil.AddHoursUTC(constants.RefreshTokenExpiryHours)
// 	return token, expiresAt, nil
// }

// func (uc *AuthUseCase) generateJWT(user *user.User) (string, int64, error) {
// 	now := timeutil.NowUTC()
// 	expirationTime := timeutil.AddHoursUTC(constants.AccessTokenExpiryHours)
// 	expiresAt := expirationTime.Unix()

// 	var roleNames []string
// 	for _, role := range user.Roles {
// 		roleNames = append(roleNames, role.RoleType)
// 	}

// 	claims := &JwtClaims{
// 		UserID:   user.ID,
// 		Email:    user.Email,
// 		Roles:    roleNames,
// 		Username: user.Name,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			IssuedAt:  jwt.NewNumericDate(now),
// 			NotBefore: jwt.NewNumericDate(now),
// 			Issuer:    constants.ProjectName,
// 			Subject:   fmt.Sprintf("%d", user.ID),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
// 	if err != nil {
// 		return "", 0, err
// 	}

// 	return tokenString, expiresAt, nil
// }
