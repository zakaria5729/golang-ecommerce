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
	DomainURL          string
	FcmServerKey       string
	FcmUrl             string
	GoogleClientID     string
	FacebookAppID      string
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
