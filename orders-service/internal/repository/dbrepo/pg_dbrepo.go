package dbrepo

import (
	"context"
	"database/sql"
	"orders-app/internal/model"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = 3 * time.Second

func (r *PostgresDBRepo) AddOrder(orderId string, isConfirmed bool, item model.OrderRequest) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var order model.Order
	err := r.DB.QueryRowContext(
		ctx,
		`
		INSERT INTO orders 
			(id, item_id, qty, isConfirmed, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6 )
		RETURNING id, isConfirmed
		`,
		orderId,
		item.ItemId,
		item.Qty,
		isConfirmed,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(
		&order.ID,
		&order.IsConfirmed,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
