package config

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv          string   `env:"APP_ENV" envDefault:"development"`
	Port            string   `env:"APP_PORT" envDefault:"8080"`
	APIKey          string   `env:"API_KEY,required"`
	LogLevel        string   `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat       string   `env:"LOG_FORMAT" envDefault:"json"`
	LogOutput       string   `env:"LOG_OUTPUT" envDefault:"stdout"`
	AllowedOrigins  []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`
	MaxUploadSizeMB int64    `env:"MAX_UPLOAD_SIZE_MB" envDefault:"50"`
	AllowedMIMEs    []string `env:"ALLOWED_MIME_TYPES" envSeparator:","`

	MinioEndpoint  string `env:"MINIO_ENDPOINT,required"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY,required"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY,required"`
	MinioUseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	MinioRegion    string `env:"MINIO_REGION" envDefault:"us-east-1"`
	MinioBucket    string `env:"MINIO_DEFAULT_BUCKET" envDefault:""`
}

func Load() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	for i, item := range cfg.AllowedMIMEs {
		cfg.AllowedMIMEs[i] = strings.TrimSpace(item)
	}

	return cfg, nil
}
