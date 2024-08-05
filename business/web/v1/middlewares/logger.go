package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/iamNilotpal/service/foundation/web"
	"go.uber.org/zap"
)

// Logger writes information about the request to the logs.
func Logger(log *zap.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(
				"request started",
				zap.String("path", path),
				zap.String("method", r.Method),
				zap.String("remoteAddress", r.RemoteAddr),
			)

			err := handler(ctx, w, r)

			log.Info(
				"request success",
				zap.String("path", path),
				zap.String("method", r.Method),
				zap.String("remoteAddress", r.RemoteAddr),
				zap.Int("statusCode", v.StatusCode), zap.Any("responseTime", time.Since(v.Now)),
			)

			return err
		}

		return h
	}

	return m
}
