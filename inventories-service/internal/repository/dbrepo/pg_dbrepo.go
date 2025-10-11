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

func (r *PostgresDBRepo) ReleaseStock(itemId string, orderId string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var trxQty int
	err = tx.QueryRowContext(ctx,
		`
		UPDATE transactions
		SET
			updated_at = $1,
			released = true
		WHERE item_id = $2 AND order_id = $3 AND released = false
		RETURNING qty
	`,
		time.Now().UTC(),
		itemId,
		orderId,
	).Scan(&trxQty)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrTrxNotFound
		}

		return nil, err
	}

	var updatedItem model.Item
	err = tx.QueryRowContext(
		ctx,
		`UPDATE items SET
			stock = stock + $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 
		RETURNING id, name, stock, version ,created_at, updated_at`,
		trxQty,
		time.Now().UTC(),
		itemId,
	).Scan(
		&updatedItem.ID,
		&updatedItem.Name,
		&updatedItem.Stock,
		&updatedItem.Version,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrItemNotFound
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &updatedItem, nil

}

func (r *PostgresDBRepo) ReduceStock(itemId string, qty int, version int, orderId string) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var updatedItem model.Item

	err = tx.QueryRowContext(
		ctx,
		`UPDATE items SET
			stock = stock - $1,
			updated_at = $2,
			version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, name, stock, version ,created_at, updated_at`,
		qty,
		time.Now().UTC(),
		itemId,
		version,
	).Scan(
		&updatedItem.ID,
		&updatedItem.Name,
		&updatedItem.Stock,
		&updatedItem.Version,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrItemConflictVersion
		}
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO 
			transactions(order_id, item_id, qty, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5)
		`,
		orderId,
		itemId,
		qty,
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &updatedItem, nil
}

func (r *PostgresDBRepo) GetItemsWithCursor(pageSize int, cursor time.Time) ([]*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
			SELECT id, name, stock, version,created_at, updated_at
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
			&tmp.Version,
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
			SELECT id, name, stock, version ,created_at, updated_at
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
			&tmp.Version,
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
			SELECT id, name, stock, version ,created_at, updated_at
			FROM items
			WHERE id = $1
		`

	row := r.DB.QueryRowContext(ctx, query, id)
	var item model.Item
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Stock,
		&item.Version,
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
		RETURNING id, name, stock, version ,created_at, updated_at`,
		item.Name,
		item.Stock,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(
		&createdItem.ID,
		&createdItem.Name,
		&createdItem.Stock,
		&createdItem.Version,
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

func (r *PostgresDBRepo) UpdateItem(id string, item model.UpdateItem) (*model.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var updatedItem model.Item

	err := r.DB.QueryRowContext(
		ctx,
		`UPDATE items SET
			name = $1,
			stock = $2,
			updated_at = $3,
			version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING id, name, stock, version ,created_at, updated_at`,
		item.Name,
		item.Stock,
		time.Now().UTC(),
		id,
		item.Version,
	).Scan(
		&updatedItem.ID,
		&updatedItem.Name,
		&updatedItem.Stock,
		&updatedItem.Version,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.ErrItemConflictVersion
		}
		return nil, err
	}

	return &updatedItem, nil
}
