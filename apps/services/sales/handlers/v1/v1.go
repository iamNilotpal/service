package v1

import (
	"net/http"

	users_handler "github.com/iamNilotpal/service/apps/services/sales/handlers/v1/users"
	"github.com/iamNilotpal/service/business/web/auth"
	"github.com/iamNilotpal/service/business/web/v1/middlewares"
	"github.com/iamNilotpal/service/foundation/web"
	"go.uber.org/zap"
)

const version string = "/api/v1/users"

type Config struct {
	Build string
	Auth  *auth.Auth
	Log   *zap.Logger
}

// Routes binds all the version 1 routes.
func SetupRoutes(app *web.App, cfg Config) {
	usersHandler := users_handler.NewHandler()

	app.Handle(
		http.MethodGet, version, "",
		usersHandler.GetUsers, middlewares.Authorize(cfg.Auth, auth.RuleAdminOnly),
	)
	app.Handle(http.MethodPost, version, "", usersHandler.GetUsers)
}
