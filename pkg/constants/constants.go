package constants

type StringKey string

const (
	UserContextKey             StringKey = "user"
	UserIDContextKey           StringKey = "user_id"
	ObjStoreProviderR2         string    = "cloudflare_r2"
	ProjectName                string    = "Easy-Commerce"
	AppHealthCheckToken        string    = "abc3WZ29@ld!43~r3_ew*qT#yz"
	SocialLoginDefaultPassword string    = "social@login#password"
)

const (
	EnvKeyActiveProfile    string = "ACTIVE_PROFILE"
	EnvActiveObjectStorage string = ObjStoreProviderR2
)

const (
	EnvKeyDomainURL      string = "DOMAIN_URL"
	EnvKeyJWTSecret      string = "JWT_SECRET"
	EnvKeyFcmUrl         string = "FCM_URL"
	EnvKeyFcmServerKey   string = "FCM_SERVER_KEY"
	EnvKeyGoogleClientID string = "GOOGLE_CLIENT_ID"
	EnvKeyFacebookAppID  string = "FACEBOOK_APP_ID"
)

const (
	AccessTokenExpiryHours        int = 24
	RefreshTokenExpiryHours       int = 48
	VerificationTokenExpiryHours  int = 6
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
	AuthTypeGoogle   string = "google"
	AuthTypeFacebook string = "facebook"
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
