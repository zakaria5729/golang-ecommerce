package constants

const (
	EnvKeyPort            string = "PORT"
	EnvKeyDBHost          string = "DB_HOST"
	EnvKeyDBPort          string = "DB_PORT"
	EnvKeyDBUser          string = "DB_USER"
	EnvKeyDBPassword      string = "DB_PASSWORD"
	EnvKeyDBName          string = "DB_NAME"
	EnvKeyDBSSLMode       string = "DB_SSL_MODE"
	EnvKeyDBShowLog       string = "DB_SHOW_LOG"
	EnvSuperAdminEmail    string = "SUPER_ADMIN_EMAIL"
	EnvSuperAdminPassword string = "SUPER_ADMIN_PASSWORD"
)

const (
	TableUser            string = "users"
	TableUserRole        string = "user_roles"
	TableRole            string = "roles"
	TableRolePermission  string = "role_permissions"
	TablePermission      string = "permissions"
	TableUserStatus      string = "user_status"
	TableCategory        string = "categories"
	TableAddress         string = "addresses"
	TableReview          string = "reviews"
	TableWishlist        string = "wishlists"
	TableProductStats    string = "product_stats"
	TableBrand           string = "brands"
	TableColor           string = "colors"
	TableAttributeType   string = "attribute_types"
	TableAttributeOption string = "attribute_options"
	TableSizeCategory    string = "size_categories"
	TableSizeOption      string = "size_options"
	TableAnalytics       string = "analytics"
)

const (
	FieldID        string = "id"
	FieldCreatedAt string = "created_at"
	FieldUpdatedAt string = "updated_at"
	FieldUpdatedBy string = "updated_by"
	FieldDeletedBy string = "deleted_by"
	FieldDeletedAt string = "deleted_at"
	FieldUserID    string = "user_id"
	FieldRoleID    string = "role_id"
)

const (
	UserEmail                       string = "email"
	UserName                        string = "name"
	UserVerified                    string = "verified"
	UserBanned                      string = "banned"
	UserLastLoginAt                 string = "last_login_at"
	UserRoles                       string = "roles"
	UserPermissions                 string = "permissions"
	UserPassword                    string = "password"
	UserPathKey                     string = "path_key"
	UserImageURL                    string = "image_url"
	UserRefreshToken                string = "refresh_token"
	UserRefreshTokenExpires         string = "refresh_token_expires"
	UserPasswordResetToken          string = "password_reset_token"
	UserPasswordResetExpires        string = "password_reset_expires"
	UserRolesCapitalized            string = "Roles"
	UserRolesPermissionsCapitalized string = "Roles.Permissions"
	UsersCapitalized                string = "Users"
	UserPurchaseCount               string = "purchase_count"
	UserTotalSpent                  string = "total_spent"
)

const (
	AddressUserID       string = "user_id"
	AddressStreet       string = "street"
	AddressCity         string = "city"
	AddressState        string = "state"
	AddressZipCode      string = "zip_code"
	AddressCountry      string = "country"
	AddressIsDefault    string = "is_default"
	AddressAddressType  string = "address_type"
	AddressTypeShipping string = "shipping"
	AddressTypeBilling  string = "billing"
)

const (
	ProductStatsProductID               string = "product_id"
	ProductStatsViewCount               string = "view_count"
	ProductStatsAddToCartCount          string = "add_to_cart_count"
	ProductStatsRemoveFromCartCount     string = "remove_from_cart_count"
	ProductStatsPurchaseCount           string = "purchase_count"
	ProductStatsWishlistCount           string = "wishlist_count"
	ProductStatsRemoveFromWishlistCount string = "remove_from_wishlist_count"
)

const (
	CategoryTitle    string = "title"
	CategorySubTitle string = "sub_title"
	CategoryImageURL string = "image_url"
	CategoryParentID string = "parent_id"
	CategoryPriority string = "priority"
)

const (
	PermissionName        string = "name"
	PermissionDescription string = "description"
)

const (
	ReviewProductID string = "product_id"
	ReviewUserID    string = "user_id"
	ReviewRating    string = "rating"
	ReviewComment   string = "comment"
)

const (
	RoleRoleName               string = "role_name"
	RoleRoleType               string = "role_type"
	RoleDescription            string = "description"
	RolePermissions            string = "role_permissions"
	RolePermissionsCapitalized string = "Permissions"
)

const (
	WishlistUserID    string = "user_id"
	WishlistProductID string = "product_id"
)

const (
	BrandName        string = "name"
	BrandDescription string = "description"
)

const (
	ColorName string = "name"
)

const (
	AttributeTypeName                  string = "name"
	AttributeOptionAttributeTypeID     string = "attribute_type_id"
	AttributeOptionAttributeOptionName string = "attribute_option_name"
)

const (
	SizeOptionName           string = "name"
	SizeOptionSortOrder      string = "sort_order"
	SizeOptionSizeCategoryID string = "size_category_id"
	SizeCategoryName         string = "name"
)
