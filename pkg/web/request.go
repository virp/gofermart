package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Param(r *http.Request, key string) string {
	return chi.URLParamFromCtx(r.Context(), key)
}

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
