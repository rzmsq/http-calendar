package models

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEventNotFound   = errors.New("event not found")
	ErrTitleIsRequired = errors.New("title is required")
	ErrExistingEvent   = errors.New("existing event")
)

type Event struct {
	UserID      uint64    `json:"user_id"`
	EventID     uint64    `json:"event_id"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
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
