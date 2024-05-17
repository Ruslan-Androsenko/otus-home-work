package internalhttp

import (
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type EventParams struct {
	ID    string `json:"id"`
	Date  string `json:"date"`
	Event Event  `json:"event"`
}

type Event struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Date         string        `json:"date"`
	DateEnd      string        `json:"date_end"`
	Description  string        `json:"description"`
	OwnerID      int           `json:"owner_id"`
	Notification time.Duration `json:"notification"`
}

// Получить параметр даты для списка событий.
func (params EventParams) getDate() (time.Time, error) {
	date, errParse := time.ParseInLocation(storage.DateFormat, params.Date, time.Local)
	if errParse != nil {
		return time.Time{}, errParse
	}

	return date, errParse
}

// Получить объект события из переданных параметров запроса.
func (params EventParams) getEvent() (*storage.Event, error) {
	eventID := params.ID

	if params.Event.ID != "" {
		eventID = params.Event.ID
	}

	date, errParse := time.ParseInLocation(storage.DateTimeFormat, params.Event.Date, time.Local)
	if errParse != nil {
		return nil, errParse
	}

	dateEnd, errParse := time.ParseInLocation(storage.DateTimeFormat, params.Event.DateEnd, time.Local)
	if errParse != nil {
		return nil, errParse
	}

	return &storage.Event{
		ID:           eventID,
		Title:        params.Event.Title,
		Date:         date,
		DateEnd:      dateEnd,
		Description:  params.Event.Description,
		OwnerID:      params.Event.OwnerID,
		Notification: params.Event.Notification,
	}, nil
}
