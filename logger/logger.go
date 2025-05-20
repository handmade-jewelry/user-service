package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	var err error
	log, err = zap.NewProduction()
	return err
}

func Sync() {
	log.Sync()
}

func Info(msg string, key, val string) {
	log.Info(msg, zap.String(key, val))
}

func Warn(msg string, key, val string) {
	log.Warn(msg, zap.String(key, val))
}

func Error(msg string, err error) {
	log.Error(msg, zap.Error(err))
}

func Fatal(msg string, err error) {
	log.Fatal(msg, zap.Error(err))
}
