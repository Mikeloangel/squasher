package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
)

type Handler struct {
	Storage *storage.Storage
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Method: %s, URL: %s, Header: %v", r.Method, r.URL.Path, r.Header)

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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("localhost:8080/" + shortened))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	shortUrl := strings.TrimPrefix(r.URL.Path, "/")

	url, err := h.Storage.Get(shortUrl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
