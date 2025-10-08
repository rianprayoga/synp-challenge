package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	appError "inventories-app/internal/error"
	"inventories-app/internal/model"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = 3 * time.Second

func (r *PostgresDBRepo) GetItemsWithCursor(pageSize int, cursor time.Time) ([]*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
			SELECT id, name, stock, created_at, updated_at
			FROM items
			WHERE created_at < $1
			ORDER BY created_at DESC
			LIMIT $2
		`

	rows, err := r.DB.QueryContext(ctx, query, cursor, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.Item

	for rows.Next() {
		var tmp model.Item
		err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Stock,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, &tmp)
	}

	return items, nil
}

func (r *PostgresDBRepo) GetItems(pageSize int) ([]*model.Item, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
			SELECT id, name, stock, created_at, updated_at
			FROM items
			ORDER BY created_at DESC
			LIMIT $1
		`

	rows, err := r.DB.QueryContext(ctx, query, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.Item

	for rows.Next() {
		var tmp model.Item
		err := rows.Scan(
			&tmp.ID,
			&tmp.Name,
			&tmp.Stock,
			&tmp.CreatedAt,
			&tmp.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, &tmp)

	}

	return items, nil
}

func (r *PostgresDBRepo) GetItem(id string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
			SELECT id, name, stock, created_at, updated_at
			FROM items
			WHERE id = $1
		`

	row := r.DB.QueryRowContext(ctx, query, id)
	var item model.Item
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Stock,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrItemNotFound
		}
		return nil, err
	}

	return &item, nil

}
