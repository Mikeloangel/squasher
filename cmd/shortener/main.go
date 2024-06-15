package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/config"
	"github.com/go-chi/chi/v5"
)

var links *storage.Storage
var conf *config.Config

func main() {
	parseEnviroment()

	var handler = &handlers.Handler{
		Storage: links,
		Config:  conf,
	}

	err := http.ListenAndServe(conf.GetServerConnectionString(), Router(handler))

	if err != nil {
		panic(err)
	}
}

func init() {
	links = storage.NewStorage()
	conf = config.NewConfig()
}

func parseEnviroment() {
	flag.StringVar(&conf.HostLocation, "b", conf.HostLocation, "Api host location to get redirect from")
	flag.Func("a", "Sets server location and port in format host:port", conf.SetServerFromString)
	flag.Parse()

	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		conf.HostLocation = baseUrl
	}

	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		err := conf.SetServerFromString(serverAddr)
		if err != nil {
			panic(err)
		}
	}
}

func Router(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.Post)
	r.Get("/{id}", handler.Get)

	return r
}
