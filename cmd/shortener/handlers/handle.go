// Package handlers provides HTTP handlers for the URL shortener service.
package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/go-chi/chi/v5"
)

// Handler embeds the application state and provides methods to handle HTTP requests.
type Handler struct {
	state.State
}

// NewHandler creates a new Handler with the given application state.
func NewHandler(appState state.State) *Handler {
	return &Handler{
		State: appState,
	}
}

// CreateShortURL handles the creation of a shortened URL.
// It reads the URL from the request body, generates a shortened version,
// and returns it to the client.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(string(body))) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	shortened := h.Links.Set(string(body))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Conf.GetHostLocation() + shortened))
}

// GetOriginalURL handles the retrieval of the original URL for a given shortened version.
// It reads the shortened URL from the request path, retrieves the original URL,
// and redirects the client to the original URL.
func (h *Handler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "id")
	url, err := h.Links.Get(t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
