package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/virp/gofermart/internal/entity"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = uuid.NewString()
	_, err := r.db.NamedExecContext(
		ctx,
		"insert into users (id, login, password_hash) VALUES (:id, :login, :password_hash)",
		user,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("insert user to db: %w", err)
	}

	return user, nil
}

func (r UserRepository) Update(ctx context.Context, user entity.User) error {
	_, err := r.db.NamedExecContext(
		ctx,
		"update users set balance = :balance, withdrawn = :withdrawn where id = :id",
		user,
	)
	if err != nil {
		return fmt.Errorf("updated user in db: %w", err)
	}

	return nil
}

func (r UserRepository) GetByID(ctx context.Context, ID string) (entity.User, error) {
	var user entity.User

	err := r.db.GetContext(
		ctx,
		&user,
		"select id, login, password_hash, balance, withdrawn from users where id = $1",
		ID,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("get user from db: %w", err)
	}

	return user, nil
}

func (r UserRepository) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	var user entity.User

	err := r.db.GetContext(
		ctx,
		&user,
		"select id, login, password_hash, balance, withdrawn from users where login = $1",
		login,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("get user from db: %w", err)
	}

	return user, nil
}
