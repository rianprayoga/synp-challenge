package appError

import "errors"

var ErrItemNotFound = errors.New("item not found")
var ErrUnexpected = errors.New("something unexpected")
