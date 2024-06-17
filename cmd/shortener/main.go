package main

import (
	"errors"
	"flag"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/storage"
	"github.com/Mikeloangel/squasher/config"
	"github.com/go-chi/chi/v5"
)

var links *storage.Storage
var conf *config.Config

func main() {
	flag.Parse()

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

	setFlags()
}

func setFlags() {
	flag.StringVar(&conf.HostLocation, "b", conf.HostLocation, "Api host location to get redirect from")
	flag.Func("a", "Sets server location and port in format host:port", func(s string) error {
		var err error

		host := strings.Split(s, ":")

		if len(host) != 2 {
			return errors.New("bad format for -b flag. Expected format host:port")
		}

		conf.ServerLocation = host[0]
		conf.ServerPort, err = strconv.Atoi(host[1])

		if err != nil {
			return err
		}

		return nil
	})
}

func Router(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.Post)
	r.Get("/{id}", handler.Get)

	return r
}
