package constants

type stringKey string

const (
	UserContextKey     stringKey = "user"
	UserIDContextKey   stringKey = "user_id"
	ObjStoreProviderR2 string    = "r2"
	ProjectName        string    = "easy-comerce"
)

const (
	StatusMessageOK               string = "Success"
	StatusMessageCreated          string = "Resource created successfully"
	StatusMessageUpdated          string = "Resource updated successfully"
	StatusMessageDeleted          string = "Resource deleted successfully"
	StatusMessageValidationFailed string = "Validation failed"
	StatusMessageNotFound         string = "Resource not found"
	StatusMessageConflict         string = "Resource conflict"
	StatusMessageInternalError    string = "Internal server error"
	StatusMessageBadRequest       string = "Bad request"
	StatusMessageUnauthorized     string = "Unauthorized"
	StatusMessageForbidden        string = "Forbidden"
	StatusMessageTooManyRequests  string = "Too many requests"
)

const (
	ValidationMessageRequired        string = "Field is required"
	ValidationMessageMinLength       string = "Field must be at least %d characters long"
	ValidationMessageMaxLength       string = "Field must not exceed %d characters"
	ValidationMessageInvalidFormat   string = "Field has invalid format"
	ValidationMessageInvalidURL      string = "Field must be a valid URL"
	ValidationMessageInvalidEmail    string = "Field must be a valid email address"
	ValidationMessageInvalidPhone    string = "Field must be a valid phone number"
	ValidationMessageInvalidUUID     string = "Field must be a valid UUID"
	ValidationMessagePositiveInteger string = "Field must be a positive integer"
)

const (
	EnvKeyJWTSecret  string = "JWT_SECRET"
	EnvActiveProfile string = "ACTIVE_PROFILE"
)

const (
	AccessTokenExpiryHours        int = 24
	RefreshTokenExpiryHours       int = 48
	PasswordResetTokenExpiryHours int = 6
)

const (
	DefaultDBTimeout       int = 30
	MaxDBConnections       int = 100
	DefaultRateLimit       int = 100
	DefaultRateLimitWindow int = 60
)

const (
	SizeInMB          int64  = 1024 * 1024 // 1MB
	AllowedImageTypes string = "image/jpeg,image/jpg,image/png,image/gif,image/svg+xml"
	AllowedDocTypes   string = "application/pdf"
)

const (
	EnvDev   string = "dev"
	EnvStage string = "stage"
	EnvProd  string = "prod"
)

const (
	EnvKeyObjStoreRegion          string = "OBJ_STORE_REGION"
	EnvKeyObjStoreBucketName      string = "OBJ_STORE_BUCKET_NAME"
	EnvKeyObjStoreAccountID       string = "OBJ_STORE_ACCOUNT_ID"
	EnvKeyObjStoreAccessKeyID     string = "OBJ_STORE_ACCESS_KEY_ID"
	EnvKeyObjStoreAccessKeySecret string = "OBJ_STORE_ACCESS_KEY_SECRET"
	EnvKeyObjStorePublicDomain    string = "OBJ_STORE_PUBLIC_DOMAIN"
)

const (
	FolderCategory string = "category"
	FolderProduct  string = "product"
	FolderUser     string = "user"
)

const (
	VersionPrefix string = "/v"
	V1            string = "v1"
)

const (
	MinReviewRating        int = 1
	MaxReviewRating        int = 5
	MaxReviewCommentLength int = 1000
	MaxPriorityLimit       int = 100
	SubcategoryDepthLimit  int = 10
)
