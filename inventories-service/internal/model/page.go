package model

type PageResponse[T any] struct {
	NextCursor string `json:"nextCursor,omitempty"`
	Data       []T    `json:"data,omitempty"`
}
