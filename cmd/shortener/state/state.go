// Package state provides the application state management for the URL shortener service.
package state

import (
	"database/sql"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

// State holds the application state, including configuration and storage.
type State struct {
	// Storage is the storage interface for managing shortened URLs.
	Storage storage.Storager
	// Conf holds the application configuration.
	Conf *config.Config
	// database connection
	DB *sql.DB
}

// NewState creates a new instance of State with the provided storage and configuration.
func NewState(storage storage.Storager, conf *config.Config, db *sql.DB) State {
	return State{
		Storage: storage,
		Conf:    conf,
		DB:      db,
	}
}
