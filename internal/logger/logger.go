// Package logger provides adaptor implementation for app logger
package logger

import (
	"go.uber.org/zap"
)

// Log is the global logger instance.
var log *zap.SugaredLogger = zap.NewNop().Sugar()

// Init initializes the logger with the given level.
func Init(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	procuctionLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	log = procuctionLogger.Sugar()
	return nil
}

func Debug(args ...interface{}) {
	log.Debugln(args...)
}

func Info(args ...interface{}) {
	log.Infoln(args)
}

func Error(args ...interface{}) {
	log.Errorln(args)
}

func Fatal(args ...interface{}) {
	log.Fatalln(args...)
}
