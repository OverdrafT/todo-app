package todo

import (
	"context"

	"github.com/silverspase/k8s-prod-service/internal/todo/model"
)

type UseCase interface {
	CreateItem(ctx context.Context, items model.Item) (string, error)
	GetAllItems(ctx context.Context) ([]model.Item, error)
	GetItem(ctx context.Context, id string) (model.Item, bool)
	UpdateItem(ctx context.Context, item model.Item) (string, error)
	DeleteItem(ctx context.Context, id string) (string, error)
}
