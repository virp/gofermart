package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/virp/gofermart/internal/entity"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r OrderRepository) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	_, err := r.db.NamedExecContext(
		ctx,
		"insert into orders (id, user_id, status, uploaded_at) values (:id, :user_id, :status, :uploaded_at)",
		order,
	)
	if err != nil {
		return entity.Order{}, fmt.Errorf("insert order to db: %w", err)
	}

	return order, nil
}

func (r OrderRepository) GetAllForUser(ctx context.Context, userID string) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.SelectContext(
		ctx,
		&orders,
		"select * from orders where user_id = $1 order by uploaded_at",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query user orders from db: %w", err)
	}

	return orders, nil
}

func (r OrderRepository) GetByNumber(ctx context.Context, number int) (entity.Order, error) {
	var order entity.Order

	err := r.db.GetContext(
		ctx,
		&order,
		"select id, user_id, status, uploaded_at, accrual from orders where id = $1",
		number,
	)
	if err != nil {
		return entity.Order{}, fmt.Errorf("get order from db: %w", err)
	}

	return order, nil
}

func (r OrderRepository) UpdateStatus(ctx context.Context, number int, status string, accrual float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	var dbAccrual sql.NullFloat64
	if accrual > 0.0 {
		_ = dbAccrual.Scan(accrual)
	}

	_, err = tx.ExecContext(
		ctx,
		"update orders set status = $1, accrual = $2 where id = $3",
		status,
		dbAccrual,
		number,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	if accrual > 0.0 {
		_, err = tx.ExecContext(
			ctx,
			"update users set balance = balance + $1 where id = (select user_id from orders where id = $2)",
			accrual,
			number,
		)
		if err != nil {
			return fmt.Errorf("update user balance: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
