package storage

import (
	"strings"

	"github.com/Mikeloangel/squasher/internal/config"
)

// NewStorageFactory creates correct storage implementation depending on config
func NewStorageFactory(cfg *config.Config) Storager {
	// not empty file storage location returns FileStorage
	if len(strings.TrimSpace(cfg.StorageFileLocation)) > 0 {
		return NewFileStorage(cfg)
	}

	// otherwise returns in memory storage
	return NewInMemoryStorage()
}
