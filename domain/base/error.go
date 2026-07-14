package domain

import "errors"

var (
	ErrNotFound     = errors.New("data not found")
	ErrConflict     = errors.New("data already exists")
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
)

func ErrorHandler(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return 404, err.Error()
	case errors.Is(err, ErrConflict):
		return 409, err.Error()
	case errors.Is(err, ErrUnauthorized):
		return 401, err.Error()
	case errors.Is(err, ErrBadRequest):
		return 400, err.Error()
	default:
		return 500, ErrInternal.Error()
	}
}
