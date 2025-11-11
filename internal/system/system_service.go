package system

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/easy-comerce/backend/db"
	m "github.com/easy-comerce/backend/internal/system/model"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type SystemService interface {
	SystemHealthCheck(ctx context.Context) *m.SystemHealthResponse
	GetSystemLogFiles(fileName string) ([]m.SystemLogFileResponse, error)
	DownloadSystemLogFile(fileName string) ([]byte, error)
	DeleteSystemLogFile(fileName string) error
	HandleGoogleLoginTemp(w http.ResponseWriter, r *http.Request)
	HandleGoogleLoginCallbackTemp(w http.ResponseWriter, r *http.Request)
	HandleFacebookLoginTemp(w http.ResponseWriter, r *http.Request)
	HandleFacebookLoginCallbackTemp(w http.ResponseWriter, r *http.Request)
}

type systemService struct {
	cfg *config.Config
}

func NewSystemService(cfg *config.Config) SystemService {
	return &systemService{
		cfg: cfg,
	}
}

func (s *systemService) SystemHealthCheck(ctx context.Context) *m.SystemHealthResponse {
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

	return &m.SystemHealthResponse{
		ServerStatus: "healthy",
		DBStatus:     dbStatus,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *systemService) GetSystemLogFiles(fileName string) ([]m.SystemLogFileResponse, error) {
	logsDir := filepath.Join(utils.GetProjectRootPath(), "logs")
	files, err := os.ReadDir(logsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs directory: %v", err)
	}

	var logFiles []m.SystemLogFileResponse
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			if fileName != "" && !strings.Contains(file.Name(), fileName) {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}

			var size string
			if info.Size() < c.SizeInMB {
				size = fmt.Sprintf("%.2f KB", float64(info.Size())/1024.0)
			} else {
				size = fmt.Sprintf("%.2f MB", float64(info.Size())/float64(c.SizeInMB))
			}

			logFiles = append(logFiles, m.SystemLogFileResponse{
				FileName:   file.Name(),
				Size:       size,
				ModifiedAt: info.ModTime().Format(time.RFC3339),
			})
		}
	}

	return logFiles, nil
}

func (s *systemService) DownloadSystemLogFile(fileName string) ([]byte, error) {
	if !strings.HasSuffix(fileName, ".log") {
		return nil, fmt.Errorf("invalid log file name")
	}

	logPath := filepath.Join(utils.GetProjectRootPath(), "logs", fileName)
	content, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load log file: %v", err)
	}

	return content, nil
}

func (s *systemService) DeleteSystemLogFile(fileName string) error {
	if !strings.HasSuffix(fileName, ".log") {
		return fmt.Errorf("invalid log file name")
	}

	logFileTime, err := time.Parse(c.LogFileFormat, strings.TrimSuffix(strings.TrimPrefix(fileName, "app-"), ".log"))
	if err != nil {
		return fmt.Errorf("failed to parse log file name: %v", err)
	}

	if logFileTime.After(timeutil.AddDaysUTC(-1)) {
		return errors.New("you can not delete today and yesterday's log files")
	}

	logPath := filepath.Join(utils.GetProjectRootPath(), "logs", fileName)
	err = os.Remove(logPath)
	if err != nil {
		return fmt.Errorf("failed to delete log file: %v", err)
	}

	return nil
}

func (s *systemService) HandleGoogleLoginTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		conf := &oauth2.Config{
			ClientID:     s.cfg.ExtServiceConfig.GoogleClientID,
			ClientSecret: "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",
			RedirectURL:  fmt.Sprintf(s.cfg.AppConfig.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeGoogle),
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *systemService) HandleGoogleLoginCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		code := r.URL.Query().Get("code")
		if code == "" {
			response.SendErrorJSON(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		clientID := s.cfg.ExtServiceConfig.GoogleClientID
		clientSecret := "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"
		redirectURL := fmt.Sprintf(s.cfg.AppConfig.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeGoogle)

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

		fmt.Fprintf(w, token.AccessToken)
	}
}

func (s *systemService) HandleFacebookLoginTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		conf := &oauth2.Config{
			ClientID:     s.cfg.ExtServiceConfig.FacebookAppID,
			ClientSecret: "ce26f2b347e99cc70d4c9f3b4b8b0bbd",
			RedirectURL:  fmt.Sprintf(s.cfg.AppConfig.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeFacebook),
			Scopes:       []string{"email"},
			Endpoint:     facebook.Endpoint,
		}

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *systemService) HandleFacebookLoginCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		conf := &oauth2.Config{
			ClientID:     s.cfg.ExtServiceConfig.FacebookAppID,
			ClientSecret: "ce26f2b347e99cc70d4c9f3b4b8b0bbd",
			RedirectURL:  fmt.Sprintf(s.cfg.AppConfig.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeFacebook),
			Scopes:       []string{"email"},
			Endpoint:     facebook.Endpoint,
		}

		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		client := conf.Client(context.Background(), token)
		resp, err := client.Get(c.FacebookUserInfoURL)
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, string(body))
	}
}
