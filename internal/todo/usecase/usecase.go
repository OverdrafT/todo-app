package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/todo"
	"github.com/silverspase/todo/internal/todo/model"
)

type itemUseCase struct {
	repo   todo.Repository
	logger *zap.Logger
}

func NewItemUseCase(logger *zap.Logger, repo todo.Repository) todo.UseCase {
	return &itemUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (i itemUseCase) CreateItem(ctx context.Context, item model.Item) (string, error) {
	return i.repo.CreateItem(ctx, item)
}

func (i itemUseCase) GetAllItems(ctx context.Context, page int) ([]model.Item, error) {
	return i.repo.GetAllItems(ctx, page)
}

func (i itemUseCase) GetItem(ctx context.Context, id string) (model.Item, error) {
	return i.repo.GetItem(ctx, id)
}

func (i itemUseCase) UpdateItem(ctx context.Context, item model.Item) (string, error) {
	return i.repo.UpdateItem(ctx, item)
}

func (i itemUseCase) DeleteItem(ctx context.Context, id string) (string, error) {
	return i.repo.DeleteItem(ctx, id)
}
