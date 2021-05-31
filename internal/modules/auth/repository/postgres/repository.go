package postgres

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/silverspase/todo/internal/modules/auth"
	"github.com/silverspase/todo/internal/modules/auth/model"
)

const pageSize = 2

type postgres struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewRepository(conn *gorm.DB, logger *zap.Logger) auth.Repository {
	return postgres{
		conn:   conn,
		logger: logger,
	}
}

func (p postgres) CreateUser(ctx context.Context, entry model.User) (string, error) {
	p.logger.Debug("CreateItem")

	res := p.conn.Create(&entry)
	if res.Error != nil {
		return "", res.Error
	}

	p.logger.Debug("created", zap.Any("item", res))

	return entry.ID, nil
}

func (p postgres) GetAllUsers(ctx context.Context, page int) (entries []model.User, err error) {
	p.logger.Debug("GetAllUsers", zap.Int("page", page))
	if page < 1 {
		page = 1
	}

	res := p.conn.Limit(pageSize).Offset((page - 1) * pageSize).Find(&entries)
	if res.Error != nil {
		return nil, res.Error
	}

	return entries, nil
}

func (p postgres) GetUser(ctx context.Context, id string) (model.User, error) {
	p.logger.Debug("GetUser", zap.String("id", id))

	item := model.User{ID: id}
	err := p.conn.First(&item).Error

	if err != nil {
		return item, err
	}

	return item, nil
}

func (p postgres) UpdateUser(ctx context.Context, newEntry model.User) (string, error) {
	p.logger.Debug("UpdateUser", zap.String("id", newEntry.ID))

	entry, err := p.GetUser(ctx, newEntry.ID)
	if err != nil {
		return "", err
	}

	entry.Name = newEntry.Name
	err = p.conn.Save(&entry).Error
	if err != nil {
		return "", err
	}

	return entry.ID, nil
}

func (p postgres) DeleteUser(ctx context.Context, id string) (string, error) {
	p.logger.Info("DeleteUser", zap.String("id", id))

	err := p.conn.Delete(&model.User{ID: id}).Error
	if err != nil {
		return "", err
	}

	return id, nil
}
