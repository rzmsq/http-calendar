package storage

import (
	"http-calendar/internal/models"
	"sync"
	"time"
)

var storage struct {
	mu sync.RWMutex
	m  map[uint]map[uint]models.Event
}

func CreateEvent(event *models.Event) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		storage.m = make(map[uint]map[uint]models.Event)
	} else {
		_, ok := storage.m[event.UserID][event.EventID]
		if ok {
			return models.ErrExistingEvent
		}
	}
	storage.m[event.UserID] = make(map[uint]models.Event)
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

func DeleteEvent(userID uint) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return models.ErrUserNotFound
	} else {
		_, ok := storage.m[userID]
		if !ok {
			return models.ErrUserNotFound
		}
	}
	delete(storage.m, userID)
	return nil
}

func GetEventsForDay(userID uint, date time.Time) ([]models.Event, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storage.m == nil {
		return nil, models.ErrUserNotFound
	}
	values, ok := storage.m[userID]
	if !ok {
		return nil, models.ErrUserNotFound
	}

	result := make([]models.Event, 0)
	for _, value := range values {
		if value.Date == date {
			result = append(result, value)
		}
	}
	return result, nil
}

func GetEventsForWeek(userID uint, startDate time.Time) ([]models.Event, error) {
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

	result := make([]models.Event, 0)
	for _, value := range values {
		eventDate := time.Date(value.Date.Year(), value.Date.Month(), value.Date.Day(), 0, 0, 0, 0, value.Date.Location())
		if (eventDate.Equal(startOfWeek) || eventDate.After(startOfWeek)) && eventDate.Before(endOfWeek) {
			result = append(result, value)
		}
	}
	return result, nil
}

func GetEventsForMonth(userID uint, startDate time.Time) ([]models.Event, error) {
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

	result := make([]models.Event, 0)
	for _, value := range values {
		eventDate := time.Date(value.Date.Year(), value.Date.Month(), value.Date.Day(), 0, 0, 0, 0, value.Date.Location())
		if (eventDate.Equal(startOfMonth) || eventDate.After(startOfMonth)) && eventDate.Before(endOfMonth) {
			result = append(result, value)
		}
	}
	return result, nil
}
