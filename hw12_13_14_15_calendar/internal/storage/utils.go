package storage

import (
	"time"
)

type Period uint8

const (
	PeriodDay Period = iota + 1
	PeriodWeek
	PeriodMonth
)

var daysOfRange = map[Period]time.Duration{
	PeriodDay:   1,
	PeriodWeek:  7,
	PeriodMonth: 30,
}

// MakeDateRange Сформировать необходимый диапазон дат.
func MakeDateRange(date time.Time, period Period) (dateFrom, dateTo time.Time) {
	var defaultDate time.Time
	currentDay := date.Format("2006-01-02")
	dateFrom, err := time.ParseInLocation("2006-01-02 15:04:05", currentDay+" 00:00:00", time.Local)
	if err != nil {
		return defaultDate, defaultDate
	}

	if days, ok := daysOfRange[period]; ok {
		dateTo = dateFrom.Add(time.Hour * 24 * days)
	} else {
		return defaultDate, defaultDate
	}

	return dateFrom, dateTo
}
