package handlers

import (
	"context"
	"net/http"
	"time"
)

// Ping handles a test connection to db sets 500 status if no connection, 200 if connection is ok
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
