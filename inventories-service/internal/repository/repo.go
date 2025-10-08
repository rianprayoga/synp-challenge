package repository

import (
	"inventories-app/internal/model"
	"time"
)

type DBRepo interface {
	GetItemsWithCursor(pageSize int, cursor time.Time) ([]*model.Item, error)
	GetItems(pageSize int) ([]*model.Item, error)
	GetItem(id string) (*model.Item, error)
	AddItem(item model.CreateItem) (*model.Item, error)
	DeleteItem(id string) error
}
