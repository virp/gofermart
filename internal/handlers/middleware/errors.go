package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/virp/gofermart/internal/usecase"
	"github.com/virp/gofermart/pkg/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			if err := handler(ctx, w, r); err != nil {
				log.Errorw(
					"ERROR",
					"traceID", v.TraceID,
					"message", err,
				)

				var status int
				var message string
				switch {
				case IsBadAPIRequestError(err):
					status = http.StatusBadRequest
				case errors.Is(err, ErrUserUnauthorized):
					status = http.StatusUnauthorized
				case errors.Is(err, usecase.ErrInvalidLoginOrPassword):
					status = http.StatusUnauthorized
				case errors.Is(err, usecase.ErrUserAlreadyExist):
					status = http.StatusConflict
				case errors.Is(err, usecase.ErrUserNotFound):
					status = http.StatusNotFound
				case errors.Is(err, usecase.ErrInvalidOrderNumber):
					status = http.StatusUnprocessableEntity
				case errors.Is(err, usecase.ErrAlreadyUploadedOrder):
					status = http.StatusOK
				case errors.Is(err, usecase.ErrAlreadyUploadedOrderByOtherUser):
					status = http.StatusConflict
				case errors.Is(err, usecase.ErrInvalidWithdrawalNumber):
					status = http.StatusUnprocessableEntity
				case errors.Is(err, usecase.ErrNotEnoughBalance):
					status = http.StatusPaymentRequired
				default:
					status = http.StatusInternalServerError
				}

				message = http.StatusText(status)

				w.WriteHeader(status)
				_, wErr := w.Write([]byte(message))
				if wErr != nil {
					return wErr
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}

type BadAPIRequestError struct {
	Message string
}

func NewBadAPIRequestError(message string) error {
	return BadAPIRequestError{message}
}

func (e BadAPIRequestError) Error() string {
	return e.Message
}

func IsBadAPIRequestError(err error) bool {
	var e BadAPIRequestError
	return errors.As(err, &e)
}
