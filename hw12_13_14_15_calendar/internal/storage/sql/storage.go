package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	dbConn *sql.DB
	logger *logger.Logger
}

func New(dbConn *sql.DB, logg *logger.Logger) *Storage {
	return &Storage{
		dbConn: dbConn,
		logger: logg,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.dbConn.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.dbConn.Close()
}

// Проверяем имеется ли событие по ID.
func (s *Storage) hasExistsById(id string) (bool, error) {
	var countRows int

	sqlQuery := fmt.Sprintf("select count(*) as countRows from `%s` where id = '%s';",
		storage.EventTableName, id)

	row := s.dbConn.QueryRow(sqlQuery)
	if err := row.Scan(&countRows); err != nil {
		return false, fmt.Errorf("can not get count rows. Error: %v", err)
	}

	return countRows > 0, nil
}

// Проверяем имеется ли событие по указанной дате.
func (s *Storage) hasExistsByDate(date time.Time) (bool, error) {
	var countRows int

	sqlQuery := fmt.Sprintf("select count(*) as countRows from `%s` where date = '%s';",
		storage.EventTableName, date.Format(storage.DateTimeFormat))

	row := s.dbConn.QueryRow(sqlQuery)
	if err := row.Scan(&countRows); err != nil {
		return false, fmt.Errorf("can not get count rows. Error: %v", err)
	}

	return countRows > 0, nil
}

// CreateEvent Создать событие.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	exists, err := s.hasExistsByDate(event.Date)
	if err != nil {
		return err
	}

	if exists {
		return storage.ErrEventDateTimeBusy
	}

	var temp time.Time

	if event.Date == temp {
		event.Date = time.Now()
	}

	fields := fmt.Sprintf("'%s', '%s', '%s', '%s', '%s', %d, %d",
		event.ID, event.Title, event.Date.Format(storage.DateTimeFormat),
		event.DateEnd.Format(storage.DateTimeFormat),
		event.Description, event.OwnerID, event.Notification)

	sqlQuery := fmt.Sprintf("insert into `%s` values(%s);", storage.EventTableName, fields)
	if _, err = s.dbConn.ExecContext(ctx, sqlQuery); err != nil {
		return fmt.Errorf("can not create event. Error: %v", err)
	}

	return nil
}

// UpdateEvent Изменить событие.
func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	exists, err := s.hasExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrEventDoesNotExist
	}

	fields := fmt.Sprintf("title = '%s', date = '%s', date_end = '%s', "+
		"description = '%s', owner_id = %d, notification = %d",
		event.Title, event.Date.Format(storage.DateTimeFormat),
		event.DateEnd.Format(storage.DateTimeFormat),
		event.Description, event.OwnerID, event.Notification)

	sqlQuery := fmt.Sprintf("update `%s` set %s where id = '%s';", storage.EventTableName, fields, id)
	if _, err = s.dbConn.ExecContext(ctx, sqlQuery); err != nil {
		return fmt.Errorf("can not update event, ID: %s. Error: %v", id, err)
	}

	return nil
}

// DeleteEvent Удалить событие.
func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	exists, err := s.hasExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return storage.ErrEventDoesNotExist
	}

	sqlQuery := fmt.Sprintf("delete from `%s` where id = '%s';", storage.EventTableName, id)
	if _, err = s.dbConn.ExecContext(ctx, sqlQuery); err != nil {
		return fmt.Errorf("can not delete event, ID: %s. Error: %v", id, err)
	}

	return nil
}

// GetEventById Получить событие по его ID.
func (s Storage) GetEventById(id string) (storage.Event, error) {
	var (
		event        storage.Event
		defaultEvent storage.Event
	)

	fields := "id, title, date, date_end, description, owner_id, notification"
	sqlQuery := fmt.Sprintf("select %s from `%s` where id = '%s';",
		fields, storage.EventTableName, id)

	row := s.dbConn.QueryRow(sqlQuery)
	if err := row.Scan(&event.ID, &event.Title, &event.Date, &event.DateEnd,
		&event.Description, &event.OwnerID, &event.Notification); err != nil {
		return defaultEvent, err
	}

	if event.ID == "" {
		return defaultEvent, storage.ErrEventDoesNotExist
	}

	return event, nil
}

// GetEventsOfDay Получить события за день.
func (s Storage) GetEventsOfDay(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodDay)
}

// GetEventsOfWeek Получить события за неделю.
func (s Storage) GetEventsOfWeek(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodWeek)
}

// GetEventsOfMonth  Получить события за месяц.
func (s Storage) GetEventsOfMonth(date time.Time) ([]storage.Event, error) {
	return s.FindEventsByPeriod(date, storage.PeriodMonth)
}

// FindEventsByPeriod Поиск списка подходящих событий по периоду.
func (s *Storage) FindEventsByPeriod(date time.Time, period storage.Period) ([]storage.Event, error) {
	var (
		events []storage.Event
		event  storage.Event
	)

	dateFrom, dateTo := storage.MakeDateRange(date, period)
	fields := "id, title, date, date_end, description, owner_id, notification"
	sqlQuery := fmt.Sprintf("select %s from `%s` where date between '%s' and '%s';",
		fields, storage.EventTableName,
		dateFrom.Format(storage.DateTimeFormat), dateTo.Format(storage.DateTimeFormat))

	rows, err := s.dbConn.Query(sqlQuery)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			s.logger.Errorf("Can not close of rows. Error: %v", err)
		}
	}()

	for rows.Next() {
		if err = rows.Scan(&event.ID, &event.Title, &event.Date, &event.DateEnd,
			&event.Description, &event.OwnerID, &event.Notification); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
