package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo  *user.UserRepository
	roleRepo  *role.RoleRepository
	jwtSecret string
}

func NewAuthService(
	jwtSecret string,
	userRepo *user.UserRepository,
	roleRepo *role.RoleRepository,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	user, err := s.userRepo.GetFullUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.Verified {
		return nil, errors.New("account is not verified yet")
	}

	if user.Banned {
		return nil, errors.New("account is banned")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if lastLoginAt, err := s.userRepo.SetRefreshTokenAndLastLoginAt(user.ID, refreshToken, refreshExpiresAt); err != nil {
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

func (s *AuthService) Register(req *RegisterRequest) (*user.UserResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Name = utils.Trim(req.Name)

	exists, err := s.userRepo.IsUserExists(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
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
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	userRole, err := s.roleRepo.GetRoleWithPermissionsByType(c.RoleTypeUser)
	if err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	user.Roles = []role.Role{*userRole}
	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ToResponse(), nil
}

func (s *AuthService) ForgotPassword(req *ForgotPasswordRequest) error {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	userID, err := s.userRepo.GetUserIdByEmail(req.Email)
	if err != nil || userID == nil {
		return nil
	}

	token, err := tokenutil.GenerateNewToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := timeutil.AddHoursUTC(c.PasswordResetTokenExpiryHours)
	if err := s.userRepo.SetPasswordResetToken(*userID, token, expiresAt); err != nil {
		return fmt.Errorf("failed to set reset token: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(req *ResetPasswordRequest) error {
	user, err := s.userRepo.GetUserByResetPasswordToken(req.Token)
	if err != nil || user == nil || user.PasswordResetExpires.Before(timeutil.NowUTC()) {
		return errors.New("invalid or expired reset token")
	}

	if user.Email == config.GetConfig().SuperAdminEmail {
		return errors.New("can not reset super admin password")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return fmt.Errorf("failed to process new password: %w", err)
	}

	err = s.userRepo.ResetPassword(user.ID, user.Password)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if !user.Verified {
		return nil, errors.New("account is not verified yet")
	}

	if user.Banned {
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userRepo.SetRefreshToken(user.ID, &refreshToken, &refreshExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to set refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *AuthService) Logout(userID uint) error {
	if err := s.userRepo.SetRefreshToken(userID, nil, nil); err != nil {
		return fmt.Errorf("failed to clear refresh token: %w", err)
	}
	return nil
}

func (s *AuthService) HealthCheck(ctx context.Context) *HealthResponse {
	dbStatus := "healthy"

	sqlDB, err := db.GetDB().DB()
	if err != nil {
		dbStatus = "unhealthy"
	} else {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(pingCtx); err != nil {
			dbStatus = "unhealthy"
		}
	}

	return &HealthResponse{
		ServerStatus: "healthy",
		DBStatus:     dbStatus,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}
