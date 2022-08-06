package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/virp/gofermart/internal/handlers/middleware"
	"github.com/virp/gofermart/internal/usecase"
	"github.com/virp/gofermart/pkg/web"
)

type Handlers struct {
	User                  usecase.User
	AppSecret             string
	AppUserAuthCookieName string
}

func (h Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ru Register
	if err := web.Decode(r, &ru); err != nil {
		return middleware.NewBadApiRequestError("decode user register request")
	}

	user, err := h.User.Register(ctx, ru.Login, ru.Password)
	if err != nil {
		return err
	}

	err = web.SetEncryptedCookie(w, h.AppUserAuthCookieName, user.ID, h.AppSecret)
	if err != nil {
		return fmt.Errorf("set cookie: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func (h Handlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var lu Login
	if err := web.Decode(r, &lu); err != nil {
		return middleware.NewBadApiRequestError("decode user login request")
	}

	user, err := h.User.Login(ctx, lu.Login, lu.Password)
	if err != nil {
		return err
	}

	err = web.SetEncryptedCookie(w, h.AppUserAuthCookieName, user.ID, h.AppSecret)
	if err != nil {
		return fmt.Errorf("set cookie: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func (h Handlers) Balance(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	user, err := h.User.GetByID(ctx, v.UserID)
	if err != nil {
		return fmt.Errorf("get user by id[%s]: %w", v.UserID, err)
	}

	balance := Balance{
		Current:   user.Balance,
		Withdrawn: user.Withdrawn,
	}

	return web.Respond(ctx, w, balance, http.StatusOK)
}
