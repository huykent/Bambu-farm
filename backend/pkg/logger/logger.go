package logger

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger(env string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	if env == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger.Sugar()
}
