package entity

import (
	"database/sql"
	"time"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	ID         int             `db:"id"`
	UserID     string          `db:"user_id"`
	Status     string          `db:"status"`
	UploadedAt time.Time       `db:"uploaded_at"`
	Accrual    sql.NullFloat64 `db:"accrual"`
}
