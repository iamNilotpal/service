// Package web contains a small web framework extension.
package web

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type App struct {
	Mux         *chi.Mux
	middlewares []Middleware
	shutdown    chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		MaxAge:           300,
		AllowCredentials: true,
		ExposedHeaders:   []string{"Link"},
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodHead,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))

	return &App{
		middlewares: mw,
		Mux:         mux,
		shutdown:    shutdown,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.middlewares, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := Values{Now: time.Now().UTC()}
		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	finalPath := path
	if group != "" {
		finalPath = group + path
	}

	a.Mux.MethodFunc(method, finalPath, h)
}

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateShutdown(err error) bool {
	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):
		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
