package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/silverspase/todo/internal/config"
)

func Init(appCfg config.Config) *zap.Logger {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(appCfg.LogLevel))
	if err != nil {
		log.Fatal(err)
	}

	var logCfg = zap.NewProductionConfig()
	logCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logCfg.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := logCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger
}
