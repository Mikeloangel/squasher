// Package state provides the application state management for the URL shortener service.
package state

import (
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

// State holds the application state, including configuration and storage.
type State struct {
	// Links is the storage interface for managing shortened URLs.
	Links storage.Storager
	// Conf holds the application configuration.
	Conf *config.Config
}

// NewState creates a new instance of State with the provided storage and configuration.
func NewState(links storage.Storager, conf *config.Config) State {
	return State{
		Links: links,
		Conf:  conf,
	}
}
