package postgres

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/silverspase/todo/internal/modules/todo"
	"github.com/silverspase/todo/internal/modules/todo/model"
)

const pageSize = 2

type postgres struct {
	conn   *gorm.DB
	logger *zap.Logger
}

func NewRepository(conn *gorm.DB, logger *zap.Logger) todo.Repository {
	return postgres{
		conn:   conn,
		logger: logger,
	}
}

func (p postgres) CreateItem(ctx context.Context, item model.Item) (string, error) {
	p.logger.Debug("CreateItem")

	res := p.conn.Create(&item)
	if res.Error != nil {
		return "", res.Error
	}

	p.logger.Debug("created", zap.Any("item", res))

	return item.ID, nil
}

func (p postgres) GetAllItems(ctx context.Context, page int) (items []model.Item, err error) {
	p.logger.Debug("GetAllItems", zap.Int("page", page))
	if page < 1 {
		page = 1
	}

	res := p.conn.Limit(pageSize).Offset((page - 1) * pageSize).Find(&items)
	if res.Error != nil {
		return nil, res.Error
	}

	return items, nil
}

func (p postgres) GetItem(ctx context.Context, id string) (model.Item, error) {
	p.logger.Debug("GetItem", zap.String("id", id))

	item := model.Item{ID: id}
	err := p.conn.First(&item).Error

	if err != nil {
		return item, err
	}

	return item, nil
}

func (p postgres) UpdateItem(ctx context.Context, newItem model.Item) (string, error) {
	p.logger.Debug("UpdateItem", zap.String("id", newItem.ID))

	item, err := p.GetItem(ctx, newItem.ID)
	if err != nil {
		return "", err
	}

	item.Title = newItem.Title
	err = p.conn.Save(&item).Error
	if err != nil {
		return "", err
	}

	return item.ID, nil
}

func (p postgres) DeleteItem(ctx context.Context, id string) (string, error) {
	p.logger.Info("DeleteItem", zap.String("id", id))

	err := p.conn.Delete(&model.Item{ID: id}).Error
	if err != nil {
		return "", err
	}

	return id, nil
}
