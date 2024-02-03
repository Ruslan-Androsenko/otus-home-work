package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

// Получение параметров запроса.
func loadParams(r *http.Request) (map[string]interface{}, error) {
	var params map[string]interface{}

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

	return params, err
}

// Сформировать параметры для создания события.
func makeCreateParams(params map[string]interface{}) (storage.Event, error) {
	var (
		event    storage.Event
		errParse error
	)

	if val, ok := params["event"]; ok {
		if eventParams, okCast := val.(map[string]interface{}); okCast {
			event, errParse = storage.MakeEventFromParams(eventParams)
		}
	} else {
		event, errParse = storage.MakeEventFromParams(params)
	}

	return event, errParse
}

// Сформировать параметры для обновления события.
func makeUpdateParams(params map[string]interface{}) (string, storage.Event, error) {
	var (
		eventID  string
		event    storage.Event
		errParse error
	)

	if val, ok := params["id"]; ok {
		if id, okCast := val.(string); okCast {
			eventID = id
		}
	}

	if val, ok := params["event"]; ok {
		if eventParams, okCast := val.(map[string]interface{}); okCast {
			event, errParse = storage.MakeEventFromParams(eventParams)
		}
	}

	return eventID, event, errParse
}

// Сформировать параметры для получения события по его ID.
func makeGetEventByIDParams(params map[string]interface{}) string {
	var eventID string

	if val, ok := params["id"]; ok {
		if id, okCast := val.(string); okCast {
			eventID = id
		}
	}

	return eventID
}

// Сформировать параметры для получения списка событий по дате.
func makeGetEventsByDateParams(params map[string]interface{}) (time.Time, error) {
	var (
		date     time.Time
		errParse error
	)

	if val, ok := params["date"]; ok {
		if eventDate, okCast := val.(string); okCast {
			date, errParse = time.ParseInLocation(storage.DateFormat, eventDate, time.Local)
			if errParse != nil {
				return time.Time{}, errParse
			}
		}
	}

	return date, errParse
}
