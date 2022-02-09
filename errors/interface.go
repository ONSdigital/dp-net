package errors

import (
	"github.com/pkg/errors"
)

type dataLogger interface {
	LogData() map[string]interface{}
}

type messager interface {
	Message() string
}

type stacktracer interface {
	StackTrace() errors.StackTrace
}

type coder interface {
	Code() int
}

type errNotFound interface {
	NotFound() bool
}

type errConflict interface {
	Conflict() bool
}

type errUnavailable interface {
	Unavailable() bool
}

type errUnauthorized interface {
	Unauthorized() bool
}

type errTimeout interface {
	Timeout() bool
}
