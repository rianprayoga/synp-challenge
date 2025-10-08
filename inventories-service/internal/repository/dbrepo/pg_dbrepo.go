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

func (r *PostgresDBRepo) AddItem(item model.CreateItem) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var createdItem model.Item

	err := r.DB.QueryRowContext(
		ctx,
		`INSERT INTO items
			(name, stock, created_at, updated_at)
		VALUES ($1,$2,$3,$4)
		RETURNING id, name, stock, created_at, updated_at`,
		item.Name,
		item.Stock,
		time.Now(),
		time.Now(),
	).Scan(
		&createdItem.ID,
		&createdItem.Name,
		&createdItem.Stock,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdItem, nil
}

func (r *PostgresDBRepo) DeleteItem(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := r.DB.ExecContext(
		ctx,
		`
		DELETE FROM items WHERE id = $1
		`,
		id,
	)

	if err != nil {
		return err
	}
	return nil
}
