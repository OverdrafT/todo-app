package auth

import (
	"context"

	"github.com/silverspase/todo/internal/modules/auth/model"
)

type UseCase interface {
	CreateUser(ctx context.Context, items model.User) (string, error)
	GetAllUsers(ctx context.Context, page int) ([]model.User, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, item model.User) (string, error)
	DeleteUser(ctx context.Context, id string) (string, error)
}
