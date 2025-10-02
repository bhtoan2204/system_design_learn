package config

import (
	"event_sourcing_bank_system_gateway/package/settings"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitLoadConfig() (*settings.Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	if env != "production" {
		if err := godotenv.Load(
			fmt.Sprintf(".env.%s", env),
		); err != nil {
			panic(fmt.Errorf("error loading .env files: %w", err))
		}
	}

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnv(v)

	var config settings.Config
	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode configuration: %w", err))
	}

	return &config, nil
}

func bindEnv(v *viper.Viper) {
	// Set up mappings for environment variables to configuration structure
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("server.mode", "SERVER_MODE")

	// Redis mappings
	v.BindEnv("redis.host", "REDIS_HOST")
	v.BindEnv("redis.port", "REDIS_PORT")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("redis.database", "REDIS_DATABASE")

	// Log mappings
	v.BindEnv("log.log_level", "LOG_LOG_LEVEL")
	v.BindEnv("log.file_path", "LOG_FILE_PATH")
	v.BindEnv("log.max_size", "LOG_MAX_SIZE")
	v.BindEnv("log.max_backups", "LOG_MAX_BACKUPS")
	v.BindEnv("log.max_age", "LOG_MAX_AGE")
	v.BindEnv("log.compress", "LOG_COMPRESS")

	// Security mappings
	v.BindEnv("security.jwt_access_secret", "SECURITY_JWT_ACCESS_SECRET")
	v.BindEnv("security.jwt_refresh_secret", "SECURITY_JWT_REFRESH_SECRET")
	v.BindEnv("security.jwt_access_expiration", "SECURITY_JWT_ACCESS_EXPIRATION")
	v.BindEnv("security.jwt_refresh_expiration", "SECURITY_JWT_REFRESH_EXPIRATION")
	v.BindEnv("security.hmac_secret", "SECURITY_HMAC_SECRET")

	// Service mappings
	v.BindEnv("service.payment_service_url", "PAYMENT_SERVICE_URL")
}
