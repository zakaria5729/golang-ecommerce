package config

type Config struct {
	AppConfig        AppConfig
	DBConfig         DBConfig
	SecretConfig     SecretConfig
	ObjStoreConfig   ObjectStoreConfig
	ExtServiceConfig ExternalServiceConfig
}

type AppConfig struct {
	Host                       string
	Port                       string
	DomainURL                  string
	UuidType                   string
	WriteLogWhen               string
	DBStatsLogType             string
	SuperAdminEmail            string
	SuperAdminPassword         string
	SocialLoginDefaultPassword string
}

type DBConfig struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	DBShowConsoleLog string
}

type SecretConfig struct {
	JWTSecret           string
	AppHealthCheckToken string
}

type ExternalServiceConfig struct {
	FcmServerKey   string
	FcmUrl         string
	GoogleClientID string
	FacebookAppID  string
}

type ObjectStoreConfig struct {
	Region          string
	BucketName      string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	PublicDomain    string
}
