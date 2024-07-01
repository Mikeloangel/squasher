package main

import (
	"flag"
	"os"

	"github.com/Mikeloangel/squasher/internal/config"
)

// parseEnviroment parses environment variables and command-line flags
// to configure the application state.
func parseEnviroment(Conf *config.Config) error {
	flag.StringVar(&Conf.HostLocation, "b", Conf.HostLocation, "Api host location to get redirect from")
	flag.StringVar(&Conf.StorageFileLocation, "f", Conf.StorageFileLocation, "Storage file location, if empty uses in memory handling")
	flag.Func("a", "Sets server location and port in format host:port", Conf.ParseServerConfig)
	flag.Parse()

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		Conf.HostLocation = baseURL
	}

	if storageFileLoaction := os.Getenv("FILE_STORAGE_PATH"); storageFileLoaction != "" {
		Conf.StorageFileLocation = storageFileLoaction
	}

	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		err := Conf.ParseServerConfig(serverAddr)
		if err != nil {
			return err
		}
	}

	return nil
}
