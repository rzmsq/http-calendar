package models

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date")
	ErrExistingEvent = errors.New("existing event")
)

type Event struct {
	UserID      uint
	EventID     uint
	Date        time.Time
	Title       string
	Description string
}
