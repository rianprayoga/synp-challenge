package repository

import "orders-app/internal/model"

type DBRepo interface {
	AddOrder(item model.OrderRequest) (*model.Order, error)
}
