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

var ErrTrxNotFound = &HttpError{
	Err:        errors.New("transaction not found"),
	StatusCode: http.StatusNotFound,
}

var ErrItemConflictVersion = &HttpError{
	Err:        errors.New("item version mismatch"),
	StatusCode: http.StatusConflict,
}

var ErrInternalServer = errors.New("something unexpected")
