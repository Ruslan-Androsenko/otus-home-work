package storage

import (
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api/proto"
)

const EventTableName = "event"

type Event struct {
	ID           string
	Title        string
	Date         time.Time
	DateEnd      time.Time
	Description  string
	OwnerID      int
	Notification time.Duration
}

// MakeStorageEvent Сформировать объект события в формате хранилища.
func MakeStorageEvent(event *proto.Event) Event {
	return Event{
		ID:           event.Id,
		Title:        event.Title,
		Date:         MakeTime(event.Date),
		DateEnd:      MakeTime(event.DateEnd),
		Description:  event.Description,
		OwnerID:      int(event.OwnerId),
		Notification: event.Notification.AsDuration(),
	}
}

// MakeProtoEvent Сформировать объект события в proto формате.
func MakeProtoEvent(event Event) *proto.Event {
	return &proto.Event{
		Id:           event.ID,
		Title:        event.Title,
		Date:         MakeProtoDateTime(event.Date),
		DateEnd:      MakeProtoDateTime(event.DateEnd),
		Description:  event.Description,
		OwnerId:      int32(event.OwnerID),
		Notification: MakeProtoDuration(event.Notification),
	}
}

// MakeProtoEventsList Сформировать массив объектов событий в proto формате.
func MakeProtoEventsList(events []Event) []*proto.Event {
	eventsList := make([]*proto.Event, 0, len(events))

	for _, event := range events {
		eventsList = append(eventsList, &proto.Event{
			Id:           event.ID,
			Title:        event.Title,
			Date:         MakeProtoDateTime(event.Date),
			DateEnd:      MakeProtoDateTime(event.DateEnd),
			Description:  event.Description,
			OwnerId:      int32(event.OwnerID),
			Notification: MakeProtoDuration(event.Notification),
		})
	}

	return eventsList
}
