package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu   sync.RWMutex //nolint:unused
	data map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		data: make(map[string]storage.Event),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

// Проверяем имеется ли событие по ID.
func (s *Storage) hasExistsById(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[id]

	return ok
}

// Проверяем имеется ли событие по указанной дате.
func (s *Storage) hasExistsByDate(date time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.data {
		if event.Date == date {
			return true
		}
	}

	return false
}

// CreateEvent Создать событие.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	if s.hasExistsByDate(event.Date) {
		return storage.ErrEventDateTimeBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var temp time.Time

	if event.Date == temp {
		event.Date = time.Now()
	}

	s.data[event.ID] = event
	return nil
}

// UpdateEvent Изменить событие.
func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	if !s.hasExistsById(id) {
		return storage.ErrEventDoesNotExist
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = event
	return nil
}

// DeleteEvent Удалить событие.
func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	if !s.hasExistsById(id) {
		return storage.ErrEventDoesNotExist
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}

// GetEventById Получить событие по его ID.
func (s *Storage) GetEventById(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if event, ok := s.data[id]; ok {
		return event, nil
	}

	return storage.Event{}, storage.ErrEventDoesNotExist
}

// GetEventsOfDay Получить события за день.
func (s *Storage) GetEventsOfDay(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodDay)
}

// GetEventsOfWeek Получить события за неделю.
func (s *Storage) GetEventsOfWeek(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodWeek)
}

// GetEventsOfMonth  Получить события за месяц.
func (s *Storage) GetEventsOfMonth(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodMonth)
}

// FindEventsByPeriod Поиск списка подходящих событий по периоду.
func (s *Storage) FindEventsByPeriod(date time.Time, period storage.Period) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	eventsList := NewEventsList(len(s.data))
	dateFrom, dateTo := storage.MakeDateRange(date, period)

	for _, event := range s.data {
		if event.Date.After(dateFrom) && event.Date.Before(dateTo) {
			eventsList.List = append(eventsList.List, event)
		}
	}

	eventsList.SortByDateAsc()
	return eventsList.List, nil
}
