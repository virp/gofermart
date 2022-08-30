package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/virp/gofermart/config"
	"github.com/virp/gofermart/internal/handlers"
	"github.com/virp/gofermart/migrations"
	"github.com/virp/gofermart/pkg/accrual"
	"github.com/virp/gofermart/pkg/logger"
	"github.com/virp/gofermart/pkg/postgres"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

const (
	defaultReadTimeout             = 5 * time.Second
	defaultWriteTimeout            = 10 * time.Second
	defaultIdleTimeout             = 120 * time.Second
	defaultShutdownTimeout         = 3 * time.Second
	defaultAppSecret               = "supertsecretkey"
	defaultAppUserAuthCookieName   = "user"
	defaultOrderStatusWorkersCount = 10
)

func main() {

	// Construct the application logger.
	log, err := logger.New("gofermart")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err = run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// ========================================================================
	// GOMAXPROCS

	// Want to see what maxprocs reports.
	opt := maxprocs.Logger(log.Infof)

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ========================================================================
	// Configuration

	cfg, err := config.NewGofermartConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// ========================================================================
	// Database

	log.Infow("startup", "status", "initializing database")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := postgres.New(ctx, cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database")
		db.Close()
	}()
	if err = migrations.Migrate(ctx, db); err != nil {
		return fmt.Errorf("database migrations: %w", err)
	}

	// ========================================================================
	// HTTP Server

	log.Infow("startup", "status", "initializing http server")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	handler := handlers.New(handlers.Config{
		Shutdown:                shutdown,
		DB:                      db,
		Log:                     log,
		AppSecret:               defaultAppSecret,
		AppUserAuthCookieName:   defaultAppUserAuthCookieName,
		AccrualSystem:           accrual.New(cfg.AccrualSystemAddress),
		OrderStatusWorkersCount: defaultOrderStatusWorkersCount,
	})

	api := http.Server{
		Addr:         cfg.RunAddress,
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "http server started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// ========================================================================
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
