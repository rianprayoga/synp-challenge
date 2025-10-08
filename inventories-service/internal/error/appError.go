package appError

import (
	"errors"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Err        error
}

func (r *HttpError) Error() string {
	return r.Err.Error()
}

var ErrItemNotFound = &HttpError{
	Err:        errors.New("item not found"),
	StatusCode: http.StatusNotFound,
}

var ErrUnexpected = errors.New("something unexpected")
