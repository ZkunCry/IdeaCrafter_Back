package config

import (
	"fmt"
	"log"
	"strings"

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

type Server struct {
    Host string `mapstructure:"host"`
    Port string `mapstructure:"port"`
}

type Config struct {
    Database Database `mapstructure:"database"`
    Server   Server   `mapstructure:"server"`
}

func (c Config) DBConnectionString() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

func LoadConfig() (Config, error) {
    var cfg Config
		if err:= godotenv.Load(); err != nil {
		log.Println("no .env file found")
		}
    // читаем config.yaml
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("internal/config/")
    viper.AddConfigPath(".")
    if err := viper.ReadInConfig(); err != nil {
        return cfg, fmt.Errorf("failed to read config.yaml: %w", err)
    }

    // связываем ENV
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    viper.AutomaticEnv()

    // биндим конкретные env
    viper.BindEnv("database.user", "DB_USER")
    viper.BindEnv("database.password", "DB_PASSWORD")
    viper.BindEnv("database.name", "DB_NAME")

    if err := viper.Unmarshal(&cfg); err != nil {
        return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    if cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.Name == "" {
        return cfg, fmt.Errorf("database user, password, and name are required")
    }

    return cfg, nil
}
