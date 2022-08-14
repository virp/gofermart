package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	_ = SetStatusCode(ctx, statusCode)

	if data == nil {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
