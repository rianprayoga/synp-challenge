package model

import "time"

type Item struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
