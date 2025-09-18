package constants

const (
	EnvKeyPort       = "PORT"
	EnvKeyDBHost     = "DB_HOST"
	EnvKeyDBPort     = "DB_PORT"
	EnvKeyDBUser     = "DB_USER"
	EnvKeyDBPassword = "DB_PASSWORD"
	EnvKeyDBName     = "DB_NAME"
	EnvKeyDBSSLMode  = "DB_SSL_MODE"
)

const (
	TableUser            = "users"
	TableUserRole        = "user_roles"
	TableRole            = "roles"
	TableRolePermission  = "role_permissions"
	TablePermission      = "permissions"
	TableUserStatus      = "user_status"
	TableCategory        = "categories"
	TableAddress         = "addresses"
	TableReview          = "reviews"
	TableWishlist        = "wishlists"
	TableBrowsingHistory = "browsing_histories"
)

const (
	FieldID        = "id"
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldDeletedAt = "deleted_at"
)

const (
	UserEmail                       = "email"
	UserName                        = "name"
	UserVerified                    = "verified"
	UserBanned                      = "banned"
	UserLastLoginAt                 = "last_login_at"
	UserRoles                       = "roles"
	UserPermissions                 = "permissions"
	UserPassword                    = "password"
	UserPathKey                     = "path_key"
	UserImageURL                    = "image_url"
	UserRefreshToken                = "refresh_token"
	UserRefreshTokenExpires         = "refresh_token_expires"
	UserPasswordResetToken          = "password_reset_token"
	UserPasswordResetExpires        = "password_reset_expires"
	UserRolesCapitalized            = "Roles"
	UserRolesPermissionsCapitalized = "Roles.Permissions"
	UsersCapitalized                = "Users"
)

const (
	AddressUserID       = "user_id"
	AddressStreet       = "street"
	AddressCity         = "city"
	AddressState        = "state"
	AddressZipCode      = "zip_code"
	AddressCountry      = "country"
	AddressIsDefault    = "is_default"
	AddressAddressType  = "address_type"
	AddressTypeShipping = "shipping"
	AddressTypeBilling  = "billing"
)

const (
	BrowsingHistoryUserID    = "user_id"
	BrowsingHistoryProductID = "product_id"
	BrowsingHistoryViewedAt  = "viewed_at"
)

const (
	CategoryTitle    = "title"
	CategorySubTitle = "sub_title"
	CategoryImageURL = "image_url"
	CategoryParentID = "parent_id"
	CategoryIsActive = "is_active"
	CategoryPriority = "priority"
)

const (
	PermissionName        = "name"
	PermissionDescription = "description"
)

const (
	ReviewProductID = "product_id"
	ReviewUserID    = "user_id"
	ReviewRating    = "rating"
	ReviewComment   = "comment"
)

const (
	RoleRoleName               = "role_name"
	RoleRoleType               = "role_type"
	RoleDescription            = "description"
	RolePermissions            = "role_permissions"
	RolePermissionsCapitalized = "Permissions"
)

const (
	WishlistUserID    = "user_id"
	WishlistProductID = "product_id"
)
