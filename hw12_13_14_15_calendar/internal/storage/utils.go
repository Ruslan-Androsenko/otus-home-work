package storage

import (
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api/proto"
	"github.com/golang/protobuf/ptypes/duration"
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

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

// MakeDateRange Сформировать необходимый диапазон дат.
func MakeDateRange(date time.Time, period Period) (dateFrom, dateTo time.Time) {
	var defaultDate time.Time
	currentDay := date.Format(DateFormat)
	dateFrom, err := time.ParseInLocation(DateTimeFormat, currentDay+" 00:00:00", time.Local)
	if err != nil {
		return defaultDate, defaultDate
	}

	if days, ok := daysOfRange[period]; ok {
		dateTo = dateFrom.Add(time.Hour * 24 * days) //nolint:durationcheck
	} else {
		return defaultDate, defaultDate
	}

	return dateFrom, dateTo
}

// MakeTime Сформировать объект для даты и времени.
func MakeTime(dateTime *proto.DateTime) time.Time {
	if dateTime != nil {
		return time.Date(
			int(dateTime.Year),
			time.Month(dateTime.Month),
			int(dateTime.Day),
			int(dateTime.Hours),
			int(dateTime.Minutes),
			int(dateTime.Seconds),
			int(dateTime.Nanos),
			time.Local,
		)
	}

	return time.Time{}
}

// MakeProtoDateTime Сформировать объект для даты и времени для ответа на запрос.
func MakeProtoDateTime(dateTime time.Time) *proto.DateTime {
	return &proto.DateTime{
		Year:    int32(dateTime.Year()),
		Month:   int32(dateTime.Month()),
		Day:     int32(dateTime.Day()),
		Hours:   int32(dateTime.Hour()),
		Minutes: int32(dateTime.Minute()),
		Seconds: int32(dateTime.Second()),
		Nanos:   int32(dateTime.Nanosecond()),
	}
}

// MakeProtoDuration Сформировать объект продолжительности для ответа на запрос.
func MakeProtoDuration(dur time.Duration) *duration.Duration {
	return &duration.Duration{
		Seconds: int64(dur),
	}
}
