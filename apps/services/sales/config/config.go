package config

import (
	"time"

	"github.com/ardanlabs/conf/v3"
)

type web struct {
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:localhost:3000"`
	DebugHost       string        `conf:"default:localhost:4000"`
}

type auth struct {
	Issuer     string `conf:"default:sales_api"`
	KeysFolder string `conf:"default:zarf/keys/"`
	ActiveKID  string `conf:"default:private_key_file"`
}

type db struct {
	MaxIdleConns int    `conf:"default:2"`
	MaxOpenConns int    `conf:"default:0"`
	DisableTLS   bool   `conf:"default:true"`
	Name         string `conf:"default:postgres"`
	User         string `conf:"default:postgres"`
	Host         string `conf:"default:localhost"`
	Password     string `conf:"default:postgres,mask"`
}

type SalesAPIConfig struct {
	DB   db
	Web  web
	Auth auth
	conf.Version
}
