package main

import (
	"flag"
	"os"

	"github.com/Mikeloangel/squasher/cmd/shortener/state"
)

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
