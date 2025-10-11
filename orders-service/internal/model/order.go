package model

import "time"

type OrderRequest struct {
	ItemId string `json:"itemId" validate:"required"`
	Qty    int    `json:"qty" validate:"required,gt=0"`
}

type Order struct {
	ID          string `json:"id"`
	ItemId      string `json:"itemId"`
	Qty         int    `json:"qty"`
	IsConfirmed bool
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
