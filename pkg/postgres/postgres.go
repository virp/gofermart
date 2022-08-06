package postgres

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, URI string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", URI)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	return db, nil
}
