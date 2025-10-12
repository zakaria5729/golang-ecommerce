package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/file_storage"
	"github.com/easy-comerce/backend/internal/feature/role"
	"github.com/easy-comerce/backend/internal/feature/user"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService struct {
	jwtSecret string
	userRepo  *user.UserRepository
	roleRepo  *role.RoleRepository
}

func NewAuthService(
	jwtSecret string,
	userRepo *user.UserRepository,
	roleRepo *role.RoleRepository,
) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Password = utils.Trim(req.Password)

	user, err := s.userRepo.GetFullUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.Verified {
		return nil, errors.New("account is not verified yet")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	return createLoginResponse(user, s.userRepo, s.jwtSecret)
}

func (s *AuthService) SocialLogin(req *SocialLoginRequest) (*LoginResponse, error) {
	var err error
	var name string
	var email string
	var imageUrl *string

	switch req.AuthType {
	case c.AuthTypeGoogle:
		name, email, imageUrl, err = getUserInfoFromGoogle(req.IdToken)
	case c.AuthTypeFacebook:
		// name, email, imageUrl, err = getUserInfoFromFacebook(req.IdToken)
	default:
		err = errors.New("invalid auth type")
	}

	if err != nil {
		return nil, err
	}

	loginUser, err := s.userRepo.GetFullUserByEmail(email)
	if loginUser != nil {
		return createLoginResponse(loginUser, s.userRepo, s.jwtSecret)
	}

	newUser := &user.User{
		Email:    email,
		Password: c.SocialLoginDefaultPassword,
		Name:     name,
		Verified: true,
		Banned:   false,
	}

	if imageUrl != nil && *imageUrl != "" {
		newUser.PathKey = getProfilePicPathKeyFromImageUrl(context.Background(), *imageUrl)
	}

	newUser, err = createRegisterResponse(newUser, s.userRepo, s.roleRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create user via %s login", req.AuthType)
	}

	return createLoginResponse(newUser, s.userRepo, s.jwtSecret)
}

func (s *AuthService) Register(req *RegisterRequest) (*user.UserResponse, error) {
	req.Name = utils.Trim(req.Name)
	req.Email = utils.Trim(strings.ToLower(req.Email))

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

	createdUser, err := createRegisterResponse(user, s.userRepo, s.roleRepo)
	if err != nil {
		return nil, err
	}
	return createdUser.ToResponse(), nil
}

func (s *AuthService) ForgotPassword(req *ForgotPasswordRequest) (string, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	userID, err := s.userRepo.GetUserIdByEmail(req.Email)
	if err != nil || userID == nil {
		return "", nil
	}

	token, err := tokenutil.GenerateNewToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	expiresAt := timeutil.AddHoursUTC(c.PasswordResetTokenExpiryHours)
	if err := s.userRepo.SetPasswordResetToken(*userID, token, expiresAt); err != nil {
		return "", fmt.Errorf("failed to set reset token: %w", err)
	}

	return token, nil
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

func createLoginResponse(user *user.User, userRepo *user.UserRepository, jwtSecret string) (*LoginResponse, error) {
	if user.Banned {
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if lastLoginAt, err := userRepo.SetRefreshTokenAndLastLoginAt(user.ID, refreshToken, refreshExpiresAt); err != nil {
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

func createRegisterResponse(user *user.User, userRepo *user.UserRepository, roleRepo *role.RoleRepository) (*user.User, error) {
	if err := user.HashPassword(); err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	userRole, err := roleRepo.GetRoleWithPermissionsByType(c.RoleTypeUser)
	if err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	user.Roles = []role.Role{*userRole}
	createdUser, err := userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func getUserInfoFromGoogle(idToken string) (name string, email string, imageUrl *string, err error) {
	payload, err := idtoken.Validate(context.Background(), idToken, config.GetConfig().GoogleClientID)
	if err != nil {
		return "", "", nil, errors.New("invalid id token")
	}

	email, ok := payload.Claims[c.UserEmail].(string)
	if !ok {
		return "", "", nil, errors.New("email not found")
	}

	name, ok = payload.Claims[c.UserName].(string)
	if !ok {
		name = "New User"
	}

	name = utils.Trim(name)
	email = utils.Trim(strings.ToLower(email))

	picture, ok := payload.Claims["picture"].(string)
	if ok {
		picture = strings.Replace(picture, "s96-c", "s512-c", 1)
		imageUrl = &picture
	}

	return name, email, imageUrl, nil
}

func getProfilePicPathKeyFromImageUrl(ctx context.Context, imageUrl string) *string {
	if imageUrl == "" {
		l.Logger.Error("Failed to create object storage client", "error", "empty image url", "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	storage, err := file_storage.NewObjectStorage()
	if err != nil {
		l.Logger.Error("Failed to create object storage client", "error", err, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		l.Logger.Error("Failed to get image from url", "error", err, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Logger.Error("Failed to read image from url", "error", err, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	uploadReq := file_storage.StorageUploadRawRequest{
		FileData:    imgData,
		FileName:    "social_user_profile_pic.png",
		Folder:      c.FolderUser,
		ContentType: http.DetectContentType(imgData),
	}

	repo := file_storage.NewFileStorageRepository(storage)
	response, err := repo.UploadRaw(ctx, &uploadReq)
	if err != nil {
		l.Logger.Error("Failed to upload image", "error", err, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	l.Logger.Info("Image uploaded successfully", "path_key", response.PathKey, "method", "getProfilePicPathKeyFromImageUrl")
	return &response.PathKey
}

// TODO: remove this after google login is implemented in frontend
func HandleGoogleLoginTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() == c.EnvDev {
		conf := &oauth2.Config{
			ClientID:     config.GetConfig().GoogleClientID,
			ClientSecret: "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",
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func HandleGoogleCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() == c.EnvDev {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		clientID := "437959137905-1oo0b6bj64tq4ji57hkb52l5epc29729.apps.googleusercontent.com"
		clientSecret := "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"
		redirectURL := "http://localhost:8080/auth/google/callback"

		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}

		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No ID token in response", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "ID Token: %s", idToken)
	}
}
