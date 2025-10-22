package storage

import (
	"errors"
	"http-calendar/internal/models"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	// Очищаем storage перед тестом
	storage.m = nil

	event := &models.Event{
		EventID: 1,
		UserID:  1,
		Title:   "Test Event",
		Date:    time.Now(),
	}

	err := CreateEvent(event)
	if err != nil {
		t.Errorf("CreateEvent() error = %v", err)
	}

	// Проверяем, что событие создано
	if storage.m[1][1].Title != "Test Event" {
		t.Errorf("Event not created correctly")
	}
}

func TestCreateEvent_ExistingEvent(t *testing.T) {
	storage.m = map[uint]map[uint]models.Event{
		1: {
			1: {EventID: 1, UserID: 1, Title: "Existing"},
		},
	}

	event := &models.Event{EventID: 1, UserID: 1, Title: "Duplicate"}

	err := CreateEvent(event)
	if !errors.Is(err, models.ErrExistingEvent) {
		t.Errorf("Expected ErrExistingEvent, got %v", err)
	}
}

func TestUpdateEvent(t *testing.T) {
	storage.m = map[uint]map[uint]models.Event{
		1: {
			1: {EventID: 1, UserID: 1, Title: "Old Title"},
		},
	}

	event := &models.Event{EventID: 1, UserID: 1, Title: "New Title"}

	err := UpdateEvent(event)
	if err != nil {
		t.Errorf("UpdateEvent() error = %v", err)
	}

	if storage.m[1][1].Title != "New Title" {
		t.Errorf("Event not updated correctly")
	}
}

func TestUpdateEvent_NotFound(t *testing.T) {
	storage.m = nil

	event := &models.Event{EventID: 1, UserID: 1}

	err := UpdateEvent(event)
	if !errors.Is(err, models.ErrUserNotFound) {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	storage.m = map[uint]map[uint]models.Event{
		1: {1: {EventID: 1, UserID: 1}},
	}

	err := DeleteEvent(1)
	if err != nil {
		t.Errorf("DeleteEvent() error = %v", err)
	}

	if _, exists := storage.m[1]; exists {
		t.Errorf("User events not deleted")
	}
}

func TestGetEventsForDay(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	storage.m = map[uint]map[uint]models.Event{
		1: {
			1: {EventID: 1, UserID: 1, Date: testDate},
			2: {EventID: 2, UserID: 1, Date: testDate.AddDate(0, 0, 1)},
		},
	}

	events, err := GetEventsForDay(1, testDate)
	if err != nil {
		t.Errorf("GetEventsForDay() error = %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestGetEventsForWeek(t *testing.T) {
	startDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	storage.m = map[uint]map[uint]models.Event{
		1: {
			1: {EventID: 1, UserID: 1, Date: startDate},
			2: {EventID: 2, UserID: 1, Date: startDate.AddDate(0, 0, 3)},
			3: {EventID: 3, UserID: 1, Date: startDate.AddDate(0, 0, 8)},
		},
	}

	events, err := GetEventsForWeek(1, startDate)
	if err != nil {
		t.Errorf("GetEventsForWeek() error = %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	startDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	storage.m = map[uint]map[uint]models.Event{
		1: {
			1: {EventID: 1, UserID: 1, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			2: {EventID: 2, UserID: 1, Date: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)},
			3: {EventID: 3, UserID: 1, Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	events, err := GetEventsForMonth(1, startDate)
	if err != nil {
		t.Errorf("GetEventsForMonth() error = %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}
