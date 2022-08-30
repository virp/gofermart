package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/virp/gofermart/internal/entity"
)

type WithdrawalRepository struct {
	db *sqlx.DB
}

func NewWithdrawalRepository(db *sqlx.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

func (r WithdrawalRepository) Create(ctx context.Context, withdrawal entity.Withdrawal) (entity.Withdrawal, error) {
	_, err := r.db.NamedExecContext(
		ctx,
		"insert into withdrawals (id, user_id, sum, processed_at) values (:id, :user_id, :sum, :processed_at)",
		withdrawal,
	)
	if err != nil {
		return entity.Withdrawal{}, fmt.Errorf("insert withdrawal to db: %w", err)
	}

	return withdrawal, nil
}

func (r WithdrawalRepository) GetAllForUser(ctx context.Context, userID string) ([]entity.Withdrawal, error) {
	var withdrawals []entity.Withdrawal
	err := r.db.SelectContext(
		ctx,
		&withdrawals,
		"select * from withdrawals where user_id = $1 order by processed_at",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query user withdrawals from db: %w", err)
	}

	return withdrawals, nil
}

func (r WithdrawalRepository) GetByNumber(ctx context.Context, number int) (entity.Withdrawal, error) {
	var withdrawal entity.Withdrawal

	err := r.db.GetContext(
		ctx,
		&withdrawal,
		"select * from withdrawals where id = $1",
		number,
	)
	if err != nil {
		return entity.Withdrawal{}, fmt.Errorf("get withdrawal from db: %w", err)
	}

	return withdrawal, nil
}
