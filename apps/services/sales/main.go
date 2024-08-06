package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf/v3"
	"github.com/iamNilotpal/service/apps/services/sales/config"
	"github.com/iamNilotpal/service/apps/services/sales/handlers"
	"github.com/iamNilotpal/service/business/web/auth"
	"github.com/iamNilotpal/service/foundation/keystore"
	"github.com/iamNilotpal/service/foundation/logger"
	"go.uber.org/zap"
)

var build string = "develop"

const description string = ""

func main() {
	log := logger.New("SALES-API")
	defer log.Sync()

	if err := run(log); err != nil {
		log.Error("startup", zap.Any("ERROR", err))
		log.Sync()
		os.Exit(1)
	}
}

func run(logger *zap.Logger) error {
	// GOMAXPROCS
	logger.Info("startup", zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))

	// Configuration
	cfg := config.SalesAPIConfig{
		Version: conf.Version{
			Build: build,
			Desc:  description,
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("error parsing Config: %w", err)
	}

	// App Starting
	logger.Info("starting service", zap.String("version", build))
	defer logger.Info("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}

	logger.Info("startup", zap.String("config", out))

	// Initialize authentication support

	logger.Info("startup", zap.String("status", "initializing authentication support"))

	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{KeyLookup: ks, Log: logger}
	auth, err := auth.New(authCfg)

	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// Start API Service
	logger.Info("startup", zap.String("status", "initializing V1 API support"))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{Build: build, Shutdown: shutdown, Log: logger, Auth: auth},
	)

	api := http.Server{
		Handler:      apiMux,
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		ErrorLog:     zap.NewStdLog(logger),
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("startup", zap.String("status", "api router started"), zap.String("host", api.Addr))
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown Server
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("shutdown", zap.String("status", "shutdown started"), zap.Any("signal", sig))
		defer logger.Info("shutdown", zap.String("status", "shutdown complete"), zap.Any("signal", sig))

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
