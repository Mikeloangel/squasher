package main

import (
	"fmt"
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
	"github.com/go-chi/chi/v5"
)

var links *storage.Storage

func main() {
	links = storage.NewStorage()
	var handler = &handlers.Handler{
		Storage: links,
	}

	dsn := fmt.Sprintf("%s:%s", config.ServerName, config.ServerPort)
	err := http.ListenAndServe(dsn, Router(handler))

	if err != nil {
		panic(err)
	}
}

func Router(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.Post)
	r.Get("/{id}", handler.Get)

	return r
}
