package constants

const (
	StatusMessageOK               = "Success"
	StatusMessageCreated          = "Resource created successfully"
	StatusMessageUpdated          = "Resource updated successfully"
	StatusMessageDeleted          = "Resource deleted successfully"
	StatusMessageValidationFailed = "Validation failed"
	StatusMessageNotFound         = "Resource not found"
	StatusMessageConflict         = "Resource conflict"
	StatusMessageInternalError    = "Internal server error"
	StatusMessageBadRequest       = "Bad request"
	StatusMessageUnauthorized     = "Unauthorized"
	StatusMessageForbidden        = "Forbidden"
	StatusMessageTooManyRequests  = "Too many requests"
)

const (
	ValidationMessageRequired        = "Field is required"
	ValidationMessageMinLength       = "Field must be at least %d characters long"
	ValidationMessageMaxLength       = "Field must not exceed %d characters"
	ValidationMessageInvalidFormat   = "Field has invalid format"
	ValidationMessageInvalidURL      = "Field must be a valid URL"
	ValidationMessageInvalidEmail    = "Field must be a valid email address"
	ValidationMessageInvalidPhone    = "Field must be a valid phone number"
	ValidationMessageInvalidUUID     = "Field must be a valid UUID"
	ValidationMessagePositiveInteger = "Field must be a positive integer"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
	Include         = "include"
	All             = "all"
	Page            = "page"
	PageSize        = "page_size"
	SortBy          = "sort_by"
	SortOrder       = "sort_order"
	SortOrderDesc   = "DESC"
	SortOrderAsc    = "ASC"
)

const (
	FieldID        = "id"
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
)

const (
	DefaultRateLimit       = 100
	DefaultRateLimitWindow = 60
)

const (
	MaxRequestSize = 10 << 20
	TokenExpiry    = 24 * 60 * 60
)

const (
	DefaultDBTimeout = 30
	MaxDBConnections = 100
)

const (
	MaxFileSize       = 5 << 20
	AllowedImageTypes = "image/jpeg,image/png,image/gif,image/webp"
	AllowedDocTypes   = "application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

const (
	DefaultCacheTTL = 300
	MaxCacheSize    = 1000
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)
