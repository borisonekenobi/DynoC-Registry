package http

import (
	"dynoc-registry/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router) {
	r.Get("/health", handlers.Health)

	// Future:
	// r.Route("/packages", func(r chi.Router) {
	//     r.Post("/", publishPackage)
	//     r.Get("/{name}", getPackage)
	// })
}
