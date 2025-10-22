package storage

import (
	"http-calendar/internal/models"
	"sync"
	"time"
)

var storage struct {
	mu sync.RWMutex
	m  map[uint64]map[uint64]models.Event
}

func CreateEvent(event *models.Event) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		storage.m = make(map[uint64]map[uint64]models.Event)
	}

	if storage.m[event.UserID] == nil {
		storage.m[event.UserID] = make(map[uint64]models.Event)
	}

	if _, exists := storage.m[event.UserID][event.EventID]; exists {
		return models.ErrExistingEvent
	}

	storage.m[event.UserID][event.EventID] = *event
	return nil
}

func UpdateEvent(event *models.Event) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return models.ErrUserNotFound
	} else {
		_, ok := storage.m[event.UserID][event.EventID]
		if !ok {
			return models.ErrEventNotFound
		}
	}
	storage.m[event.UserID][event.EventID] = *event
	return nil
}

func DeleteEvent(userID, eventID uint64) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return models.ErrUserNotFound
	} else {
		_, ok := storage.m[userID][eventID]
		if !ok {
			return models.ErrEventNotFound
		}
	}
	delete(storage.m[userID], eventID)
	return nil
}

func GetEventsForDay(userID uint64, date time.Time) ([]models.Event, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return nil, models.ErrUserNotFound
	}
	values, ok := storage.m[userID]
	if !ok {
		return nil, models.ErrUserNotFound
	}

	result := make([]models.Event, 0, len(values))
	for _, value := range values {
		if date.Equal(value.Date) {
			result = append(result, value)
		}
	}
	return result, nil
}

func GetEventsForWeek(userID uint64, startDate time.Time) ([]models.Event, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return nil, models.ErrUserNotFound
	}
	values, ok := storage.m[userID]
	if !ok {
		return nil, models.ErrUserNotFound
	}

	startOfWeek := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	result := make([]models.Event, 0, len(values))
	for _, value := range values {
		eventDate := time.Date(value.Date.Year(), value.Date.Month(), value.Date.Day(), 0, 0, 0, 0, value.Date.Location())
		if (eventDate.Equal(startOfWeek) || eventDate.After(startOfWeek)) && eventDate.Before(endOfWeek) {
			result = append(result, value)
		}
	}
	return result, nil
}

func GetEventsForMonth(userID uint64, startDate time.Time) ([]models.Event, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return nil, models.ErrUserNotFound
	}
	values, ok := storage.m[userID]
	if !ok {
		return nil, models.ErrUserNotFound
	}

	startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	result := make([]models.Event, 0, len(values))
	for _, value := range values {
		eventDate := time.Date(value.Date.Year(), value.Date.Month(), value.Date.Day(), 0, 0, 0, 0, value.Date.Location())
		if (eventDate.Equal(startOfMonth) || eventDate.After(startOfMonth)) && eventDate.Before(endOfMonth) {
			result = append(result, value)
		}
	}
	return result, nil
}

func GetNewEventID() uint64 {
	return uint64(time.Now().UnixNano())
}

func Clear() {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.m = make(map[uint64]map[uint64]models.Event)
}
