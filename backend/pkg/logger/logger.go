// internal/logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(env string) error {
	var config zap.Config

	if env == "prod" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		// Добавляем stack trace только для ERROR и выше (не для WARN)
		config.Development = true
		config.DisableStacktrace = false
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Stack trace только для ERROR и PANIC, не для WARN
	var err error
	Log, err = config.Build(zap.AddStacktrace(zapcore.ErrorLevel))

	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	_ = Log.Sync()
}
