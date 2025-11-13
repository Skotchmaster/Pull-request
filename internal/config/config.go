package config

import (
	"fmt"
	"os"
	"pull_request/internal/logging"

)

type Config struct {
	JWTSecret []byte
	DBDSN     string
	Port      string
}

func mustEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("missing required env %s", name)
	}
	return v, nil
}

func Load(logger logging.Logger) (*Config, error) {
	jwt, err := mustEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	dsn, err := mustEnv("DB_DSN")
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Infof("PORT not set, using default %s", port)
	}

	return &Config{
		JWTSecret: []byte(jwt),
		DBDSN:     dsn,
		Port:      port,
	}, nil
}