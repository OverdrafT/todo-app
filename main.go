package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/silverspase/todo/internal"
	"github.com/silverspase/todo/internal/config"
	appLogger "github.com/silverspase/todo/internal/logger"
	"github.com/silverspase/todo/internal/metadata/version"
	"github.com/silverspase/todo/internal/todo"
	"github.com/silverspase/todo/internal/todo/repository/memory"
	"github.com/silverspase/todo/internal/todo/repository/postgres"
	gorilla "github.com/silverspase/todo/internal/todo/transport/gorilla-mux"
	"github.com/silverspase/todo/internal/todo/usecase"
)

func main() {
	cfg := config.Init()
	logger := appLogger.Init(cfg)
	logger.Info("Starting server", zap.String("params:",
		fmt.Sprintf("port: %s, log level: %s, repo: %s", cfg.Port, cfg.LogLevel, cfg.Repository)))

	var err error
	var repo todo.Repository
	switch cfg.Repository {
	case config.MemoryRepo:
		repo = memory.NewMemoryStorage(logger)
	case config.PostgresRepo:
		repo, err = postgres.NewPostgres(logger, cfg)
		if err != nil {
			logger.Fatal("Postgres Init failed", zap.Error(err))
		}
	default:
		logger.Fatal("unable to define repo type")
	}

	useCase := usecase.NewItemUseCase(logger, repo)
	transport := gorilla.NewTransport(logger, useCase)
	server := internal.NewServer(logger, transport)
	router := server.GorillaMuxRouter(version.BuildTime, version.Commit, version.Release)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown:
	shutdown := make(chan struct{}, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			shutdown <- struct{}{}
			log.Printf("%v", err)
		}
	}()
	logger.Info("The service is ready to listen and serve", zap.String("port:", cfg.Port))

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			logger.Warn("Got SIGINT...")
		case syscall.SIGTERM:
			logger.Warn("Got SIGTERM...")
		}
	case <-shutdown:
		logger.Error("Got an error...")
	}

	logger.Info("The service is shutting down...")
	srv.Shutdown(context.Background())
	logger.Info("Done")
}
