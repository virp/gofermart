package withdrawal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/virp/gofermart/internal/handlers/middleware"
	"github.com/virp/gofermart/internal/usecase"
	"github.com/virp/gofermart/pkg/web"
)

type Handlers struct {
	Withdrawal usecase.Withdrawal
}

func (h Handlers) Withdraw(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var wd Withdraw
	if err := web.Decode(r, &wd); err != nil {
		return middleware.NewBadAPIRequestError("decode withdraw request")
	}

	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	orderNumber, err := strconv.Atoi(wd.Order)
	if err != nil {
		return middleware.NewBadAPIRequestError("withdraw order number not an integer")
	}

	_, err = h.Withdrawal.Upload(ctx, orderNumber, wd.Sum, v.UserID)
	if err != nil {
		return fmt.Errorf("upload withdraw: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusOK)
}

func (h Handlers) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	entityWithdrawals, err := h.Withdrawal.List(ctx, v.UserID)
	if err != nil {
		return fmt.Errorf("get user withdrawals: %w", err)
	}

	var withdrawals []Withdrawal
	for _, w := range entityWithdrawals {
		withdrawal := Withdrawal{
			Order:       strconv.Itoa(w.ID),
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if len(withdrawals) == 0 {
		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}

	return web.Respond(ctx, w, withdrawals, http.StatusOK)
}
