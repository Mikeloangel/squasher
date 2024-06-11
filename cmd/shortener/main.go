package main

import (
	"fmt"
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
)

var links *storage.Storage

func main() {
	links = storage.NewStorage()

	handler := &handlers.Handler{
		Storage: links,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.Handle)

	dsn := fmt.Sprintf("%s:%s", config.ServerName, config.ServerPort)
	err := http.ListenAndServe(dsn, mux)

	if err != nil {
		panic(err)
	}
}
