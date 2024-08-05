// Package handlers manages the different versions of the API.
package handlers

import (
	"net/http"
	"os"

	v1 "github.com/iamNilotpal/service/apps/services/sales/handlers/v1"
	"github.com/iamNilotpal/service/business/web/v1/middlewares"
	"github.com/iamNilotpal/service/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Log      *zap.Logger
	Shutdown chan os.Signal
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	app := web.NewApp(cfg.Shutdown, middlewares.Logger(cfg.Log), middlewares.Panics())

	v1.SetupRoutes(app, v1.Config{Log: cfg.Log, Build: cfg.Build})

	return app.Mux
}
