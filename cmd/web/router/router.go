package router

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	errors2 "github.com/usmanzaheer1995/devconnect-go-v2/internal/errors"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"log"
	"net/http"
)

type rootHandler func(http.ResponseWriter, *http.Request) error

func errorHandler(fn rootHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r) // Call handler function
		if err == nil {
			return
		}
		// This is where our error handling logic starts.
		log.Printf("[ERROR OCCURRED]: %v", err) // Log the error.

		clientError, ok := err.(errors2.ClientError) // Check if it is a ClientError.
		if !ok {
			// If the error is not ClientError, assume that it is ServerError.
			// return 500 Internal Server Error.
			log.Printf("An internal error accured: %v", err)
			utils.ERROR(w, http.StatusInternalServerError, errors.New("internal server error"), nil)
			return
		}
		log.Printf("[ERROR OCCURRED]: %v", err)

		status := clientError.ResponseStatus() // Get http status code and headers.
		utils.ERROR(w, status, clientError, clientError.ResponseData())
		return
	}
}

// NewRouter returns new instance of router
func NewRouter(services *postgres.Services) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/auth", authRoutes(services.User))
	r.Mount("/api/profile", profileRoutes(services.User, services.Profile))

	return r
}
