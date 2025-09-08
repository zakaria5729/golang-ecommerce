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

// Context Keys
const (
	UserContextKey = "user"
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
	EnvDev   = "dev"
	EnvStage = "stage"
	EnvProd  = "prod"
)

const (
	EnvActiveProfile = "ACTIVE_PROFILE"
	EnvKeyPort       = "PORT"
	EnvKeyDBHost     = "DB_HOST"
	EnvKeyDBPort     = "DB_PORT"
	EnvKeyDBUser     = "DB_USER"
	EnvKeyDBPassword = "DB_PASSWORD"
	EnvKeyDBName     = "DB_NAME"
	EnvKeyDBSSLMode  = "DB_SSL_MODE"
	EnvKeyJWTSecret  = "JWT_SECRET"
)

// Role Types
const (
	RoleTypeSuperAdmin = "SUPER_ADMIN"
	RoleTypeAdmin      = "ADMIN"
	RoleTypeManager    = "MANAGER"
	RoleTypeSeller     = "SELLER"
	RoleTypeUser       = "USER"
)

// Role Names
const (
	RoleNameSuperAdmin = "Super Admin"
	RoleNameAdmin      = "Admin"
	RoleNameManager    = "Manager"
	RoleNameSeller     = "Seller"
	RoleNameUser       = "User"
)

// Permission Names
const (
	// User permissions
	PermissionUserCreate = "user.create"
	PermissionUserRead   = "user.read"
	PermissionUserUpdate = "user.update"
	PermissionUserDelete = "user.delete"

	// Role permissions
	PermissionRoleCreate = "role.create"
	PermissionRoleRead   = "role.read"
	PermissionRoleUpdate = "role.update"
	PermissionRoleDelete = "role.delete"

	// Permission permissions
	PermissionPermissionRead = "permission.read"

	// Category permissions
	PermissionCategoryCreate = "category.create"
	PermissionCategoryRead   = "category.read"
	PermissionCategoryUpdate = "category.update"
	PermissionCategoryDelete = "category.delete"

	// Address permissions
	PermissionAddressCreate = "address.create"
	PermissionAddressRead   = "address.read"
	PermissionAddressUpdate = "address.update"
	PermissionAddressDelete = "address.delete"

	// Review permissions
	PermissionReviewCreate = "review.create"
	PermissionReviewRead   = "review.read"
	PermissionReviewUpdate = "review.update"
	PermissionReviewDelete = "review.delete"

	// Wishlist permissions
	PermissionWishlistCreate = "wishlist.create"
	PermissionWishlistRead   = "wishlist.read"
	PermissionWishlistUpdate = "wishlist.update"
	PermissionWishlistDelete = "wishlist.delete"

	// Browsing History permissions
	PermissionBrowsingHistoryCreate = "browsing_history.create"
	PermissionBrowsingHistoryRead   = "browsing_history.read"
	PermissionBrowsingHistoryUpdate = "browsing_history.update"
	PermissionBrowsingHistoryDelete = "browsing_history.delete"

	// System permissions
	PermissionSystemAdmin      = "system.admin"
	PermissionSystemSuperAdmin = "system.super_admin"
)
