package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/iamNilotpal/service/business/web/auth"
	v1 "github.com/iamNilotpal/service/business/web/v1"
	"github.com/iamNilotpal/service/foundation/web"
)

// Set of error variables for handling user group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Authenticate validates a JWT from the `Authorization` header.
func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := a.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return auth.NewAuthError(http.StatusUnauthorized, "authenticate: failed: %s", err)
			}

			ctx = auth.SetClaims(ctx, claims)
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Authorize validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func Authorize(a *auth.Auth, rule string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := auth.GetClaims(ctx)

			if claims.Subject == "" {
				return auth.NewAuthError(http.StatusUnauthorized, "authorize: you are not authorized for that action, no claims")
			}

			// I will use an zero valued user id if it doesn't exist.
			var userID int
			id := web.Param(r, "user_id")

			if id != "" {
				var err error
				userID, err = strconv.Atoi(id)

				if err != nil {
					return v1.NewRequestError(ErrInvalidID.Error(), http.StatusBadRequest)
				}
				ctx = auth.SetUserID(ctx, userID)
			}

			if err := a.Authorize(ctx, claims, userID, rule); err != nil {
				return auth.NewAuthError(http.StatusForbidden, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
