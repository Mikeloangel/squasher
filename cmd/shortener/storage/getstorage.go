package storage

import (
	"database/sql"
	"strings"

	"github.com/Mikeloangel/squasher/internal/config"
	"github.com/Mikeloangel/squasher/internal/logger"
)

// NewStorageFactory creates correct storage implementation depending on config
func NewStorageFactory(cfg *config.Config, db *sql.DB) Storager {
	// if DBDSN is not empty then use db storager
	if len(strings.TrimSpace(cfg.DBDSN)) > 0 {
		logger.Info("Using db storage!")
		return NewDBStorage(db)
	}

	// not empty file storage location returns FileStorage
	if len(strings.TrimSpace(cfg.StorageFileLocation)) > 0 {
		logger.Info("Using file storage!")
		return NewFileStorage(cfg)
	}

	// otherwise returns in memory storage
	logger.Info("In memory storage mode is on")
	return NewInMemoryStorage()
}
