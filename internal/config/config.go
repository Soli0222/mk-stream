package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持する
type Config struct {
	Host  string
	Token string
}

// LoadConfig は.envファイルと環境変数から設定を読み込む
func LoadConfig(logger *slog.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("Not using .env file")
	}

	requiredVars := []string{
		"HOST",
		"TOKEN",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return nil, errors.New("required environment variable is not set: " + v)
		}
	}

	return &Config{
		Host:  os.Getenv("HOST"),
		Token: os.Getenv("TOKEN"),
	}, nil
}
