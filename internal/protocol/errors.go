package protocol

import "errors"

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrItemConflict    = errors.New("item already exists")
	ErrInvalidArgs     = errors.New("invalid arguments")
	ErrOperationFailed = errors.New("operation failed")
	ErrDatabase        = errors.New("database error")
	ErrInternal        = errors.New("internal error")
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrUnauthorized    = errors.New("unauthorized")
)
