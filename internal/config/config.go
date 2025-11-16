package config

import (
	"fmt"
	"os"

	"pull_request/internal/logging"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
	Port   string
}

func mustEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("missing required env %s", name)
	}
	return v, nil
}

func Load(logger logging.Logger) (*Config, error) {
	user, err := mustEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	pass, err := mustEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	name, err := mustEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
		logger.Infof("DB_HOST not set, using default %s", host)
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
		logger.Infof("DB_PORT not set, using default %s", port)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
		logger.Infof("PORT not set, using default %s", httpPort)
	}

	return &Config{
		DBUser: user,
		DBPass: pass,
		DBHost: host,
		DBPort: port,
		DBName: name,
		Port:   httpPort,
	}, nil
}

func (c *Config) DBURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
