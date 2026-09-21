package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Database struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SSLMode     string `mapstructure:"ssl_mode"`
	PoolSize    int    `mapstructure:"pool_size"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
}

type JWT struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	CORSOrigins string `mapstructure:"cors_origins"`
}

type S3Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Config struct {
	Env      string   `mapstructure:"env"`
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
	JWT      JWT      `mapstructure:"jwt"`
	S3       S3Config `mapstructure:"s3"`
}

func (c Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password,
		c.Database.Name, c.Database.SSLMode)
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.Env, "production")
}

type AppConfig struct {
	Config
	S3Client *s3.Client
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("env", "development")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.pool_size", 10)
	v.SetDefault("database.max_idle_conn", 2)

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "3001")
	v.SetDefault("server.cors_origins", "http://localhost:3000")
}

func bindEnv(v *viper.Viper) {

	binds := map[string]string{
		"env": "APP_ENV",

		"database.host":          "DB_HOST",
		"database.port":          "DB_PORT",
		"database.user":          "DB_USER",
		"database.password":      "DB_PASSWORD",
		"database.name":          "DB_NAME",
		"database.ssl_mode":      "DB_SSLMODE",
		"database.pool_size":     "DB_POOL_SIZE",
		"database.max_idle_conn": "DB_MAX_IDLE_CONN",

		"server.host":         "SERVER_HOST",
		"server.port":         "SERVER_PORT",
		"server.cors_origins": "CORS_ORIGINS",

		"jwt.access_secret":  "ACCESS_SECRET",
		"jwt.refresh_secret": "REFRESH_SECRET",

		"s3.access_key": "S3_ACCESS_KEY",
		"s3.secret_key": "S3_SECRET_KEY",
		"s3.region":     "S3_REGION",
		"s3.bucket":     "S3_BUCKET",
		"s3.endpoint":   "S3_ENDPOINT",
	}
	for key, env := range binds {
		_ = v.BindEnv(key, env)
	}
}

func (c Config) validate() error {
	var missing []string
	required := map[string]string{
		"DB_USER":        c.Database.User,
		"DB_PASSWORD":    c.Database.Password,
		"DB_NAME":        c.Database.Name,
		"ACCESS_SECRET":  c.JWT.AccessSecret,
		"REFRESH_SECRET": c.JWT.RefreshSecret,
		"S3_BUCKET":      c.S3.Bucket,
		"S3_ACCESS_KEY":  c.S3.AccessKey,
		"S3_SECRET_KEY":  c.S3.SecretKey,
		"S3_ENDPOINT":    c.S3.Endpoint,
		"S3_REGION":      c.S3.Region,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	if c.IsProduction() && c.JWT.AccessSecret == c.JWT.RefreshSecret {
		return fmt.Errorf("ACCESS_SECRET and REFRESH_SECRET must differ in production")
	}

	return nil
}

func LoadConfig() (*AppConfig, error) {
	var cfg AppConfig

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading configuration from environment")
	}

	v := viper.New()
	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/app")
	v.AddConfigPath("internal/platform/config/")
	if path := strings.TrimSpace(os.Getenv("CONFIG_PATH")); path != "" {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("failed to read config.yaml: %w", err)
		}
		log.Println("no config.yaml found, using defaults and environment")
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnv(v)

	if err := v.Unmarshal(&cfg.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Config.validate(); err != nil {
		return nil, err
	}

	s3Client, err := newS3Client(cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to build s3 client: %w", err)
	}
	cfg.S3Client = s3Client

	return &cfg, nil
}

func newS3Client(cfg Config) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.S3.Endpoint,
			SigningRegion:     cfg.S3.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.S3.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKey,
			cfg.S3.SecretKey,
			"",
		)),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	}), nil
}
