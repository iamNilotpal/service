package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/ardanlabs/conf/v3"
	"github.com/iamNilotpal/service/apps/services/sales/config"
	"github.com/iamNilotpal/service/foundation/logger"
	"go.uber.org/zap"
)

var build = "develop"

const description string = ""

func main() {
	log := logger.New("SALES-API")

	if err := run(log); err != nil {
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

	return nil
}
