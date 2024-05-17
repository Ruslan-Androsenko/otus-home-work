package storage

import "errors"

var (
	ErrEventDoesNotExist  = errors.New("event does not exist")
	ErrEventIDCantChanged = errors.New("event ID cannot be changed")
	ErrEventDateTimeBusy  = errors.New("time is already occupied by another event")
	ErrSQLNoRows          = "sql: no rows in result set"
)
