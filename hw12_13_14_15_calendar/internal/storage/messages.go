package storage

import (
	"fmt"
	"time"
)

func MakeMessageOfSuccessfullyEventsList(typeQuery string, date time.Time, countEvents int) string {
	return fmt.Sprintf("Events for the specified %s date: %v were successfully received, count: %d",
		typeQuery, date, countEvents)
}
