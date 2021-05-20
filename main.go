package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/silverspase/k8s-prod-service/internal"
	"github.com/silverspase/k8s-prod-service/internal/metadata/version"
	"github.com/silverspase/k8s-prod-service/internal/todo/repository/memory"
	gorilla "github.com/silverspase/k8s-prod-service/internal/todo/transport/gorilla-mux"
	"github.com/silverspase/k8s-prod-service/internal/todo/usecase"
	"go.uber.org/zap"
)

const defaultPort = "8000"

func main() {
	logger, _ := zap.NewProduction()

	log.Printf("Starting the service...\n commit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repo := memory.NewMemoryStorage(logger)
	useCase := usecase.NewItemUseCase(logger, repo)
	transport := gorilla.NewTransport(logger, useCase)
	server := internal.NewServer(logger, transport)
	router := server.GorillaMuxRouter(version.BuildTime, version.Commit, version.Release)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// graceful shutdown:
	shutdown := make(chan struct{}, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			shutdown <- struct{}{}
			log.Printf("%v", err)
		}
	}()
	log.Print("The service is ready to listen and serve.")

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Print("Got SIGINT...")
		case syscall.SIGTERM:
			log.Print("Got SIGTERM...")
		}
	case <-shutdown:
		log.Printf("Got an error...")
	}

	log.Print("The service is shutting down...")
	srv.Shutdown(context.Background())
	log.Print("Done")
}
