package system

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type SystemService struct {
	cfg *config.Config
}

func NewSystemService(cfg *config.Config) *SystemService {
	return &SystemService{
		cfg: cfg,
	}
}

func (s *SystemService) SystemHealthCheck(ctx context.Context) *SystemHealthResponse {
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

	return &SystemHealthResponse{
		ServerStatus: "healthy",
		DBStatus:     dbStatus,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *SystemService) HandleGoogleLoginTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		conf := &oauth2.Config{
			ClientID:     s.cfg.GoogleClientID,
			ClientSecret: "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",
			RedirectURL:  fmt.Sprintf(s.cfg.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeGoogle),
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}

		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *SystemService) HandleGoogleLoginCallbackTemp(w http.ResponseWriter, r *http.Request) {
	if config.GetActiveProfile() != c.EnvProd {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		clientID := s.cfg.GoogleClientID
		clientSecret := "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"
		redirectURL := fmt.Sprintf(s.cfg.DomainURL+"/system/social-flow/callback?auth_type=%s", c.AuthTypeGoogle)

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
