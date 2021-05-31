package memory

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/modules/auth"
	"github.com/silverspase/todo/internal/modules/auth/model"
)

type memoryStorage struct {
	users map[string]model.User // TODO change to sync.Map
	// usersArray []model.User // TODO use this for pagination in GetAllUsers
	logger *zap.Logger
}

func NewMemoryStorage(logger *zap.Logger) auth.Repository {
	return &memoryStorage{
		users:  make(map[string]model.User),
		logger: logger,
	}
}

func (m memoryStorage) CreateUser(ctx context.Context, entry model.User) (string, error) {
	m.logger.Debug("CreateUser")
	m.users[entry.ID] = entry

	return entry.ID, nil
}

func (m memoryStorage) GetAllUsers(ctx context.Context, page int) (res []model.User, err error) {
	for _, user := range m.users {
		res = append(res, user)
	}

	return res, nil
}

func (m memoryStorage) GetUser(ctx context.Context, id string) (model.User, error) {
	m.logger.Debug("GetItem")
	item, ok := m.users[id]
	if !ok {
		return item, errors.New("not found")
	}

	return item, nil
}

func (m memoryStorage) UpdateUser(ctx context.Context, item model.User) (string, error) {
	_, ok := m.users[item.ID]
	if !ok {
		return "", errors.New("entry with given id not found, nothing to update")
	}

	m.users[item.ID] = item

	return item.ID, nil
}

func (m memoryStorage) DeleteUser(ctx context.Context, id string) (string, error) {
	_, ok := m.users[id]
	if !ok {
		return "", errors.New("item with given id not found, nothing to delete")
	}

	delete(m.users, id)

	return id, nil
}
