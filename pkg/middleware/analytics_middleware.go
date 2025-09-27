package middleware

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/easy-comerce/backend/internal/feature/analytics"
	"github.com/easy-comerce/backend/pkg/logger"
)

// AnalyticsMiddleware tracks requests automatically
func AnalyticsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Skip analytics for certain paths
			if shouldSkipAnalytics(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Create a custom response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Process the request
			next.ServeHTTP(rw, r)

			// Track the request asynchronously
			go func() {
				ctx := context.Background()
				uc := analytics.NewAnalyticsUseCase()

				// Generate visitor ID from IP + User Agent
				visitorID := generateVisitorID(r)
				sessionID := generateSessionID(r)

				// Create page view
				ipAddress := getClientIP(r)
				userAgent := r.UserAgent()
				pageView := &analytics.AnalyticsPageView{
					VisitorID:    visitorID,
					SessionID:    sessionID,
					Path:         r.URL.Path,
					PageTitle:    getPageTitle(r.URL.Path),
					HTTPMethod:   r.Method,
					StatusCode:   rw.statusCode,
					ResponseTime: int(time.Since(start).Milliseconds()),
					IPAddress:    &ipAddress,
					UserAgent:    &userAgent,
					Referrer:     getReferrer(r),
					Country:      getCountry(ipAddress),
					City:         getCity(ipAddress),
					Region:       getRegion(ipAddress),
					Timezone:     getTimezone(ipAddress),
					IsBot:        isBot(r.UserAgent()),
					BotName:      getBotName(r.UserAgent()),
					VisitedAt:    time.Now(),
				}

				// Track the page view
				err := uc.TrackPageView(ctx, pageView)
				if err != nil {
					logger.Logger.Error("Failed to track page view", "error", err, "path", r.URL.Path)
				}
			}()
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// shouldSkipAnalytics determines if analytics should be skipped for this path
func shouldSkipAnalytics(path string) bool {
	skipPaths := []string{
		"/health",
		"/analytics",
		"/favicon.ico",
		"/robots.txt",
		"/sitemap.xml",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

// generateVisitorID creates a unique visitor ID from IP and User Agent
func generateVisitorID(r *http.Request) string {
	ip := getClientIP(r)
	userAgent := r.UserAgent()
	data := fmt.Sprintf("%s-%s", ip, userAgent)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// generateSessionID creates a session ID (simplified - in production, use proper session management)
func generateSessionID(r *http.Request) string {
	ip := getClientIP(r)
	userAgent := r.UserAgent()
	now := time.Now().Format("2006-01-02")
	data := fmt.Sprintf("%s-%s-%s", ip, userAgent, now)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// getClientIP extracts the real client IP address
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			// Handle IPv6 addresses by removing brackets
			if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
				ip = ip[1 : len(ip)-1]
			}
			return ip
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		// Handle IPv6 addresses by removing brackets
		if strings.HasPrefix(xri, "[") && strings.HasSuffix(xri, "]") {
			return xri[1 : len(xri)-1]
		}
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Handle IPv6 addresses by removing brackets
	if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
		ip = ip[1 : len(ip)-1]
	} else if idx := strings.LastIndex(ip, ":"); idx != -1 {
		// For IPv4, remove port
		ip = ip[:idx]
	}
	return ip
}

// getReferrer extracts the referrer from the request
func getReferrer(r *http.Request) *string {
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		return nil
	}
	return &referrer
}

// getQueryParams extracts query parameters as JSON string
func getQueryParams(r *http.Request) *string {
	if len(r.URL.Query()) == 0 {
		return nil
	}

	params := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil
	}
	jsonStr := string(jsonBytes)
	return &jsonStr
}

// getLanguage extracts the language from Accept-Language header
func getLanguage(r *http.Request) *string {
	lang := r.Header.Get("Accept-Language")
	if lang == "" {
		return nil
	}

	// Take the first language preference
	if idx := strings.Index(lang, ","); idx != -1 {
		lang = lang[:idx]
	}
	return &lang
}

// getPageTitle generates a page title based on the path
func getPageTitle(path string) *string {
	if path == "/" {
		title := "Home"
		return &title
	}

	// Remove leading slash and capitalize
	title := strings.TrimPrefix(path, "/")
	if title == "" {
		title = "Home"
	}

	// Replace hyphens and underscores with spaces
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	title = strings.Join(words, " ")

	return &title
}

// isBot checks if the request is from a bot
func isBot(userAgent string) bool {
	botPatterns := []string{
		"googlebot", "bingbot", "slurp", "duckduckbot", "baiduspider",
		"yandexbot", "facebookexternalhit", "twitterbot", "linkedinbot",
		"whatsapp", "telegram", "skype", "discord", "slack",
		"curl", "wget", "postman", "insomnia", "httpie",
		"pingdom", "uptimerobot", "statuscake", "newrelic",
		"ahrefs", "semrush", "moz", "screaming frog",
	}

	ua := strings.ToLower(userAgent)
	for _, pattern := range botPatterns {
		if strings.Contains(ua, pattern) {
			return true
		}
	}

	return false
}

// getCountry returns the country for the given IP address
func getCountry(ip string) *string {
	// This is a simplified implementation. For production, consider using
	// a proper IP geolocation service like MaxMind GeoIP2 or similar.
	// For now, return nil to indicate no country data available.
	return nil
}

// getCity returns the city for the given IP address
func getCity(ip string) *string {
	// This is a simplified implementation. For production, consider using
	// a proper IP geolocation service like MaxMind GeoIP2 or similar.
	// For now, return nil to indicate no city data available.
	return nil
}

// getRegion returns the region for the given IP address
func getRegion(ip string) *string {
	// This is a simplified implementation. For production, consider using
	// a proper IP geolocation service like MaxMind GeoIP2 or similar.
	// For now, return nil to indicate no region data available.
	return nil
}

// getTimezone returns the timezone for the given IP address
func getTimezone(ip string) *string {
	// This is a simplified implementation. For production, consider using
	// a proper IP geolocation service like MaxMind GeoIP2 or similar.
	// For now, return nil to indicate no timezone data available.
	return nil
}

// getBotName returns the bot name if the user agent is from a known bot
func getBotName(userAgent string) *string {
	botPatterns := map[string]string{
		"googlebot":           "Googlebot",
		"bingbot":             "Bingbot",
		"slurp":               "Yahoo! Slurp",
		"duckduckbot":         "DuckDuckBot",
		"baiduspider":         "Baiduspider",
		"yandexbot":           "YandexBot",
		"facebookexternalhit": "Facebook External Hit",
		"twitterbot":          "Twitterbot",
		"linkedinbot":         "LinkedInBot",
		"whatsapp":            "WhatsApp",
		"telegram":            "Telegram",
		"skype":               "Skype",
		"discord":             "Discord",
		"slack":               "Slack",
		"curl":                "cURL",
		"wget":                "Wget",
		"postman":             "Postman",
		"insomnia":            "Insomnia",
		"httpie":              "HTTPie",
		"pingdom":             "Pingdom",
		"uptimerobot":         "UptimeRobot",
		"statuscake":          "StatusCake",
		"newrelic":            "New Relic",
		"ahrefs":              "AhrefsBot",
		"semrush":             "SemrushBot",
		"moz":                 "Moz",
		"screaming frog":      "Screaming Frog",
	}

	ua := strings.ToLower(userAgent)
	for pattern, name := range botPatterns {
		if strings.Contains(ua, pattern) {
			return &name
		}
	}

	return nil
}
