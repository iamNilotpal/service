package users_handler

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/iamNilotpal/service/business/web/v1"
)

type handlers struct{}

func NewHandler() *handlers {
	return &handlers{}
}

func (h *handlers) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("UNTRUSTED ERROR")
	}

	return v1.NewRequestError(errors.New("bad request"), http.StatusBadRequest)
	// return web.Respond(ctx, w, map[string]any{"ok": true, "statusCode": http.StatusOK}, http.StatusOK)
}
