package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type (
	server struct {
		Port         string
		Domain       string
		TimeoutRead  time.Duration
		TimeoutWrite time.Duration
	}
	Config struct {
		Server   server
		Database string
	}
)

const (
	serverPort         = "SERVER_PORT"
	serverDomain       = "SERVER_DOMAIN"
	serverTimeoutRead  = "SERVER_TIMEOUT_READ"
	serverTimeoutWrite = "SERVER_TIMEOUT_WRITE"

	databaseURL = "DATABASE_URL"
)

var (
	ErrNoServerData   = errors.New("config: did not find configs for server")
	ErrNoDatabaseData = errors.New("config: did not find configs for database")
)

func NewConfig(filenames ...string) (*Config, error) {
	if len(filenames) > 0 {
		if err := godotenv.Load(filenames...); err != nil {
			return nil, err
		}
	}
	tR, err := time.ParseDuration(os.Getenv(serverTimeoutRead))
	if err != nil {
		return nil, ErrNoServerData
	}
	tW, err := time.ParseDuration(os.Getenv(serverTimeoutWrite))
	if err != nil {
		return nil, ErrNoServerData
	}
	cfg := Config{
		Server: server{
			Port:         os.Getenv(serverPort),
			Domain:       os.Getenv(serverDomain),
			TimeoutRead:  tR,
			TimeoutWrite: tW,
		},
		Database: os.Getenv(databaseURL),
	}
	if cfg.Server.Port == "" || cfg.Server.Domain == "" {
		return nil, ErrNoServerData
	}
	if cfg.Database == "" {
		return nil, ErrNoDatabaseData
	}
	return &cfg, nil
}
