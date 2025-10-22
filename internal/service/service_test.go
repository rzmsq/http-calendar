package service

import (
	"http-calendar/internal/storage"
	"strconv"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	storage.Clear()

	tests := []struct {
		name        string
		userID      string
		dateStr     string
		title       string
		description string
		wantErr     bool
	}{
		{
			name:        "valid event",
			userID:      "1",
			dateStr:     "2024-01-15",
			title:       "Meeting",
			description: "Team meeting",
			wantErr:     false,
		},
		{
			name:        "empty title",
			userID:      "1",
			dateStr:     "2024-01-15",
			title:       "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid",
			dateStr:     "2024-01-15",
			title:       "Meeting",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid date",
			userID:      "1",
			dateStr:     "invalid-date",
			title:       "Meeting",
			description: "Description",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := CreateEvent(tt.userID, tt.dateStr, tt.title, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if event == nil {
					t.Error("CreateEvent() returned nil event")
					return
				}

				// Verify event fields
				expectedUserID, _ := strconv.ParseUint(tt.userID, 10, 64)
				expectedDate, _ := time.Parse(DateFormat, tt.dateStr)

				if event.UserID != expectedUserID {
					t.Errorf("Event has wrong userID = %v, want %v", event.UserID, expectedUserID)
				}
				if event.Title != tt.title {
					t.Errorf("Event has wrong title = %v, want %v", event.Title, tt.title)
				}
				if event.Description != tt.description {
					t.Errorf("Event has wrong description = %v, want %v", event.Description, tt.description)
				}
				if !event.Date.Equal(expectedDate) {
					t.Errorf("Event has wrong date = %v, want %v", event.Date, expectedDate)
				}
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	storage.Clear()

	// Create a test event first
	originalEvent, err := CreateEvent("1", "2024-01-15", "Original", "Original description")
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	eventIDStr := strconv.FormatUint(originalEvent.EventID, 10)

	tests := []struct {
		name        string
		userID      string
		eventID     string
		dateStr     string
		title       string
		description string
		wantErr     bool
	}{
		{
			name:        "update existing event",
			userID:      "1",
			eventID:     eventIDStr,
			dateStr:     "2024-01-10", // Change date
			title:       "Updated",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "update with empty title",
			userID:      "1",
			eventID:     eventIDStr,
			dateStr:     "2024-01-15",
			title:       "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "update non-existing event",
			userID:      "1",
			eventID:     "999999",
			dateStr:     "2024-01-15",
			title:       "Title",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid event ID",
			userID:      "1",
			eventID:     "invalid",
			dateStr:     "2024-01-15",
			title:       "Title",
			description: "Description",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := UpdateEvent(tt.userID, tt.eventID, tt.dateStr, tt.title, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && event != nil {
				// Verify updated fields
				if event.Title != tt.title {
					t.Errorf("Updated event has wrong title = %v, want %v", event.Title, tt.title)
				}
				if event.Description != tt.description {
					t.Errorf("Updated event has wrong description = %v, want %v", event.Description, tt.description)
				}

				expectedDate, _ := time.Parse(DateFormat, tt.dateStr)
				if !event.Date.Equal(expectedDate) {
					t.Errorf("Updated event has wrong date = %v, want %v", event.Date, expectedDate)
				}
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	storage.Clear()

	// Create a test event first
	event, err := CreateEvent("1", "2024-01-15", "Meeting", "Description")
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	eventIDStr := strconv.FormatUint(event.EventID, 10)

	tests := []struct {
		name    string
		userID  string
		eventID string
		wantErr bool
	}{
		{
			name:    "delete existing event",
			userID:  "1",
			eventID: eventIDStr,
			wantErr: false,
		},
		{
			name:    "delete already deleted event",
			userID:  "1",
			eventID: eventIDStr,
			wantErr: true,
		},
		{
			name:    "delete non-existing event",
			userID:  "1",
			eventID: "999999",
			wantErr: true,
		},
		{
			name:    "invalid user ID",
			userID:  "invalid",
			eventID: eventIDStr,
			wantErr: true,
		},
		{
			name:    "invalid event ID",
			userID:  "1",
			eventID: "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteEvent(tt.userID, tt.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEventsForDay(t *testing.T) {
	storage.Clear()

	// Create test events
	event1, err := CreateEvent("1", "2024-01-15", "Event 1", "Description 1")
	if err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	event2, err := CreateEvent("1", "2024-01-15", "Event 2", "Description 2")
	if err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}
	_, err = CreateEvent("1", "2024-01-10", "Event 3", "Description 3")
	if err != nil {
		t.Fatalf("Failed to create test event 3: %v", err)
	}

	tests := []struct {
		name      string
		userID    string
		dateStr   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get events for day with 2 events",
			userID:    "1",
			dateStr:   "2024-01-15",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "get events for day with 1 event",
			userID:    "1",
			dateStr:   "2024-01-10",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "get events for day with no events",
			userID:    "1",
			dateStr:   "2024-01-01",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "invalid user ID",
			userID:    "invalid",
			dateStr:   "2024-01-15",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "invalid date",
			userID:    "1",
			dateStr:   "invalid-date",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := GetEventsForDay(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(events) != tt.wantCount {
					t.Errorf("GetEventsForDay() got %d events, want %d", len(events), tt.wantCount)
				}

				// Check if returned events match the expected date
				if len(events) > 0 {
					expectedDate, _ := time.Parse(DateFormat, tt.dateStr)
					for i, e := range events {
						if !e.Date.Equal(expectedDate) {
							t.Errorf("Event[%d] has wrong date = %v, want %v", i, e.Date, expectedDate)
						}
					}

					// For the first day (2024-01-15), verify we got our two specific events
					if tt.dateStr == "2024-01-15" {
						found1, found2 := false, false
						for _, e := range events {
							if e.EventID == event1.EventID {
								found1 = true
							}
							if e.EventID == event2.EventID {
								found2 = true
							}
						}
						if !found1 || !found2 {
							t.Errorf("Did not find expected events in results")
						}
					}
				}
			}
		})
	}
}

func TestGetEventsForWeek(t *testing.T) {
	storage.Clear()

	// Create test events with days within valid hour range (0-23)
	_, err := CreateEvent("1", "2024-01-15", "Event 1", "Description 1")
	if err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	_, err = CreateEvent("1", "2024-01-17", "Event 2", "Description 2")
	if err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}
	_, err = CreateEvent("1", "2024-01-23", "Event 3", "Description 3")
	if err != nil {
		t.Fatalf("Failed to create test event 3: %v", err)
	}

	tests := []struct {
		name      string
		userID    string
		dateStr   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get events for week starting 2024-01-15",
			userID:    "1",
			dateStr:   "2024-01-15",
			wantCount: 2, // Should include events on 15th and 17th
			wantErr:   false,
		},
		{
			name:      "get events for week starting 2024-01-22",
			userID:    "1",
			dateStr:   "2024-01-22",
			wantCount: 1, // Should include only event on 23rd
			wantErr:   false,
		},
		{
			name:      "invalid user ID",
			userID:    "invalid",
			dateStr:   "2024-01-15",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "invalid date",
			userID:    "1",
			dateStr:   "invalid-date",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := GetEventsForWeek(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForWeek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(events) != tt.wantCount {
				t.Errorf("GetEventsForWeek() got %d events, want %d", len(events), tt.wantCount)
			}
		})
	}
}

func TestGetEventsForMonth(t *testing.T) {
	storage.Clear()

	_, err := CreateEvent("1", "2024-01-15", "Event 1", "Description 1")
	if err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	_, err = CreateEvent("1", "2024-01-20", "Event 2", "Description 2")
	if err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}
	_, err = CreateEvent("1", "2024-02-15", "Event 3", "Description 3")
	if err != nil {
		t.Fatalf("Failed to create test event 3: %v", err)
	}

	tests := []struct {
		name      string
		userID    string
		dateStr   string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "get events for January 2024",
			userID:    "1",
			dateStr:   "2024-01-15",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "get events for February 2024",
			userID:    "1",
			dateStr:   "2024-02-01",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "invalid user ID",
			userID:    "invalid",
			dateStr:   "2024-01-15",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "invalid date",
			userID:    "1",
			dateStr:   "invalid-date",
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := GetEventsForMonth(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(events) != tt.wantCount {
				t.Errorf("GetEventsForMonth() got %d events, want %d", len(events), tt.wantCount)
			}
		})
	}
}

func TestParseUserIDAndDate(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		dateStr string
		wantErr bool
	}{
		{
			name:    "valid input",
			userID:  "1",
			dateStr: "2024-01-15",
			wantErr: false,
		},
		{
			name:    "invalid user ID",
			userID:  "invalid",
			dateStr: "2024-01-15",
			wantErr: true,
		},
		{
			name:    "invalid date",
			userID:  "1",
			dateStr: "invalid-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, date, err := parseUserIDAndDate(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUserIDAndDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedUID, _ := strconv.ParseUint(tt.userID, 10, 64)
				expectedDate, _ := time.Parse(DateFormat, tt.dateStr)

				if userID != expectedUID {
					t.Errorf("parseUserIDAndDate() userID = %v, want %v", userID, expectedUID)
				}
				if !date.Equal(expectedDate) {
					t.Errorf("parseUserIDAndDate() date = %v, want %v", date, expectedDate)
				}
			}
		})
	}
}

func TestValidateAndParse(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		dateStr string
		title   string
		wantErr bool
	}{
		{
			name:    "valid input",
			userID:  "1",
			dateStr: "2024-01-15",
			title:   "Meeting",
			wantErr: false,
		},
		{
			name:    "empty title",
			userID:  "1",
			dateStr: "2024-01-15",
			title:   "",
			wantErr: true,
		},
		{
			name:    "invalid user ID",
			userID:  "invalid",
			dateStr: "2024-01-15",
			title:   "Meeting",
			wantErr: true,
		},
		{
			name:    "invalid date",
			userID:  "1",
			dateStr: "invalid-date",
			title:   "Meeting",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, date, err := validateAndParse(tt.userID, tt.dateStr, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedUID, _ := strconv.ParseUint(tt.userID, 10, 64)
				expectedDate, _ := time.Parse(DateFormat, tt.dateStr)

				if userID != expectedUID {
					t.Errorf("validateAndParse() userID = %v, want %v", userID, expectedUID)
				}
				if !date.Equal(expectedDate) {
					t.Errorf("validateAndParse() date = %v, want %v", date, expectedDate)
				}
			}
		})
	}
}
