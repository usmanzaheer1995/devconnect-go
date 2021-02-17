package router

import (
	"github.com/go-chi/chi"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/controllers"
	"github.com/usmanzaheer1995/devconnect-go-v2/cmd/web/middlewares"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/profile"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/models/postgres/user"
	"net/http"
)

func profileRoutes(us user.UserService, ps profile.ProfileService) http.Handler {
	r := chi.NewRouter()

	pc := controllers.NewProfileController(us, ps)
	r.
		With(middlewares.AuthJwtVerify).
		Post("/", errorHandler(pc.Create))
	r.
		With(middlewares.AuthJwtVerify).
		Put("/", errorHandler(pc.Update))

	return r
}
