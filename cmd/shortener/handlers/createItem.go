package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/models"
)

// CreateShortURL handles the creation of a shortened URL.
// It reads the URL from the request body, generates a shortened version,
// and returns it to the client.
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var url string

	// parses request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}
	url = string(body)

	if len(strings.TrimSpace(url)) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	status := http.StatusCreated

	// adds link
	shortened, err := h.Storage.StoreURL(url)
	var ae *storage.ItemAlreadyExistsError
	if errors.As(err, &ae) {
		status = http.StatusConflict
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := h.Conf.GetHostLocation() + shortened.Shorten

	// sends response
	w.WriteHeader(status)
	w.Write([]byte(result))

}

// CreateShortURLJson handles the creation of a shortened URL.
// uses json format for request and response
// It reads the URL from the request body, generates a shortened version,
// and returns it to the client.
func (h *Handler) CreateShortURLJson(w http.ResponseWriter, r *http.Request) {
	var url string

	// parses request
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

	if len(strings.TrimSpace(url)) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// adds link
	status := http.StatusCreated
	shortened, err := h.Storage.StoreURL(url)
	var ae *storage.ItemAlreadyExistsError
	if errors.As(err, &ae) {
		status = http.StatusConflict
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := h.Conf.GetHostLocation() + shortened.Shorten

	// sends response
	response := models.CreateShortURLResponse{Result: result}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

}
