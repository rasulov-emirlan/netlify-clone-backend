package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type (
	minio struct {
		URL             string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
	}
	server struct {
		Port         string
		Domain       string
		TimeoutRead  time.Duration
		TimeoutWrite time.Duration
	}
	Config struct {
		Server   server
		Database string
		Minio    minio
	}
)

const (
	serverPort         = "SERVER_PORT"
	serverDomain       = "SERVER_DOMAIN"
	serverTimeoutRead  = "SERVER_TIMEOUT_READ"
	serverTimeoutWrite = "SERVER_TIMEOUT_WRITE"

	databaseURL = "DATABASE_URL"

	minioURL          = "MINIO_URL"
	minioAcceessKeyID = "MINIO_ACCESS_KEY_ID"
	minioSecretKey    = "MINIO_SECRET_KEY"
	minioUseSSL       = "MINIO_USE_SSL"
)

var (
	ErrNoServerData   = errors.New("config: did not find configs for server")
	ErrNoDatabaseData = errors.New("config: did not find configs for database")
	ErrNoMinioData    = errors.New("config: did not find configs for minio")
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
	b, err := strconv.ParseBool(os.Getenv(minioUseSSL))
	if err != nil {
		return nil, err
	}
	cfg := Config{
		Server: server{
			Port:         os.Getenv(serverPort),
			Domain:       os.Getenv(serverDomain),
			TimeoutRead:  tR,
			TimeoutWrite: tW,
		},
		Database: os.Getenv(databaseURL),
		Minio: minio{
			URL:             os.Getenv(minioURL),
			AccessKeyID:     os.Getenv(minioAcceessKeyID),
			SecretAccessKey: os.Getenv(minioSecretKey),
			UseSSL:          b,
		},
	}
	if cfg.Server.Port == "" || cfg.Server.Domain == "" {
		return nil, ErrNoServerData
	}
	cfg.Server.Port = ":" + cfg.Server.Port
	if cfg.Database == "" {
		return nil, ErrNoDatabaseData
	}
	if cfg.Minio.AccessKeyID == "" || cfg.Minio.SecretAccessKey == "" || cfg.Minio.URL == "" {
		return nil, ErrNoMinioData
	}
	return &cfg, nil
}
