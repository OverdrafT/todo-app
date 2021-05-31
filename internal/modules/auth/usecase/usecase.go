package usecase

import (
	"context"

	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/modules/auth"
	"github.com/silverspase/todo/internal/modules/auth/model"
)

type useCase struct {
	repo   auth.Repository
	logger *zap.Logger
}

func NewUseCase(logger *zap.Logger, repo auth.Repository) auth.UseCase {
	return &useCase{
		repo:   repo,
		logger: logger,
	}
}

func (u useCase) CreateUser(ctx context.Context, entry model.User) (string, error) {
	return u.repo.CreateUser(ctx, entry)
}

func (u useCase) GetAllUsers(ctx context.Context, page int) ([]model.User, error) {
	return u.repo.GetAllUsers(ctx, page)
}

func (u useCase) GetUser(ctx context.Context, id string) (model.User, error) {
	return u.repo.GetUser(ctx, id)
}

func (u useCase) UpdateUser(ctx context.Context, entry model.User) (string, error) {
	return u.repo.UpdateUser(ctx, entry)
}

func (u useCase) DeleteUser(ctx context.Context, id string) (string, error) {
	return u.repo.DeleteUser(ctx, id)
}
