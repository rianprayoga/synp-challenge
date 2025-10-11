package model

import "time"

type OrderRequest struct {
	ItemId string `json:"itemId" validate:"required"`
	Qty    int    `json:"qty" validate:"required,gt=0"`
}

type OrderResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type Order struct {
	ID          string
	ItemId      string
	Qty         int
	IsConfirmed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
