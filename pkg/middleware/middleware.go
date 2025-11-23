package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/timeutil"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isExcludedPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := timeutil.NowUTC()
		l.Info("🟢 START-REQUEST 🟢", "path", r.URL.Path, "http-method", r.Method)

		next.ServeHTTP(w, r)
		l.Info("✅ END-REQUEST ✅", "path", r.URL.Path, "http-method", r.Method, "duration", time.Since(start).String())
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errStack := string(debug.Stack())

				if strings.EqualFold(config.GetActiveProfile(), c.EnvLocal) {
					fmt.Printf(`"❌ Panic recovered", "error" %v, "method" %v, "path" %v`, errStack, r.Method, r.URL.Path)
				} else {
					l.Error("❌ Panic recovered", "error", errStack, "method", r.Method, "path", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func TrailingSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

func isExcludedPath(path string) bool {
	excluded := []string{
		"/favicon.ico",
		"/app-health",
		"/metrics",
		"/system/logs",
		"/auth/social-flow",
		"/auth/social-flow/callback",
	}

	for _, prefix := range excluded {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
