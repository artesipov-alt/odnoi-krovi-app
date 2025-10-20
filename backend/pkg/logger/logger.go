// internal/logger/logger.go
package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(env string) error {
	var err error

	if env == "prod" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	_ = Log.Sync()
}
