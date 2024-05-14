package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Получение параметров запроса.
func loadParams[T any](r *http.Request) (*T, error) {
	var params T

	defer func() {
		if err := r.Body.Close(); err != nil {
			logg.Errorf("Cannot close body: %v", err)
		}
	}()

	buffer := make([]byte, 1024)
	read, err := r.Body.Read(buffer)
	if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to load request params. Error: %w", err)
	}

	data := buffer[:read]
	errDecode := json.Unmarshal(data, &params)
	if errDecode != nil {
		return nil, fmt.Errorf("failed to deserialize data: %s. Error: %w", string(data), errDecode)
	}

	return &params, err
}
