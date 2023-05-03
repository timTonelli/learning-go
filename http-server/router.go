package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createRouter() *chi.Mux {
	// Create new Chi router
	router := chi.NewRouter()

	// Register middlewares
	router.Use(
		middleware.Logger,
		middleware.Heartbeat("/health"),
	)

	// Register endpoints
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from endpoint: %s", r.URL.Path)
	})

	return router
}
