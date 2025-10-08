package model

import "time"

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateItem struct {
	Name  string `json:"name" validate:"required,min=3"`
	Stock int    `json:"stock" validate:"required,gt=0"`
}
