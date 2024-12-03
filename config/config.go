package config

import (
	"effectiveMobileTask/lib/logger"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
)

func LoadConfigEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Error("not found .env file", slog.Any("err", err))
		log.Fatal("Error loading .env file")
	}
}
