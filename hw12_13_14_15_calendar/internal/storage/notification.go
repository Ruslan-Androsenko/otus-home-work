package storage

import "time"

type Notification struct {
	ID      string
	Title   string
	Date    time.Time
	OwnerID string
}
