package sql

import (
	"errors"
	"fmt"
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/silverspase/todo/internal/config"
	"github.com/silverspase/todo/internal/modules/todo/model"
)

const (
	HOST = "localhost"
	PORT = 5432
)

func NewConn(logger *zap.Logger, cfg config.Config) (*gorm.DB, error) {
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

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to exec gorm.Open function")
		return nil, err
	}

	logger.Info("Migrating", zap.String("model", "Item"))

	err = conn.AutoMigrate(&model.Item{})
	if err != nil {
		logger.Error("Error during migrating Item struct", zap.Error(err))
		return nil, err
	}

	log.Println("database connection established")

	return conn, nil
}
