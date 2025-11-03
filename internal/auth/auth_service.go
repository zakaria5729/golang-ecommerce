package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/easy-comerce/backend/internal/auth/model"
	"github.com/easy-comerce/backend/internal/file_storage"
	fm "github.com/easy-comerce/backend/internal/file_storage/model"
	"github.com/easy-comerce/backend/internal/role"
	"github.com/easy-comerce/backend/internal/user"
	userModel "github.com/easy-comerce/backend/internal/user/model"
	apperror "github.com/easy-comerce/backend/pkg/app_error"
	e "github.com/easy-comerce/backend/pkg/app_error"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	httpclient "github.com/easy-comerce/backend/pkg/http_client"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	SocialLogin(req *model.SocialLoginRequest) (*model.LoginResponse, error)
	Register(req *model.RegisterRequest) (*userModel.UserResponse, error)
	ForgotPassword(req *model.ForgotPasswordRequest) (string, error)
	ResetPassword(req *model.ResetPasswordRequest) error
	RefreshToken(req *model.RefreshTokenRequest) (*model.LoginResponse, error)
	Logout(userID uint) error
	VerifyEmail(token string) error
	ResendVerifyLink(req *model.ResendVerifyLinkRequest) (verificationLink string, verified bool, err error)
}

type authService struct {
	jwtSecret string
	userRepo  user.UserRepository
	roleRepo  role.RoleRepository
}

func NewAuthService(
	jwtSecret string,
	userRepo user.UserRepository,
	roleRepo role.RoleRepository,
) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
	}
}

func (s *authService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))
	req.Password = utils.Trim(req.Password)

	userEntity, err := s.userRepo.GetFullUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !userEntity.Verified {
		return nil, errors.New("account is not verified yet")
	}

	if !userEntity.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	return createLoginResponse(userEntity, s.userRepo, s.jwtSecret)
}

func (s *authService) SocialLogin(req *model.SocialLoginRequest) (*model.LoginResponse, error) {
	var err error
	var name string
	var email string
	var imageUrl *string

	switch req.AuthType {
	case c.AuthTypeGoogle:
		name, email, imageUrl, err = getUserInfoFromGoogle(req.AccessToken)
	case c.AuthTypeFacebook:
		name, email, imageUrl, err = getUserInfoFromFacebook(req.AccessToken)
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

	newUser := &user.UserEntity{
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
		return nil, fmt.Errorf("failed to create user via %s login: %w", req.AuthType, err)
	}

	return createLoginResponse(newUser, s.userRepo, s.jwtSecret)
}

func (s *authService) Register(req *model.RegisterRequest) (*userModel.UserResponse, error) {
	req.Name = utils.Trim(req.Name)
	req.Email = utils.Trim(strings.ToLower(req.Email))

	exists, err := s.userRepo.IsUserExists(req.Email)
	if exists {
		return nil, errors.New("user with this email already exists")
	}

	user := &user.UserEntity{
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

func (s *authService) ForgotPassword(req *model.ForgotPasswordRequest) (string, error) {
	req.Email = utils.Trim(strings.ToLower(req.Email))

	userID, err := s.userRepo.GetUserIdByEmail(req.Email)
	if err != nil || userID == nil {
		return "", err
	}

	token, err := tokenutil.GenerateNewToken()
	if err != nil {
		return "", err
	}

	expiresAt := timeutil.AddHoursUTC(c.PasswordResetTokenExpiryHours)
	if err := s.userRepo.SetPasswordResetToken(*userID, token, expiresAt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ResetPassword(req *model.ResetPasswordRequest) error {
	user, err := s.userRepo.GetUserByResetPasswordToken(req.Token)
	if err != nil || user == nil || user.PasswordResetExpires.Before(timeutil.NowUTC()) {
		return e.NewServerError("invalid or expired reset token")
	}

	if user.Email == config.GetConfig().SuperAdminEmail {
		return errors.New("can not reset super admin password")
	}

	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return e.WrapServerError("failed to process new password", err)
	}

	if err := s.userRepo.ResetPassword(user.ID, user.Password); err != nil {
		return e.NewServerError("failed to reset user password")
	}

	return nil
}

func (s *authService) RefreshToken(req *model.RefreshTokenRequest) (*model.LoginResponse, error) {
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
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.RefreshTokenExpiryHours)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.SetRefreshToken(user.ID, &refreshToken, &refreshExpiresAt); err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *authService) Logout(userID uint) error {
	if err := s.userRepo.SetRefreshToken(userID, nil, nil); err != nil {
		return e.WrapServerError("failed to clear refresh token", err)
	}
	return nil
}

func (s *authService) VerifyEmail(token string) error {
	userID, err := s.userRepo.GetUserByVerificationToken(token)
	if err != nil || userID == nil {
		return err
	}

	if err := s.userRepo.SetVerifiedAndVerificationToken(*userID, true, nil, nil); err != nil {
		return err
	}

	return nil
}

func (s *authService) ResendVerifyLink(req *model.ResendVerifyLinkRequest) (verificationLink string, verified bool, err error) {
	userID, verified, err := s.userRepo.GetUserIdAndVerifiedByEmail(req.Email, nil)
	if err != nil || userID == nil {
		return "", false, e.NewServerError("failed to get user id and verified")
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

func createLoginResponse(user *user.UserEntity, userRepo user.UserRepository, jwtSecret string) (*model.LoginResponse, error) {
	if user.Banned {
		return nil, errors.New("account is banned")
	}

	accessToken, expiresAt, err := tokenutil.GenerateNewJwtToken(user, jwtSecret)
	if err != nil {
		return nil, apperror.WrapServerError("failed to generate token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.RefreshTokenExpiryHours)
	if err != nil {
		return nil, apperror.WrapServerError("failed to generate refresh token: %w", err)
	}

	if lastLoginAt, err := userRepo.SetRefreshTokenAndLastLoginAt(user.ID, refreshToken, refreshExpiresAt); err != nil {
		return nil, err
	} else {
		user.LastLoginAt = &lastLoginAt
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresAt:    expiresAt,
	}, nil
}

func createRegisterResponse(user *user.UserEntity, userRepo user.UserRepository, roleRepo role.RoleRepository) (*user.UserEntity, error) {
	if err := user.HashPassword(); err != nil {
		return nil, apperror.WrapServerError("failed to process password", err)
	}

	userRole, err := roleRepo.GetRoleWithPermissionsByType(c.RoleTypeUser)
	if err != nil {
		return nil, err
	}

	user.Roles = []role.RoleEntity{*userRole}
	createdUser, err := userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func setVerificationToken(userID uint, userRepo user.UserRepository) (string, error) {
	verificationToken, expiresAt, err := tokenutil.GenerateNewTokenWithExpiryTime(c.VerificationTokenExpiryHours)
	if err != nil {
		l.Logger.Error("❌ Failed to generate verification token", "method", "Register", "error", err, "userID", userID)
		return "", e.WrapServerError("Failed to generate verification token", err)
	}

	if verificationToken != "" {
		if err := userRepo.SetVerifiedAndVerificationToken(userID, false, &verificationToken, &expiresAt); err != nil {
			l.Logger.Error("❌ Failed to set verification token", "method", "Register", "error", err, "userID", userID)

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return "", e.WrapServerError("Failed to set verification token", err)
			}

			return "", errors.New("Failed to set verification token")
		}
	}

	return verificationToken, nil
}

func getUserInfoFromGoogle(idToken string) (name string, email string, imageUrl *string, err error) {
	googleClientId := config.GetConfig().GoogleClientID
	if googleClientId == "" {
		return "", "", nil, errors.New("Google login is not properly configured")
	}

	var res map[string]any
	if err = httpclient.New().Retry(1).Do(httpclient.Request{
		URL:      c.GoogleUserInfoURL,
		Response: &res,
		Headers: map[string]string{
			c.Authorization: c.Bearer + " " + idToken,
		},
	}); err != nil {
		return "", "", nil, fmt.Errorf("❌ failed to get user info from Google: %v", err)
	}

	email, ok := res[c.UserEmail].(string)
	if !ok {
		return "", "", nil, errors.New("email not found")
	}
	name, ok = res[c.UserName].(string)
	if !ok {
		name = "New User"
	}

	picture, ok := res["picture"].(string)
	if ok {
		picture = strings.Replace(picture, "s96-c", "s512-c", 1)
		imageUrl = &picture
	}

	return utils.Trim(name), utils.Trim(strings.ToLower(email)), imageUrl, nil
}

func getUserInfoFromFacebook(accessToken string) (name string, email string, imageUrl *string, err error) {
	facebookAppID := config.GetConfig().FacebookAppID
	if facebookAppID == "" {
		return "", "", nil, errors.New("Facebook login is not properly configured")
	}

	var res map[string]any
	if err = httpclient.New().Retry(1).Do(httpclient.Request{
		URL:      c.FacebookUserInfoURL,
		Response: &res,
		QueryParams: map[string]string{
			"access_token": accessToken,
		},
	}); err != nil {
		return "", "", nil, fmt.Errorf("❌ failed to get user info from Facebook: %v", err)
	}

	l.Logger.Info("✅ Got user info from Facebook", "method", "getUserInfoFromFacebook", "res", res)
	return name, email, imageUrl, nil
}

func getProfilePicPathKeyFromImageUrl(ctx context.Context, imageUrl string) *string {
	if imageUrl == "" {
		l.Logger.Error("❌ Failed to create object storage client", "error", "empty image url", "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	storage, err := file_storage.NewObjectStorage()
	if err != nil {
		l.Logger.Error("❌ Failed to create object storage client", "error", err, "image_url", imageUrl, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	imageUrl = utils.Trim(imageUrl)
	var imgData []byte
	err = httpclient.New().Retry(2).Do(httpclient.Request{
		URL:      imageUrl,
		Response: &imgData,
	})
	if err != nil {
		l.Logger.Error("❌ Failed to get image from url", "error", err, "image_url", imageUrl, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	uploadReq := fm.StorageUploadRawRequest{
		FileData:    imgData,
		FileName:    "social_user_profile_pic.png",
		Folder:      c.FolderUser,
		ContentType: http.DetectContentType(imgData),
	}

	repo := file_storage.NewFileStorageRepository(storage)
	response, err := repo.UploadRaw(ctx, &uploadReq)
	if err != nil {
		l.Logger.Error("❌ Failed to upload image", "error", err, "image_url", imageUrl, "method", "getProfilePicPathKeyFromImageUrl")
		return nil
	}

	return &response.PathKey
}
