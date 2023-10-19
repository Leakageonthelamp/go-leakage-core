package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const EnvFileName = ".env"
const EnvTestFileName = "test.env"

type IENV interface {
	Config() *ENVConfig
	IsDev() bool
	IsTest() bool
	IsMock() bool
	IsProd() bool
	Bool(key string) bool
	Int(key string) int
	String(key string) string
	All() map[string]string
}

type ENVConfig struct {
	LogLevel logrus.Level
	LogHost  string `mapstructure:"log_host"`
	LogPort  string `mapstructure:"log_port"`

	Host    string `mapstructure:"host"`
	ENV     string `mapstructure:"env"`
	Service string `mapstructure:"service"`

	SentryDSN string `mapstructure:"sentry_dsn"`

	DBDriver   string `mapstructure:"db_driver"`
	DBHost     string `mapstructure:"db_host"`
	DBName     string `mapstructure:"db_name"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBPort     string `mapstructure:"db_port"`

	DBMongoHost     string `mapstructure:"db_mongo_host"`
	DBMongoName     string `mapstructure:"db_mongo_name"`
	DBMongoUserName string `mapstructure:"db_mongo_username"`
	DBMongoPassword string `mapstructure:"db_mongo_password"`
	DBMongoPort     string `mapstructure:"db_mongo_port"`

	MQHost     string `mapstructure:"mq_host"`
	MQUser     string `mapstructure:"mq_user"`
	MQPassword string `mapstructure:"mq_password"`
	MQPort     string `mapstructure:"mq_port"`

	CachePort string `mapstructure:"cache_port"`
	CacheHost string `mapstructure:"cache_host"`

	ABCIEndpoint      string `mapstructure:"abci_endpoint"`
	DIDMethodDefault  string `mapstructure:"did_method_default"`
	DIDKeyTypeDefault string `mapstructure:"did_key_type_default"`

	WinRMHost     string `mapstructure:"winrm_host"`
	WinRMUser     string `mapstructure:"winrm_user"`
	WinRMPassword string `mapstructure:"winrm_password"`

	S3Endpoint  string `mapstructure:"s3_endpoint"`
	S3AccessKey string `mapstructure:"s3_access_key"`
	S3SecretKey string `mapstructure:"s3_secret_key"`
	S3Bucket    string `mapstructure:"s3_bucket"`
	S3Region    string `mapstructure:"s3_region"`
	S3IsHTTPS   bool   `mapstructure:"s3_https"`

	EmailServer   string `mapstructure:"email_server"`
	EmailPort     int    `mapstructure:"email_port"`
	EmailUsername string `mapstructure:"email_username"`
	EmailPassword string `mapstructure:"email_password"`
	EmailSender   string `mapstructure:"email_sender"`
}

type ENVType struct {
	config *ENVConfig
}

func NewEnv() IENV {
	return NewENVPath(".")
}

func NewENVPath(path string) IENV {
	if os.Getenv("APP_ENV") == "test" {
		viper.SetConfigName(EnvTestFileName)
	} else {
		viper.SetConfigName(EnvFileName)
	}

	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	envKeys := []string{
		"LOG_HOST",
		"LOG_PORT",
		"HOST", "ENV", "SERVICE",
		"SENTRY_DSN",
		"DB_DRIVER", "DB_HOST", "DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_PORT",
		"DB_MONGO_HOST", "DB_MONGO_NAME", "DB_MONGO_USERNAME", "DB_MONGO_PASSWORD", "DB_MONGO_PORT",
		"MQ_HOST", "MQ_USER", "MQ_PASSWORD", "MQ_PORT",
		"WINRM_HOST", "WINRM_USER", "WINRM_PASSWORD", "WINRM_PORT",
		"CACHE_PORT", "CACHE_HOST",
		"ABCI_ENDPOINT", "DID_METHOD_DEFAULT", "DID_KEY_TYPE_DEFAULT", "S3_ENDPOINT",
		"S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET", "S3_HTTPS", "S3_REGION",
		"EMAIL_SERVER", "EMAIL_PORT", "EMAIL_USERNAME", "EMAIL_PASSWORD", "EMAIL_SENDER",
	}

	for _, key := range envKeys {
		viper.BindEnv(key)
	}

	env := &ENVConfig{}
	err := viper.Unmarshal(env)
	if err != nil {
		// NewLoggerSimple().Debug(err.Error())
		panic(err)
	}

	env.LogLevel, _ = logrus.ParseLevel(viper.GetString("log_level"))
	return &ENVType{
		config: env,
	}
}

func (e ENVType) Config() *ENVConfig {
	return e.config
}

// IsDev config  is Dev config
func (e ENVType) IsDev() bool {
	return e.String("env") == "dev"
}

func (e ENVType) IsMock() bool {
	return e.String("env") == "mock"
}

// IsTest config  is Test config
func (e ENVType) IsTest() bool {
	return e.String("env") == "test"
}

// IsProd config  is production config
func (e ENVType) IsProd() bool {
	return e.String("env") == "prod"
}

func (e ENVType) Bool(key string) bool {
	return viper.GetBool(strings.ToLower(key))
}

func (e ENVType) Int(key string) int {
	return viper.GetInt(strings.ToLower(key))
}

func (e ENVType) String(key string) string {
	return viper.GetString(strings.ToLower(key))
}
func (e ENVType) All() map[string]string {
	mapEnvs := make(map[string]string, 0)
	envs := viper.AllSettings()
	for key, value := range envs {
		mapEnvs[key] = fmt.Sprintf("%v", value)
	}

	return mapEnvs
}
