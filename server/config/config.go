package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl    string `env:"DATABASE_URL"`
	RedisUrl       string `env:"REDIS_URL"`
	Port           string `env:"PORT"`
	SessionSecret  string `env:"SECRET"`
	Domain         string `env:"DOMAIN"`
	CorsOrigin     string `env:"CORS_ORIGIN"`
	AccessKey      string `env:"AWS_ACCESS_KEY"`
	SecretKey      string `env:"SECRET_KEY"`
	BucketName     string `env:"BUCKET_NAME"`
	Region         string `env:"REGION"`
	GmailUser      string `env:"GMAIL_USER"`
	GmailPassword  string `env:"GMAIL_PASSWORD"`
	HandlerTimeOut int64  `env:"HANDLER_TIMEOUT"`
	MaxBodyBytes   int64  `env:"MAX_BODY_BYTES"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	if err != nil {
		return
	}

	return
}
