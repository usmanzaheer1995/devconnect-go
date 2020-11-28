package router

import (
	"github.com/go-chi/chi"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/controllers"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/middlewares"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/user"
	"net/http"
)

func userRoutes(us user.UserService) http.Handler {
	r := chi.NewRouter()

	uc := controllers.NewUserController(us)
	r.Get("/", uc.Find)

	r.
		With(middlewares.MyMiddleware, middlewares.SecondMiddleware).
		Post("/register", uc.Create)

	r.Post("/login", uc.Login)

	return r
}
