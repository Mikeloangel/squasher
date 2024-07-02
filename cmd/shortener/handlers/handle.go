// Package handlers provides HTTP handlers for the URL shortener service.
package handlers

import (
	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

// Handler embeds the application state and provides methods to handle HTTP requests.
type Handler struct {
	Storage storage.Storager
	Conf    *config.Config
}

// NewHandler creates a new Handler with the given application state.
func NewHandler(appState state.State) *Handler {
	return &Handler{
		Storage: appState.Storage,
		Conf:    appState.Conf,
	}
}
