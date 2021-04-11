package router

import (
	"github.com/go-chi/chi"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/controllers"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/middlewares"
	"github.com/usmanzaheer1995/devconnect-go-v2/internal/models/postgres/user"
	"net/http"
)

func authRoutes(us user.UserService) http.Handler {
	r := chi.NewRouter()

	uc := controllers.NewUserController(us)
	r.
		With(middlewares.AuthJwtVerify).
		Get("/", errorHandler(uc.FindByID))

	r.Post("/register", errorHandler(uc.Create))

	r.Post("/login", errorHandler(uc.Login))

	return r
}
