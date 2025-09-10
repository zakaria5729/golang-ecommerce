package constants

const (
	TableUser            = "users"
	TableRole            = "roles"
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
)

const (
	EnvKeyPort       = "PORT"
	EnvKeyDBHost     = "DB_HOST"
	EnvKeyDBPort     = "DB_PORT"
	EnvKeyDBUser     = "DB_USER"
	EnvKeyDBPassword = "DB_PASSWORD"
	EnvKeyDBName     = "DB_NAME"
	EnvKeyDBSSLMode  = "DB_SSL_MODE"
)
