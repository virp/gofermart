package hash

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHash struct {
}

func NewPasswordHash() PasswordHash {
	return PasswordHash{}
}

func (ph PasswordHash) Make(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}

	return string(bytes), nil
}

func (ph PasswordHash) Check(ctx context.Context, password string, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil, nil
}
