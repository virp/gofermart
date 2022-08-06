package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/virp/gofermart/pkg/web"
)

var ErrUserUnauthorized = errors.New("user unauthorized")

func Auth(cookie string, secret string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			userID, err := web.GetEncryptedCookie(r, cookie, secret)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return ErrUserUnauthorized
				}
				return fmt.Errorf("get encrypted cookie: %w", err)
			}
			v.UserID = userID

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
