// Package singleton logger add a logging functionality to app
package logger

import (
	"go.uber.org/zap"
)

// Log is the global logger instance.
var Log *zap.Logger = zap.NewNop()

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

	Log = procuctionLogger
	return nil
}
