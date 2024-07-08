package main

import (
	"github.com/Mikeloangel/squasher/cmd/shortener/handlers"
	"github.com/Mikeloangel/squasher/cmd/shortener/middlewares"
	"github.com/Mikeloangel/squasher/internal/logger"
	"github.com/go-chi/chi/v5"
)

// Router sets up the HTTP routes
func Router(handler *handlers.Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.GzipMiddleware)
	r.Use(logger.WithLoggerMiddleware)

	registerShortURLRoutes(r, handler)
	registerAPIRoutes(r, handler)

	return r
}

// registerAPIRoutes registers API routes
func registerAPIRoutes(r chi.Router, handler *handlers.Handler) {
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", handler.CreateShortURLJson)
		r.Post("/batch", handler.CreateBatchUrls)
	})
}

// registerShortURLRoutes registers general app routes
func registerShortURLRoutes(r chi.Router, handler *handlers.Handler) {
	r.Post("/", handler.CreateShortURL)
	r.Get("/{id}", handler.GetOriginalURL)
	r.Get("/ping", handler.Ping)
}
