package storage

import (
	"time"
)

type Event struct {
	ID           string
	Title        string
	Date         time.Time
	DateEnd      time.Time
	Description  string
	OwnerID      string
	Notification time.Duration
}
