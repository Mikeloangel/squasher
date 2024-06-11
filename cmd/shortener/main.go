package main

import (
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
)

var links *storage.Storage

func main() {
	links = storage.NewStorage()

	handler := &handlers.Handler{
		Storage: links,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.Handle)
	// mux.Handle("/", middlewares.TextPlain(http.HandlerFunc(handler.Handle)))

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
