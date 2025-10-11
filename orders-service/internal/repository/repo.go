package repository

import "orders-app/internal/model"

type DBRepo interface {
	AddOrder(orderId string, isConfirmed bool, item model.OrderRequest) (*model.Order, error)
}
