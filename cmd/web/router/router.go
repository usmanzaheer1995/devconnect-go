package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres"
)

// NewRouter returns new instance of router
func NewRouter(services *postgres.Services) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/auth", userRoutes(services.User))

	return r
}
