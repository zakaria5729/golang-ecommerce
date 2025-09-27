package middleware

import (
	"net/http"
	"time"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"github.com/easy-comerce/backend/pkg/tokenutil"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := tokenutil.GenerateNewToken()
		if err != nil {
			requestID = timeutil.NowUTC().Format("20060102150405")
		}
		start := timeutil.NowUTC()
		logger.Logger.Info("🚀🚀🚀 START REQUEST 🚀🚀🚀", "request_id", requestID, "path", r.URL.Path)

		next.ServeHTTP(w, r)
		logger.Logger.Info("✅✅✅ END REQUEST ✅✅✅", "path", "request_id", requestID, r.URL.Path, "duration", time.Since(start).String())
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == constants.OPTIONS {
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
				logger.Logger.Error("Panic recovered", "error", err, "method", r.Method, "path", r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"success":false,"error":{"message":"Internal server error","code":"INTERNAL_ERROR"}}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	clients := make(map[string]*ClientInfo)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		client, exists := clients[clientIP]
		if !exists {
			client = &ClientInfo{
				IP:           clientIP,
				RequestCount: 0,
				LastRequest:  timeutil.NowUTC(),
			}
			clients[clientIP] = client
		}

		now := timeutil.NowUTC()
		if now.Sub(client.LastRequest) > time.Minute {
			client.RequestCount = 0
			client.LastRequest = now
		}

		if client.RequestCount >= 100 {
			logger.Logger.Warn("Rate limit exceeded", "clientIP", clientIP, "path", r.URL.Path)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"success":false,"error":{"message":"Rate limit exceeded","code":"RATE_LIMIT_EXCEEDED"}}`))
			return
		}

		client.RequestCount++
		next.ServeHTTP(w, r)
	})
}

type ClientInfo struct {
	IP           string
	RequestCount int
	LastRequest  time.Time
}
