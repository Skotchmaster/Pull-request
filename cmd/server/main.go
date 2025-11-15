package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pull_request/internal/config"
	appdb "pull_request/internal/db"
	"pull_request/internal/handlers"
	httproutes "pull_request/internal/http"
	"pull_request/internal/logging"
	"pull_request/internal/repo"
	"pull_request/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	logger := logging.NewStdLogger()

	cfg, err := config.Load(logger)
	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	gdb, err := appdb.Init(ctx, cfg.DBURL)
	if err != nil {
		logger.Errorf("failed to init db: %v", err)
		os.Exit(1)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		logger.Errorf("failed to get sql.DB: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("failed to close db: %v", err)
		}
	}()

	prRepo := &repo.PullRequestRepo{DB: gdb}
	teamRepo := &repo.TeamRepo{DB: gdb}
	userRepo := &repo.UserRepo{DB: gdb}

	prService := &service.PullRequestService{Repo: prRepo}
	teamService := &service.TeamService{Repo: teamRepo}
	userService := &service.UserService{Repo: userRepo}

	prHandler := handlers.PullRequestHandler{Service: prService, Logger: logger}
	teamHandler := handlers.TeamHandler{Service: teamService, Logger: logger}
	userHandler := handlers.UserHandler{Service: userService, Logger: logger}

	allHandlers := handlers.Handlers{
		PRHandler:   prHandler,
		TeamHandler: teamHandler,
		UserHandler: userHandler,
	}

	e := echo.New()

	httproutes.Register(e, &httproutes.Deps{
		Handlers:  allHandlers,
	})

	serverErr := make(chan error, 1)

	go func() {
		addr := ":" + cfg.Port
		logger.Infof("starting HTTP server on %s", addr)
		if err := e.Start(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Warnf("shutdown signal received, shutting down...")
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("http server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("http server forced to shutdown: %v", err)
	} else {
		logger.Infof("http server stopped gracefully")
	}
}
