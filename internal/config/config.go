package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	PoolSize string `mapstructure:"pool_size"`
}

type JWT struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type S3Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
	JWT      JWT      `mapstructure:"jwt"`
	S3       S3Config `mapstructure:"s3"`
}

func (c Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

type AppConfig struct {
	Config
	S3Client *s3.Client
}

func LoadConfig() (*AppConfig, error) {
	var cfg AppConfig
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/config/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return &cfg, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("jwt.access_secret", "ACCESS_SECRET")
	viper.BindEnv("jwt.refresh_secret", "REFRESH_SECRET")

	viper.BindEnv("s3.access_key", "S3_ACCESS_KEY")
	viper.BindEnv("s3.secret_key", "S3_SECRET_KEY")
	viper.BindEnv("s3.region", "S3_REGION")
	viper.BindEnv("s3.bucket", "S3_BUCKET")
	viper.BindEnv("s3.endpoint", "S3_ENDPOINT")

	if err := viper.Unmarshal(&cfg.Config); err != nil {
		return &cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.Name == "" {
		return &cfg, fmt.Errorf("database user, password, and name are required")
	}
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
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	cfg.S3Client = s3Client
	return &cfg, nil
}
