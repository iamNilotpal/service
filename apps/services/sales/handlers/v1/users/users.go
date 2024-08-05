package users_handler

import (
	"context"
	"net/http"

	"github.com/iamNilotpal/service/foundation/web"
)

type handlers struct{}

func NewHandler() *handlers {
	return &handlers{}
}

func (h *handlers) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, map[string]any{"ok": true, "statusCode": http.StatusOK}, http.StatusOK)
}
