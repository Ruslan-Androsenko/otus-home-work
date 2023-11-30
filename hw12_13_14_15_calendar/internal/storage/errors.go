package storage

import "errors"

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventDoesNotExist  = errors.New("event does non exist")
	ErrEventIDCantChanged = errors.New("event ID cannot be changed")
	ErrEventDateTimeBusy  = errors.New("time is already occupied by another event")
)
