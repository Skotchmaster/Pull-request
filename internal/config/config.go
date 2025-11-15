package config

import (
	"fmt"
	"os"
	"pull_request/internal/logging"

)

type Config struct {
	DBURL     string
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
	dsn, err := mustEnv("DB_URL")
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Infof("PORT not set, using default %s", port)
	}

	return &Config{
		DBURL:     dsn,
		Port:      port,
	}, nil
}