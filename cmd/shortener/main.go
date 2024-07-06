// Package main provides an entry point for the URL shortener service
package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"

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
		"",
	)

	// Parses enviroment flags and command line flags
	err = config.ParseEnvironment(cfg)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// Get db
	db := config.GetDB(cfg)
	defer db.Close()

	// Gets storage implementation
	storage := storage.NewStorageFactory(cfg, db)

	// Inits storage
	err = storage.Init()
	if err != nil {
		logger.Fatal(err)
		return
	}

	// Initializes application state
	appState := state.NewState(storage, cfg, db)

	// Creates a new handler for application state
	handler := handlers.NewHandler(appState)

	// Starts HTTP server
	err = http.ListenAndServe(appState.Conf.GetServerConnectionString(), Router(handler))

	if err != nil {
		logger.Fatal(err)
	}
}
