package main

import (
	"context"
	"os"
	"os/signal"
	"pull_request/internal/config"
	"pull_request/internal/db"
	"pull_request/internal/logging"
	"syscall"
)

func main() {
	logger := logging.NewStdLogger()

	cfg, err := config.Load(logger)
	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	gdb, err := db.Init(ctx, cfg.DBDSN)
	if err != nil {
		logger.Errorf("failed to init db: %v", err)
		os.Exit(1)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		logger.Errorf("failed to get sql.DB: %v", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
}