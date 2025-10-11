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

func (r *PostgresDBRepo) AddOrder(item model.OrderRequest) (*model.Order, error) {
	_, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	return nil, nil
}
