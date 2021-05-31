package app

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/silverspase/todo/internal/app/repository/sql"
	"github.com/silverspase/todo/internal/config"
	appLogger "github.com/silverspase/todo/internal/logger"
	"github.com/silverspase/todo/internal/modules/auth"
	authMemory "github.com/silverspase/todo/internal/modules/auth/repository/memory"
	authRepo "github.com/silverspase/todo/internal/modules/auth/repository/postgres"
	authTransport "github.com/silverspase/todo/internal/modules/auth/transport/gorilla-mux"
	"github.com/silverspase/todo/internal/modules/todo"
	"github.com/silverspase/todo/internal/modules/todo/repository/memory"
	"github.com/silverspase/todo/internal/modules/todo/repository/postgres"
	todoTransport "github.com/silverspase/todo/internal/modules/todo/transport/gorilla-mux"

	authUseCase "github.com/silverspase/todo/internal/modules/auth/usecase"
	todoUseCase "github.com/silverspase/todo/internal/modules/todo/usecase"
)

type App struct {
	Todo todo.Transport
	Auth auth.Transport

	Srv     *http.Server
	Logger  *zap.Logger
	Cfg     config.Config
	isReady *atomic.Value
}

func Init() (*App, error) {
	cfg := config.Init()

	logger := appLogger.Init(cfg)
	logger.Info("Starting server", zap.String("params:",
		fmt.Sprintf("port: %s, log level: %s, repo: %s", cfg.Port, cfg.LogLevel, cfg.Repository)))

	sqlConn, err := sql.NewConn(logger, cfg)
	if err != nil {
		return nil, err
	}

	application := &App{
		Todo:   initTodoModule(cfg, logger, sqlConn),
		Auth:   initAuthModule(cfg, logger, sqlConn),
		Logger: logger,
		Cfg:    cfg,
	}

	application.Srv = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: gorillaMuxRouter(application),
	}

	return application, nil
}

func initTodoModule(cfg config.Config, logger *zap.Logger, sqlConn *gorm.DB) todo.Transport {
	var repo todo.Repository
	switch cfg.Repository {
	case config.MemoryRepo:
		repo = memory.NewMemoryStorage(logger)
	case config.PostgresRepo:
		repo = postgres.NewRepository(sqlConn, logger)
	default:
		logger.Fatal("unable to define repo type")
	}

	useCase := todoUseCase.NewItemUseCase(logger, repo)

	return todoTransport.NewTransport(logger, useCase) // add support of several transports
}

func initAuthModule(cfg config.Config, logger *zap.Logger, sqlConn *gorm.DB) auth.Transport {
	var repo auth.Repository
	switch cfg.Repository {
	case config.MemoryRepo:
		repo = authMemory.NewMemoryStorage(logger)
	case config.PostgresRepo:
		repo = authRepo.NewRepository(sqlConn, logger)
	default:
		logger.Fatal("unable to define repo type")
	}

	useCase := authUseCase.NewUseCase(logger, repo)

	return authTransport.NewTransport(logger, useCase) // add support of several transports
}
