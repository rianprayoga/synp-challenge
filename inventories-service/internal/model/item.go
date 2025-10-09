package model

import "time"

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateItem struct {
	Name  string `json:"name" validate:"required,min=3"`
	Stock int    `json:"stock" validate:"required,gt=0"`
}

type UpdateItem struct {
	Name    string `json:"name" validate:"required,min=3"`
	Stock   int    `json:"stock" validate:"required,gt=0"`
	Version int    `json:"version" validate:"required,gte=1"`
}
