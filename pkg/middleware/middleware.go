package middleware

import (
	"net/http"
	"time"

	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, _ := tokenutil.GenerateNewToken(true)
		start := timeutil.NowUTC()
		l.Logger.Info("🚀🚀🚀 START REQUEST 🚀🚀🚀", "request_id", requestID, "path", r.URL.Path, "method", r.Method)

		next.ServeHTTP(w, r)
		l.Logger.Info("✅✅✅ END REQUEST ✅✅✅", "request_id", requestID, "path", r.URL.Path, "method", r.Method, "duration", time.Since(start).String())
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Logger.Info("CORS Middleware111--------", "path", r.URL.Path, "method", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			l.Logger.Info("CORS Middleware222--------", "path", r.URL.Path, "method", r.Method)
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
				l.Logger.Error("Panic recovered", "error", err, "method", r.Method, "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				response.SendErrorJSON(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
