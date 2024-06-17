package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/config"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage *storage.Storage
	Config  *config.Config
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(string(body))) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	shortened := h.Storage.Set(string(body))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(h.Config.GetHostLocation() + shortened))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "id")
	url, err := h.Storage.Get(t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
