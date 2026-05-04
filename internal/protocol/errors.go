package protocol

import "errors"

var (
	ErrInvalidArgs          = errors.New("invalid arguments")
	ErrTelemetryUnavailable = errors.New("telemetry unavailable")
	ErrTelemetryFailed      = errors.New("telemetry collection failed")
	ErrInternal             = errors.New("internal error")
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrUnauthorized    = errors.New("unauthorized")
)
