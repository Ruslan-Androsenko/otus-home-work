package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	memStorage *Storage
	events     = []storage.Event{
		{
			ID:    uuid.NewString(),
			Title: "First my event",
			Date:  time.Now(),
		},
		{
			ID:    uuid.NewString(),
			Title: "Second my event",
			Date:  time.Now(),
		},
		{
			ID:    uuid.NewString(),
			Title: "Third my event",
			Date:  time.Now(),
		},
		{
			ID:    uuid.NewString(),
			Title: "Fourth my event",
			Date:  time.Now(),
		},
		{
			ID:    uuid.NewString(),
			Title: "Fifth my event",
			Date:  time.Now(),
		},
	}
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	memStorage = New()
	require.NoError(t, memStorage.Connect(ctx))

	// Наполняем хранилище событиями
	fillEventsTest(ctx, t)

	// Редактирование события
	updEventIndex := 2
	require.Less(t, updEventIndex, len(events))
	require.NotNil(t, events[updEventIndex])
	updateEventTest(ctx, t, events[updEventIndex].ID)

	// Удаляем события
	deleteEventTest(ctx, t, events[1].ID)
	deleteEventTest(ctx, t, events[4].ID)

	// Закррываем соедиение с хранилищем
	require.NoError(t, memStorage.Close())
}

// Наполнение хранилища.
func fillEventsTest(ctx context.Context, t *testing.T) {
	t.Helper()

	// Создание новых событий
	for _, eventItem := range events {
		require.NotNil(t, eventItem)
		err := memStorage.CreateEvent(ctx, eventItem)
		require.NoError(t, err)
		checkEventItemTest(t, eventItem.ID, eventItem)
	}
}

// Проверяем наличие добавленного события в хранилище.
func checkEventItemTest(t *testing.T, eventID string, event storage.Event) {
	t.Helper()

	require.True(t, memStorage.hasExistsByID(eventID))
	require.True(t, memStorage.hasExistsByDate(event.Date))

	// Получаем событие из хранилища по его ID
	findEvent, err := memStorage.GetEventByID(eventID)
	require.NoError(t, err)
	require.NotNil(t, findEvent)

	// Проверяем содержимое основных полей
	require.Equal(t, eventID, findEvent.ID)
	require.Equal(t, event.Title, findEvent.Title)
	require.Equal(t, event.Date, findEvent.Date)
}

// Редактирование события.
func updateEventTest(ctx context.Context, t *testing.T, eventID string) {
	t.Helper()

	updEvent := storage.Event{
		ID:    uuid.NewString(),
		Title: "Update event",
		Date:  time.Now(),
	}

	err := memStorage.UpdateEvent(ctx, eventID, updEvent)
	require.NoError(t, err)
	checkEventItemTest(t, eventID, updEvent)
}

// Удалние события.
func deleteEventTest(ctx context.Context, t *testing.T, eventID string) {
	t.Helper()

	err := memStorage.DeleteEvent(ctx, eventID)
	require.NoError(t, err)
	require.False(t, memStorage.hasExistsByID(eventID))
}
