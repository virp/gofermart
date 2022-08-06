package usecase

import (
	"context"
	"time"

	"github.com/virp/gofermart/internal/entity"
	"github.com/virp/gofermart/pkg/accrual"
)

type OrderStatusUseCase struct {
	repository OrderRepository
	accrual    accrual.SDK
	queue      <-chan int
}

func NewOrderStatusUseCase(repository OrderRepository, accrual accrual.SDK, queue <-chan int, workersCount int) OrderStatusUseCase {
	uc := OrderStatusUseCase{
		repository: repository,
		accrual:    accrual,
		queue:      queue,
	}

	for i := 0; i < workersCount; i++ {
		go func() {
			for order := range uc.queue {
				uc.checkOrderStatus(order)
			}
		}()
	}

	return uc
}

func (uc *OrderStatusUseCase) checkOrderStatus(number int) {
	status, delay, err := uc.accrual.GetStatus(number)
	if err != nil {
		return
	}
	if delay > 0 {
		time.AfterFunc(
			time.Duration(delay)*time.Second,
			func() { uc.checkOrderStatus(number) },
		)
		return
	}

	err = uc.repository.UpdateStatus(context.Background(), number, status.Status, status.Accrual)
	if err != nil {
		return
	}

	if status.Status == entity.OrderStatusProcessing {
		uc.checkOrderStatus(number)
	}
}
