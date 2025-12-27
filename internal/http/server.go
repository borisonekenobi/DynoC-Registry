package http

import (
	"context"
	"dynoc-registry/internal/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer() http.Handler {
	ctx := context.Background()

	pool, err := db.NewPool(ctx)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx = context.WithValue(r.Context(), "db", pool)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	registerRoutes(r)
	return r
}
