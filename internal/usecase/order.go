package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/virp/gofermart/internal/entity"
	"github.com/virp/gofermart/pkg/luhn"
)

var (
	ErrInvalidOrderNumber              = errors.New("invalid order number")
	ErrAlreadyUploadedOrder            = errors.New("already uploaded order by user")
	ErrAlreadyUploadedOrderByOtherUser = errors.New("already uploaded order by other user")
)

type OrderUseCase struct {
	repository OrderRepository
	queue      chan<- int
}

func NewOrderUseCase(repository OrderRepository, queue chan<- int) OrderUseCase {
	return OrderUseCase{
		repository: repository,
		queue:      queue,
	}
}

func (uc OrderUseCase) Upload(ctx context.Context, number int, userID string) (entity.Order, error) {
	if !luhn.IsValid(number) {
		return entity.Order{}, ErrInvalidOrderNumber
	}

	alreadyExistOrder, err := uc.repository.GetByNumber(ctx, number)
	if err == nil {
		if alreadyExistOrder.UserID == userID {
			return entity.Order{}, ErrAlreadyUploadedOrder
		}
		return entity.Order{}, ErrAlreadyUploadedOrderByOtherUser
	}

	order := entity.Order{
		ID:         number,
		UserID:     userID,
		Status:     entity.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	order, err = uc.repository.Create(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	go func(orderNumber int) {
		uc.queue <- orderNumber
	}(order.ID)

	return order, nil
}

func (uc OrderUseCase) List(ctx context.Context, userID string) ([]entity.Order, error) {
	orders, err := uc.repository.GetAllForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user orders list: %w", err)
	}

	return orders, nil
}
