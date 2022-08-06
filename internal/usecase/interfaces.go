package usecase

import (
	"context"

	"github.com/virp/gofermart/internal/entity"
)

type User interface {
	Register(ctx context.Context, login string, password string) (entity.User, error)
	Login(ctx context.Context, login string, password string) (entity.User, error)
	GetByID(ctx context.Context, ID string) (entity.User, error)
}

type Order interface {
	Upload(ctx context.Context, number int, userID string) (entity.Order, error)
	List(ctx context.Context, userID string) ([]entity.Order, error)
}

type Withdrawal interface {
	Upload(ctx context.Context, number int, sum float64, userID string) (entity.Withdrawal, error)
	List(ctx context.Context, userID string) ([]entity.Withdrawal, error)
}

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, ID string) (entity.User, error)
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	GetAllForUser(ctx context.Context, userID string) ([]entity.Order, error)
	GetByNumber(ctx context.Context, number int) (entity.Order, error)
	UpdateStatus(ctx context.Context, number int, status string, accrual float64) error
}

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal entity.Withdrawal) (entity.Withdrawal, error)
	GetAllForUser(ctx context.Context, userID string) ([]entity.Withdrawal, error)
	GetByNumber(ctx context.Context, number int) (entity.Withdrawal, error)
}

type PasswordHash interface {
	Make(ctx context.Context, password string) (string, error)
	Check(ctx context.Context, password string, passwordHash string) (bool, error)
}
