package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/internal/models"
)

// CreateShortURL handles the creation of a shortened URL.
// It reads the URL from the request body, generates a shortened version,
// and returns it to the client.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	h.handleCreateShortURL(w, r, false)
}

// CreateShortURLJson handles the creation of a shortened URL.
// uses json format for request and response
// It reads the URL from the request body, generates a shortened version,
// and returns it to the client.
func (h *Handler) CreateShortURLJson(w http.ResponseWriter, r *http.Request) {
	h.handleCreateShortURL(w, r, true)
}

// handleCreateShortURL is a helper function for creating a link depending on request type: JSON or plaint/text
func (h *Handler) handleCreateShortURL(w http.ResponseWriter, r *http.Request, isJSON bool) {
	var url string

	// parses request
	if isJSON {
		var request models.CreateShortURLRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := models.ValidateCreateShortURLRequest(&request); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
			return
		}

		url = request.URL
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Body error", http.StatusBadRequest)
			return
		}
		url = string(body)
	}

	if len(strings.TrimSpace(url)) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// adds link
	shortened := h.Links.Set(url)
	result := h.Conf.GetHostLocation() + shortened.Shorten

	// sends response
	if isJSON {
		response := models.CreateShortURLResponse{Result: result}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(result))
	}
}
