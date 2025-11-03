package config

type Config struct {
	AppConfig        AppConfig
	DBConfig         DBConfig
	SecretConfig     SecretConfig
	ObjStoreConfig   ObjectStoreConfig
	ExtServiceConfig ExternalServiceConfig
}

type AppConfig struct {
	Port                       string
	DomainURL                  string
	SuperAdminEmail            string
	SuperAdminPassword         string
	SocialLoginDefaultPassword string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	ShowLog  string
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
<<<<<<< HEAD

// type Config struct {
// 	Port                       string
// 	DBHost                     string
// 	DBPort                     string
// 	DBUser                     string
// 	DBPassword                 string
// 	DBName                     string
// 	DBSSLMode                  string
// 	DBShowLog                  string
// 	JWTSecret                  string
// 	SuperAdminEmail            string
// 	SuperAdminPassword         string
// 	DomainURL                  string
// 	FcmServerKey               string
// 	FcmUrl                     string
// 	GoogleClientID             string
// 	FacebookAppID              string
// 	AppHealthCheckToken        string
// 	SocialLoginDefaultPassword string
// 	ObjStore                   ObjectStoreConfig
// }

// type ObjectStoreConfig struct {
// 	Region          string
// 	BucketName      string
// 	AccountID       string
// 	AccessKeyID     string
// 	AccessKeySecret string
// 	PublicDomain    string
// }
=======
>>>>>>> 930864ae816b312ba19eeb0fd7b8763f7726aac3
