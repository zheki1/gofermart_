package helpers

import (
	"context"
	"errors"

	"gofermart_/internal/http/middleware"
)

var ErrUnauthorized = errors.New("unauthorized")

func GetUserID(ctx context.Context) (int, error) {
	val := ctx.Value(middleware.UserIDKey)

	uid, ok := val.(int)
	if !ok {
		return 0, ErrUnauthorized
	}

	return uid, nil
}
