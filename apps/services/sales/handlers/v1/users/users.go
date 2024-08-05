package users_handler

import (
	"context"
	"math/rand"
	"net/http"

	v1 "github.com/iamNilotpal/service/business/web/v1"
	"github.com/iamNilotpal/service/foundation/web"
)

type handlers struct{}

func NewHandler() *handlers {
	return &handlers{}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (h *handlers) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return v1.NewRequestError(
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
		)
	}

	return web.Respond(
		ctx, w, web.NewSuccessResponse(User{ID: 1, Username: "iamNilotpal"}), http.StatusOK,
	)
}
