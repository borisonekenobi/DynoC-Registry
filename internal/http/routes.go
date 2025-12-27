package http

import (
	"dynoc-registry/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handlers.Login)
		r.Get("/renew", handlers.RenewToken)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", handlers.CreateAccount)
		r.Get("/{name}", handlers.ReadAccount)
		r.Put("/{name}", handlers.UpdateAccount)
		r.Delete("/{name}", handlers.DeleteAccount)
	})

	r.Route("/packages", func(r chi.Router) {
		r.Post("/", handlers.CreatePackage)
		r.Post("/{name}/versions", handlers.CreatePackageVersion)

		r.Get("/{name}/latest", handlers.ReadLatest)
		r.Get("/{name}/versions", handlers.ReadVersions)
		r.Get("/{name}/{version}", handlers.ReadPackage)

		r.Put("/{name}", handlers.UpdatePackage)
		r.Put("/{name}/versions/{version}", handlers.UpdatePackageVersion)

		r.Delete("/{name}", handlers.DeletePackage)
		r.Delete("/{name}/versions/{version}", handlers.DeletePackageVersion)
	})

	r.Get("/search", handlers.FindPackages)
}
