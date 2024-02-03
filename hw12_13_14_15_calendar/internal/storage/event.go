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

// MakeEventFromParams Сформировать объект события в формате хранилища из параметров запроса.
func MakeEventFromParams(eventParams map[string]interface{}) (Event, error) {
	var (
		event    Event
		errParse error
	)

	if val, ok := eventParams["id"]; ok {
		if eventID, okCast := val.(string); okCast {
			event.ID = eventID
		}
	}

	if val, ok := eventParams["title"]; ok {
		if title, okCast := val.(string); okCast {
			event.Title = title
		}
	}

	if val, ok := eventParams["date"]; ok {
		if eventDate, okCast := val.(string); okCast {
			event.Date, errParse = time.ParseInLocation(DateTimeFormat, eventDate, time.Local)
			if errParse != nil {
				return Event{}, errParse
			}
		}
	}

	if val, ok := eventParams["date_end"]; ok {
		if eventDateEnd, okCast := val.(string); okCast {
			event.DateEnd, errParse = time.ParseInLocation(DateTimeFormat, eventDateEnd, time.Local)
			if errParse != nil {
				return Event{}, errParse
			}
		}
	}

	if val, ok := eventParams["description"]; ok {
		if description, okCast := val.(string); okCast {
			event.Description = description
		}
	}

	if val, ok := eventParams["owner_id"]; ok {
		if ownerID, okCast := val.(float64); okCast {
			event.OwnerID = int(ownerID)
		}
	}

	if val, ok := eventParams["notification"]; ok {
		if notification, okCast := val.(float64); okCast {
			event.Notification = time.Duration(notification)
		}
	}

	return event, nil
}
