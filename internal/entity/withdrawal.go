package entity

import (
	"time"
)

type Withdrawal struct {
	ID          int       `db:"id"`
	UserID      string    `db:"user_id"`
	Sum         float64   `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}
