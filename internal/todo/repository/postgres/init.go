package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/config"
	"github.com/silverspase/todo/internal/todo"
)

const (
	HOST = "localhost"
	PORT = 5432
)

type postgres struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewPostgres(logger *zap.Logger, cfg config.Config) (todo.Repository, error) {
	if cfg.Username == "" {
		return nil, errors.New("env var POSTGRES_USER not set")
	}

	if cfg.Password == "" {
		return nil, errors.New("env var POSTGRES_PASSWORD not set")
	}

	if cfg.DB == "" {
		return nil, errors.New("env var POSTGRES_DB not set")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, cfg.Username, cfg.Password, cfg.DB)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to exec sql.Open function")
		return nil, err
	}

	db := postgres{
		conn:   conn,
		logger: logger,
	}

	err = db.conn.Ping()
	if err != nil {
		logger.Error("failed to ping db")
		return db, err
	}

	log.Println("database connection established")

	return db, nil
}
