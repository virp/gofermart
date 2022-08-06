package order

import (
	"encoding/json"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	return json.Marshal(&struct {
		*OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: (*OrderAlias)(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}
