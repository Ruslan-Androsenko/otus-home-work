package storage

import "errors"

var (
	ErrEventDoesNotExist  = errors.New("event does non exist")
	ErrEventIDCantChanged = errors.New("event ID cannot be changed")
	ErrEventDateTimeBusy  = errors.New("time is already occupied by another event")
)
