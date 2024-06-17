package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

type Handler struct {
	Storage *storage.Storage
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	shortened := h.Storage.Set(string(body))
	host := fmt.Sprintf("%s://%s:%s/", config.ServerProtocol, config.ServerName, config.ServerPort)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(host + shortened))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	url, err := h.Storage.Get(shortURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
