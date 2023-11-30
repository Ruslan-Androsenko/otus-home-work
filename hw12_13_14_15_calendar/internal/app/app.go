package app

import (
	"context"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Fatal(msg string)
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)

	Fatalf(format string, values ...any)
	Errorf(format string, values ...any)
	Warningf(format string, values ...any)
	Infof(format string, values ...any)
	Debugf(format string, values ...any)
}

type Storage interface {
	CreateEvent(event storage.Event) error
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error

	GetEventById(id string) (storage.Event, error)
	GetEventsOfDay(date time.Time) ([]storage.Event, error)
	GetEventsOfWeek(date time.Time) ([]storage.Event, error)
	GetEventsOfMonth(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

// CreateEvent Создать событие.
func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	err := a.storage.CreateEvent(event)
	if err == nil {
		a.logger.Debugf("Event has been created, id: %s", event.ID)
	}

	return err
}

// UpdateEvent Изменить событие.
func (a *App) UpdateEvent(id string, event storage.Event) error {
	if id != event.ID {
		return storage.ErrEventIDCantChanged
	}

	err := a.storage.UpdateEvent(id, event)
	if err == nil {
		a.logger.Debugf("Event has been updated, id: %s", id)
	}

	return err
}

// DeleteEvent Удалить событие.
func (a *App) DeleteEvent(id string) error {
	err := a.storage.DeleteEvent(id)
	if err == nil {
		a.logger.Debugf("Event has been deleted, id: %s", id)
	}

	return err
}

// GetEventById Получить событие по его ID.
func (a *App) GetEventById(id string) (storage.Event, error) {
	return a.storage.GetEventById(id)
}

// GetEventsOfDay Получить события за день.
func (a *App) GetEventsOfDay(date time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsOfDay(date)
}

// GetEventsOfWeek Получить события за неделю.
func (a *App) GetEventsOfWeek(date time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsOfWeek(date)
}

// GetEventsOfMonth  Получить события за месяц.
func (a *App) GetEventsOfMonth(date time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsOfMonth(date)
}
