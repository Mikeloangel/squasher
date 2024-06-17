// Package main provides an entry point for the URL shortener service
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/state"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/config"
	"github.com/go-chi/chi/v5"
)

// main is the entry point to application
func main() {
	var err error

	// Initializes application state
	appState := state.NewState(
		storage.NewStorage(),
		config.NewConfig("localhost", 8080, "http://localhost:8080"),
	)

	// Parses enviroment flags and command line flags
	err = parseEnviroment(appState)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Creates a new handler for application state
	handler := handlers.NewHandler(appState)

	// Starts HTTP server
	err = http.ListenAndServe(appState.Conf.GetServerConnectionString(), Router(handler))

	if err != nil {
		log.Fatal(err)
	}
}

// parseEnviroment parses environment variables and command-line flags
// to configure the application state.
func parseEnviroment(state state.State) error {
	flag.StringVar(&state.Conf.HostLocation, "b", state.Conf.HostLocation, "Api host location to get redirect from")
	flag.Func("a", "Sets server location and port in format host:port", state.Conf.ParseServerConfig)
	flag.Parse()

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		state.Conf.HostLocation = baseURL
	}

	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		err := state.Conf.ParseServerConfig(serverAddr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Router sets up the HTTP routes
func Router(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.CreateShortURL)
	r.Get("/{id}", handler.GetOriginalURL)

	return r
}
