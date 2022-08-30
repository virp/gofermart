package migrations

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var schema string

func Migrate(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	return nil
}
