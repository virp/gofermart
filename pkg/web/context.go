package web

import (
	"context"
	"errors"
	"time"
)

type ctxKey int

const key ctxKey = 1

type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
	UserID     string
}

func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return nil, errors.New("web value missing from context")
	}

	return v, nil
}

func SetStatusCode(ctx context.Context, statusCode int) error {
	v, err := GetValues(ctx)
	if err != nil {
		return err
	}
	v.StatusCode = statusCode

	return nil
}
