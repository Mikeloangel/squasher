package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

	w.Header().Set("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
