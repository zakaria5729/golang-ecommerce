package constants

type userContextKey string

const (
	UserContextKey     userContextKey = "user"
	UserIDContextKey   userContextKey = "user_id"
	ObjStoreProviderR2 string         = "r2"
	ProjectName        string         = "easy-comerce"
)

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
	EnvKeyJWTSecret  string = "JWT_SECRET"
	EnvActiveProfile string = "ACTIVE_PROFILE"
)

const (
	DefaultRateLimit       = 100
	DefaultRateLimitWindow = 60
)

const (
	MaxRequestSizeMB              = 5 * 1024 * 1024 // 5MB
	AccessTokenExpiryHours        = 12
	RefreshTokenExpiryHours       = 48
	PasswordResetTokenExpiryHours = 5
)

const (
	DefaultDBTimeout = 30
	MaxDBConnections = 100
)

const (
	MaxFileSizeMB     = 5 * 1024 * 1024 // 5MB
	AllowedImageTypes = "image/jpeg,image/jpg,image/png,image/gif,image/svg+xml"
	AllowedDocTypes   = "application/pdf"
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
	EnvDev   string = "dev"
	EnvStage string = "stage"
	EnvProd  string = "prod"
)

const (
	EnvKeyObjStoreRegion          = "OBJ_STORE_REGION"
	EnvKeyObjStoreBucketName      = "OBJ_STORE_BUCKET_NAME"
	EnvKeyObjStoreAccountID       = "OBJ_STORE_ACCOUNT_ID"
	EnvKeyObjStoreAccessKeyID     = "OBJ_STORE_ACCESS_KEY_ID"
	EnvKeyObjStoreAccessKeySecret = "OBJ_STORE_ACCESS_KEY_SECRET"
	EnvKeyObjStorePublicDomain    = "OBJ_STORE_PUBLIC_DOMAIN"
)

const (
	EnvSuperAdminEmail    = "SUPER_ADMIN_EMAIL"
	EnvSuperAdminPassword = "SUPER_ADMIN_PASSWORD"
)

const (
	FolderCategory = "category"
	FolderProduct  = "product"
	FolderUser     = "user"
	FolderReview   = "review"
	FolderAddress  = "address"
	FolderWishlist = "wishlist"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

const (
	MinReviewRating = 1
	MaxReviewRating = 5
)
