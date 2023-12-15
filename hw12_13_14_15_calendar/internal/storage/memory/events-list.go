package memorystorage

import (
	"sort"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type EventsList struct {
	List []storage.Event
}

// NewEventsList Создать новый список событий указанной емкости.
func NewEventsList(capacity int) EventsList {
	return EventsList{
		List: make([]storage.Event, 0, capacity),
	}
}

// SortByDateAsc Отсортировать список событий по возрастанию даты создания.
func (events *EventsList) SortByDateAsc() {
	sort.Slice(events.List, func(i, j int) bool {
		return events.List[i].Date.Before(events.List[j].Date)
	})
}

// SortByDateDesc Отсортировать список событий по убыванию даты создания.
func (events *EventsList) SortByDateDesc() {
	sort.Slice(events.List, func(i, j int) bool {
		return events.List[i].Date.After(events.List[j].Date)
	})
}
