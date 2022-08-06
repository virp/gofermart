package accrual

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	statusPath = "/api/orders/%d"
)

type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type SDK struct {
	Address string
}

func New(address string) SDK {
	return SDK{
		Address: address,
	}
}

func (a *SDK) GetStatus(number int) (OrderStatus, int, error) {
	fullPath := fmt.Sprintf(a.Address+statusPath, number)

	res, err := http.Get(fullPath)
	if err != nil {
		return OrderStatus{}, 0, fmt.Errorf("request: %w", err)
	}

	if res.StatusCode == http.StatusTooManyRequests {
		retryHeader := res.Header.Get("Retry-After")
		retry, err := strconv.Atoi(retryHeader)
		if err != nil {
			return OrderStatus{}, 0, fmt.Errorf("missed Retry-After header")
		}

		return OrderStatus{}, retry, nil
	}

	if res.StatusCode != http.StatusOK {
		return OrderStatus{}, 0, fmt.Errorf("unknown status code: %d", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return OrderStatus{}, 0, fmt.Errorf("read body: %w", err)
	}

	var status OrderStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return OrderStatus{}, 0, fmt.Errorf("unmarshalling: %w", err)
	}

	return status, 0, nil
}
