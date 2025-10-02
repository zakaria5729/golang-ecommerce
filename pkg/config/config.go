package config

type Config struct {
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	DBShowLog          string
	JWTSecret          string
	SuperAdminEmail    string
	SuperAdminPassword string
	ObjStore           ObjectStoreConfig
}

type ObjectStoreConfig struct {
	Region          string
	BucketName      string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	PublicDomain    string
}
