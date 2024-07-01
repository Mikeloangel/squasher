// Package main provides an entry point for the URL shortener service
package main

import (
	"log"
	"net/http"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/internal/config"
	"github.com/Mikeloangel/squasher/internal/logger"
)

// main is the entry point to application
func main() {
	var err error

	// Initializing app logger
	logger.Init("info")

	// Inits config
	cfg := config.NewConfig(
		"localhost",
		8080,
		"http://localhost:8080",
		"/tmp/short-url-db.json",
	)

	// Parses enviroment flags and command line flags
	err = parseEnviroment(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Gets storage implementation
	storage := storage.NewStorageFactory(cfg)

	// Inits storage
	err = storage.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initializes application state
	appState := state.NewState(storage, cfg)

	// Creates a new handler for application state
	handler := handlers.NewHandler(appState)

	// Starts HTTP server
	err = http.ListenAndServe(appState.Conf.GetServerConnectionString(), Router(handler))

	if err != nil {
		log.Fatal(err)
	}
}
