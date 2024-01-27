package storage

import (
	"time"
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
