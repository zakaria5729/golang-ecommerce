package constants

type StringKey string

const (
	UserContextKey     StringKey = "user"
	UserIDContextKey   StringKey = "user_id"
	ObjStoreProviderR2 string    = "cloudflare_r2"
	ProjectName        string    = "Easy-Commerce"
)

const (
	EnvKeyPort             string = "PORT"
	EnvKeyHost             string = "HOST"
	EnvKeyActiveProfile    string = "ACTIVE_PROFILE"
	EnvActiveObjectStorage string = ObjStoreProviderR2
)

const (
	EnvSuperAdminEmail               string = "SUPER_ADMIN_EMAIL"
	EnvSuperAdminPassword            string = "SUPER_ADMIN_PASSWORD"
	EnvKeyDomainURL                  string = "DOMAIN_URL"
	EnvKeyJWTSecret                  string = "JWT_SECRET"
	EnvKeyFcmUrl                     string = "FCM_URL"
	EnvKeyFcmServerKey               string = "FCM_SERVER_KEY"
	EnvKeyFacebookAppID              string = "FACEBOOK_APP_ID"
	EnvKeyGoogleClientID             string = "GOOGLE_CLIENT_ID"
	EnvKeyAppHealthCheckToken        string = "APP_HEALTH_CHECK_TOKEN"
	EnvKeySocialLoginDefaultPassword string = "SOCIAL_LOGIN_DEFAULT_PASSWORD"
)

const (
	AuthTypeGoogle      string = "google"
	AuthTypeFacebook    string = "facebook"
	GoogleUserInfoURL   string = "https://www.googleapis.com/oauth2/v2/userinfo"
	FacebookUserInfoURL string = "https://graph.facebook.com/me?fields=email,name,first_name,last_name,picture.width(250).height(250)"
)

const (
	Bearer        string = "Bearer"
	Authorization string = "Authorization"
)

const (
	AccessTokenExpiryHours        int = 24
	RefreshTokenExpiryHours       int = 48
	VerificationTokenExpiryHours  int = 6
	PasswordResetTokenExpiryHours int = 6
)

const (
	DBMaxIdleConns               int   = 5
	DBMaxOpenConns               int   = 20
	DBConnMaxLifeTimeHour        int64 = 1
	DBConnMaxIdleTimeMinute      int   = 15
	LogFileDeleteProhibitedLimit int   = 2
)

const (
	SizeInMB          int64  = 1024 * 1024 // 1MB
	AllowedImageTypes string = "image/jpeg,image/jpg,image/png,image/gif,image/svg+xml"
	AllowedDocTypes   string = "application/pdf"
	LogFileFormat     string = "2006-01-02"
	LogFolderName     string = "logs"
	LogFileExt        string = "jsonl"
)

const (
	EnvLocal string = "local"
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
	RoleTypeSuperAdmin string = "SUPER_ADMIN"
	RoleTypeAdmin      string = "ADMIN"
	RoleTypeMaintainer string = "MAINTAINER"
	RoleTypeSeller     string = "SELLER"
	RoleTypeUser       string = "USER"
)

const (
	RoleNameSuperAdmin string = "Super Admin"
	RoleNameAdmin      string = "Admin"
	RoleNameMaintainer string = "Maintainer"
	RoleNameSeller     string = "Seller"
	RoleNameUser       string = "User"
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

const (
	PhoneNumRegex string = `^(88)?01\d{9}$`
	EmailRegex    string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)
