package middlewares

import (
	"context"
	"net/http"

	"github.com/iamNilotpal/service/business/web/auth"
	v1 "github.com/iamNilotpal/service/business/web/v1"
	"github.com/iamNilotpal/service/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error("ERROR", zap.Any("message", err))

				var status int
				var er v1.ErrorResponse

				switch {
				case v1.IsRequestError(err):
					reqErr := v1.GetRequestError(err)
					er = v1.ErrorResponse{Message: reqErr.Error()}
					status = reqErr.Status

				case auth.IsAuthError(err):
					authErr := auth.GetAuthError(err)
					er = v1.ErrorResponse{
						Message: http.StatusText(authErr.Code),
					}
					status = authErr.Code

				default:
					er = v1.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, web.NewErrorResponse(er), status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
