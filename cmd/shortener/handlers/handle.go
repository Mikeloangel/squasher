// Package handlers provides HTTP handlers for the URL shortener service.
package handlers

import (
	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

// Handler embeds the application state and provides methods to handle HTTP requests.
type Handler struct {
	Links storage.Storager
	Conf  *config.Config
}

// NewHandler creates a new Handler with the given application state.
func NewHandler(appState state.State) *Handler {
	return &Handler{
		Links: appState.Links,
		Conf:  appState.Conf,
	}
}
