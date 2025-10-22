package models

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEventNotFound   = errors.New("event not found")
	ErrTitleIsRequired = errors.New("title is required")
	ErrInvalidDate     = errors.New("invalid date")
	ErrExistingEvent   = errors.New("existing event")
)

type Event struct {
	UserID      uint64
	EventID     uint64
	Date        time.Time
	Title       string
	Description string
}

func NewEvent(userID uint64, eventID uint64, date time.Time, title, description string) *Event {
	return &Event{
		UserID:      userID,
		EventID:     eventID,
		Date:        date,
		Title:       title,
		Description: description,
	}
}
