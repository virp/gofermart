package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/virp/gofermart/internal/entity"
)

type UserUseCase struct {
	repository   UserRepository
	passwordHash PasswordHash
}

var (
	ErrInvalidLoginOrPassword = errors.New("invalid user login or password")
	ErrUserAlreadyExist       = errors.New("user already exist")
	ErrUserNotFound           = errors.New("user not found")
)

func NewUserUseCase(repository UserRepository, passwordHash PasswordHash) UserUseCase {
	return UserUseCase{
		repository:   repository,
		passwordHash: passwordHash,
	}
}

func (uc UserUseCase) Register(ctx context.Context, login string, password string) (entity.User, error) {
	_, err := uc.repository.GetByLogin(ctx, login)
	if err == nil {
		return entity.User{}, ErrUserAlreadyExist
	}

	passwordHash, err := uc.passwordHash.Make(ctx, password)
	if err != nil {
		return entity.User{}, fmt.Errorf("make password hash: %w", err)
	}
	user := entity.User{
		Login:        login,
		PasswordHash: passwordHash,
	}
	user, err = uc.repository.Create(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (uc UserUseCase) Login(ctx context.Context, login string, password string) (entity.User, error) {
	user, err := uc.repository.GetByLogin(ctx, login)
	if err != nil {
		return entity.User{}, ErrInvalidLoginOrPassword
	}

	isValidPassword, err := uc.passwordHash.Check(ctx, password, user.PasswordHash)
	if err != nil {
		return entity.User{}, fmt.Errorf("check password hash: %w", err)
	}
	if !isValidPassword {
		return entity.User{}, ErrInvalidLoginOrPassword
	}

	return user, nil
}

func (uc UserUseCase) GetByID(ctx context.Context, ID string) (entity.User, error) {
	user, err := uc.repository.GetByID(ctx, ID)
	if err != nil {
		return entity.User{}, ErrUserNotFound
	}

	return user, nil
}
