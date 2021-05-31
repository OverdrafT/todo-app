package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	application "github.com/silverspase/todo/internal/app"
)

func main() {
	app, err := application.Init()
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown:
	shutdown := make(chan struct{}, 1)
	go func() {
		err := app.Srv.ListenAndServe()
		if err != nil {
			shutdown <- struct{}{}
			log.Printf("%v", err)
		}
	}()
	app.Logger.Info("The service is ready to listen and serve", zap.String("port:", app.Cfg.Port))

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			app.Logger.Warn("Got SIGINT...")
		case syscall.SIGTERM:
			app.Logger.Warn("Got SIGTERM...")
		}
	case <-shutdown:
		app.Logger.Error("Got an error...")
	}

	app.Logger.Info("The service is shutting down...")
	app.Srv.Shutdown(context.Background())
	app.Logger.Info("Done")
}
