package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/silverspase/todo/internal/todo/model"
)

const pageSize = 2

func (p postgres) CreateItem(ctx context.Context, item model.Item) (string, error) {
	p.logger.Debug("CreateItem")
	var id int
	query := `INSERT INTO todo (title, is_deleted)  VALUES ($1, $2) RETURNING id;`
	err := p.conn.QueryRow(query, item.Title, false).Scan(&id)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

// TODO handle null result from SQL (offset more then rows)
func (p postgres) GetAllItems(ctx context.Context, page int) (items []model.Item, err error) {
	p.logger.Debug("GetAllItems", zap.Int("page", page))
	if page < 1 {
		page = 1
	}

	query := `SELECT * FROM todo WHERE is_deleted = false LIMIT $1 OFFSET $2;`
	rows, err := p.conn.Query(query, pageSize, (page-1)*pageSize)
	if err != nil {
		if err == sql.ErrNoRows {
			return items, errors.New("not found")
		} else {
			return items, err
		}
	}
	defer rows.Close()

	var item model.Item
	for rows.Next() {
		item = model.Item{}
		err = rows.Scan(&item.ID, &item.Title, &item.IsDeleted, &item.CreatedAt)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	err = rows.Err() // get any error encountered ing iteration
	if err != nil {
		return items, err
	}

	return items, nil
}

func (p postgres) GetItem(ctx context.Context, id string) (model.Item, error) {
	p.logger.Debug("GetItem", zap.String("id", id))

	var item model.Item
	query := `SELECT * FROM todo WHERE id = $1 AND is_deleted = false;`
	row := p.conn.QueryRow(query, id)

	err := row.Scan(&item.ID, &item.Title, &item.IsDeleted, &item.CreatedAt)

	switch err {
	case nil:
		return item, nil
	case sql.ErrNoRows:
		return item, errors.New("not found")
	default:
		return item, err
	}
}

func (p postgres) UpdateItem(ctx context.Context, item model.Item) (string, error) {
	p.logger.Debug("UpdateItem", zap.String("id", item.ID))

	query := `UPDATE todo SET title=$1 WHERE id = $2;` // TODO ignore updating deleted items
	_, err := p.conn.Exec(query, item.Title, item.ID)
	if err != nil {
		p.logger.Error("Error during Item updating", zap.Error(err))
		return "", err
	}

	return item.ID, nil
}

func (p postgres) DeleteItem(ctx context.Context, id string) (string, error) {
	p.logger.Info("DeleteItem", zap.String("id", id))

	query := `UPDATE todo SET is_deleted=true WHERE id = $1;`
	_, err := p.conn.Exec(query, id)

	if err != nil {
		p.logger.Error("Error during Item deletion", zap.Error(err))
		return "", err
	}

	return id, nil
}
