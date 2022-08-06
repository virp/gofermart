package withdrawal

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w *Withdrawal) MarshalJSON() ([]byte, error) {
	type WithdrawalAlias Withdrawal
	return json.Marshal(&struct {
		*WithdrawalAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawalAlias: (*WithdrawalAlias)(w),
		ProcessedAt:     w.ProcessedAt.Format(time.RFC3339),
	})
}
