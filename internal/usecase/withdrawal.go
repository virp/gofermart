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
	ErrInvalidWithdrawalNumber = errors.New("invalid withdrawal number")
	ErrNotEnoughBalance        = errors.New("not enough balance")
)

type WithdrawalUseCase struct {
	withdrawalRepository WithdrawalRepository
	userRepository       UserRepository
}

func NewWithdrawalUseCase(withdrawalRepository WithdrawalRepository, userRepository UserRepository) WithdrawalUseCase {
	return WithdrawalUseCase{
		withdrawalRepository: withdrawalRepository,
		userRepository:       userRepository,
	}
}

func (uc WithdrawalUseCase) Upload(ctx context.Context, number int, sum float64, userID string) (entity.Withdrawal, error) {
	if !luhn.IsValid(number) {
		return entity.Withdrawal{}, ErrInvalidWithdrawalNumber
	}

	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		return entity.Withdrawal{}, fmt.Errorf("user not found: %w", err)
	}

	if user.Balance < sum {
		return entity.Withdrawal{}, ErrNotEnoughBalance
	}

	withdrawal := entity.Withdrawal{
		ID:          number,
		UserID:      userID,
		Sum:         sum,
		ProcessedAt: time.Now(),
	}

	withdrawal, err = uc.withdrawalRepository.Create(ctx, withdrawal)
	if err != nil {
		return entity.Withdrawal{}, err
	}

	user.Balance = user.Balance - sum
	user.Withdrawn = user.Withdrawn + sum
	err = uc.userRepository.Update(ctx, user)
	if err != nil {
		return entity.Withdrawal{}, err
	}

	return withdrawal, nil
}

func (uc WithdrawalUseCase) List(ctx context.Context, userID string) ([]entity.Withdrawal, error) {
	withdrawals, err := uc.withdrawalRepository.GetAllForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
