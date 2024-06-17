package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	state.State
}

func NewHandler(appState state.State) *Handler {
	return &Handler{
		State: appState,
	}
}

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
