package order

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/virp/gofermart/internal/handlers/middleware"
	"github.com/virp/gofermart/internal/usecase"
	"github.com/virp/gofermart/pkg/web"
)

type Handlers struct {
	Order usecase.Order
}

func (h Handlers) Upload(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return middleware.NewBadAPIRequestError("get order upload request body")
	}

	orderNumber, err := strconv.Atoi(string(b))
	if err != nil {
		return middleware.NewBadAPIRequestError("order number is not integer")
	}

	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	_, err = h.Order.Upload(ctx, orderNumber, v.UserID)
	if err != nil {
		return fmt.Errorf("upload order: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusAccepted)
}

func (h Handlers) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	entityOrders, err := h.Order.List(ctx, v.UserID)
	if err != nil {
		return fmt.Errorf("get user[%s] orders: %w", v.UserID, err)
	}
	var orders []Order
	for _, o := range entityOrders {
		order := Order{
			Number:     strconv.Itoa(o.ID),
			Status:     o.Status,
			Accrual:    o.Accrual.Float64,
			UploadedAt: o.UploadedAt,
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}

	return web.Respond(ctx, w, orders, http.StatusOK)
}
