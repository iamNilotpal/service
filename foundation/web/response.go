package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type ResponseStatus string

const (
	ResponseOK    ResponseStatus = "OK"
	ResponseError ResponseStatus = "ERROR"
)

type responseSuccess[T any] struct {
	Data   T              `json:"data"`
	Status ResponseStatus `json:"status"`
}

type responseError[T any] struct {
	Error  T              `json:"error"`
	Status ResponseStatus `json:"status"`
}

func NewSuccessResponse[T any](data T) *responseSuccess[T] {
	return &responseSuccess[T]{Data: data, Status: ResponseOK}
}

func NewErrorResponse[T any](err T) *responseError[T] {
	return &responseError[T]{Error: err, Status: ResponseError}
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}
