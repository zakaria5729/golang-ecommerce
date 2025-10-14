package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/file_storage"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
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

	user.Sanitize()
	createdUser, err := createRegisterResponse(user, s.userRepo, s.roleRepo)
	if err != nil {
		return nil, err
	}

	verificationToken, _ := setVerificationToken(createdUser.ID, s.userRepo)
	userResponse := createdUser.ToResponse()
	if config.GetActiveProfile() != c.EnvProd {
		link := config.GetConfig().DomainURL + "/auth/verify-account?verification_token=" + verificationToken
		userResponse.VerificationLink = &link
	}

	return userResponse, nil
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

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.RefreshTokenExpiryHours)
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

func (s *AuthService) VerifyEmail(token string) error {
	userID, err := s.userRepo.GetUserByVerificationToken(token)
	if err != nil || userID == nil {
		return errors.New("invalid or expired verification token")
	}

	if err := s.userRepo.SetVerifiedAndVerificationToken(*userID, true, nil, nil); err != nil {
		return errors.New("failed to verify user")
	}

	return nil
}

func (s *AuthService) ResendVerifyLink(req *ResendVerifyLinkRequest) (verificationLink string, verified bool, err error) {
	userID, verified, err := s.userRepo.GetUserIdAndVerifiedByEmail(req.Email, nil)
	if err != nil || userID == nil {
		return "", false, err
	}

	if !verified {
		verificationToken, err := setVerificationToken(*userID, s.userRepo)
		if err != nil {
			return "", false, err
		}

		if config.GetActiveProfile() != c.EnvProd {
			return config.GetConfig().DomainURL + "/auth/verify-account?verification_token=" + verificationToken, false, nil
		}
	}

	return "", verified, nil
}

func createLoginResponse(user *user.User, userRepo *user.UserRepository, jwtSecret string) (*LoginResponse, error) {
	if user.Banned {
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.RefreshTokenExpiryHours)
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

func setVerificationToken(userID uint, userRepo *user.UserRepository) (string, error) {
	verificationToken, expiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.VerificationTokenExpiryHours)
	if err != nil {
		l.Logger.Error("❌ Failed to generate verification token", "method", "Register", "error", err, "userID", userID)
		return "", errors.New("Failed to generate verification token")
	}

	if verificationToken != "" {
		if err := userRepo.SetVerifiedAndVerificationToken(userID, false, &verificationToken, &expiresAt); err != nil {
			l.Logger.Error("❌ Failed to set verification token", "method", "Register", "error", err, "userID", userID)
			return "", errors.New("Failed to set verification token")
		}
	}

	return verificationToken, nil
}

func getUserInfoFromGoogle(idToken string) (name string, email string, imageUrl *string, err error) {
	googleClientId := config.GetConfig().GoogleClientID
	if googleClientId == "" {
		return "", "", nil, errors.New("Social login is not enabled")
	}

	payload, err := idtoken.Validate(context.Background(), idToken, googleClientId)
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
