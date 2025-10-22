package service

import (
	"http-calendar/internal/storage"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	storage.Clear()

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
			name:        "valid event",
			userID:      "1",
			eventID:     "1",
			dateStr:     "2024-01-15",
			title:       "Meeting",
			description: "Team meeting",
			wantErr:     false,
		},
		{
			name:        "empty title",
			userID:      "1",
			eventID:     "2",
			dateStr:     "2024-01-15",
			title:       "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid user ID",
			userID:      "invalid",
			eventID:     "1",
			dateStr:     "2024-01-15",
			title:       "Meeting",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid event ID",
			userID:      "1",
			eventID:     "invalid",
			dateStr:     "2024-01-15",
			title:       "Meeting",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "invalid date",
			userID:      "1",
			eventID:     "1",
			dateStr:     "invalid-date",
			title:       "Meeting",
			description: "Description",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := CreateEvent(tt.userID, tt.eventID, tt.dateStr, tt.title, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && event == nil {
				t.Error("CreateEvent() returned nil event")
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	storage.Clear()

	event, err := CreateEvent("1", "1", "2024-01-15", "Original", "Original description")
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	if event == nil {
		t.Fatal("Created event is nil")
	}

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
			eventID:     "1",
			dateStr:     "2024-01-15",
			title:       "Updated",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "update with empty title",
			userID:      "1",
			eventID:     "1",
			dateStr:     "2024-01-15",
			title:       "",
			description: "Description",
			wantErr:     true,
		},
		{
			name:        "update non-existing event",
			userID:      "1",
			eventID:     "999",
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
			if !tt.wantErr && event == nil {
				t.Error("UpdateEvent() returned nil event")
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	storage.Clear()

	event, err := CreateEvent("1", "1", "2024-01-15", "Meeting", "Description")
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
	if event == nil {
		t.Fatal("Created event is nil")
	}

	tests := []struct {
		name    string
		userID  string
		eventID string
		wantErr bool
	}{
		{
			name:    "delete existing event",
			userID:  "1",
			eventID: "1",
			wantErr: false,
		},
		{
			name:    "delete non-existing event",
			userID:  "1",
			eventID: "999",
			wantErr: true,
		},
		{
			name:    "invalid user ID",
			userID:  "invalid",
			eventID: "1",
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

	if _, err := CreateEvent("1", "1", "2024-01-15", "Event 1", "Description 1"); err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	if _, err := CreateEvent("1", "2", "2024-01-15", "Event 2", "Description 2"); err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}
	if _, err := CreateEvent("1", "3", "2024-01-16", "Event 3", "Description 3"); err != nil {
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
			dateStr:   "2024-01-16",
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
			events, err := GetEventsForDay(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(events) != tt.wantCount {
				t.Errorf("GetEventsForDay() got %d events, want %d", len(events), tt.wantCount)
			}
		})
	}
}

func TestGetEventsForWeek(t *testing.T) {
	storage.Clear()

	if _, err := CreateEvent("1", "1", "2024-01-15", "Event 1", "Description 1"); err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	if _, err := CreateEvent("1", "2", "2024-01-20", "Event 2", "Description 2"); err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}

	tests := []struct {
		name    string
		userID  string
		dateStr string
		wantErr bool
	}{
		{
			name:    "get events for week",
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
			_, err := GetEventsForWeek(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForWeek() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEventsForMonth(t *testing.T) {
	storage.Clear()

	if _, err := CreateEvent("1", "1", "2024-01-15", "Event 1", "Description 1"); err != nil {
		t.Fatalf("Failed to create test event 1: %v", err)
	}
	if _, err := CreateEvent("1", "2", "2024-02-15", "Event 2", "Description 2"); err != nil {
		t.Fatalf("Failed to create test event 2: %v", err)
	}

	tests := []struct {
		name    string
		userID  string
		dateStr string
		wantErr bool
	}{
		{
			name:    "get events for month",
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
			_, err := GetEventsForMonth(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsForMonth() error = %v, wantErr %v", err, tt.wantErr)
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
			_, _, err := parseUserIDAndDate(tt.userID, tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUserIDAndDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAndParse(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		eventID string
		dateStr string
		title   string
		wantErr bool
	}{
		{
			name:    "valid input",
			userID:  "1",
			eventID: "1",
			dateStr: "2024-01-15",
			title:   "Meeting",
			wantErr: false,
		},
		{
			name:    "empty title",
			userID:  "1",
			eventID: "1",
			dateStr: "2024-01-15",
			title:   "",
			wantErr: true,
		},
		{
			name:    "invalid user ID",
			userID:  "invalid",
			eventID: "1",
			dateStr: "2024-01-15",
			title:   "Meeting",
			wantErr: true,
		},
		{
			name:    "invalid event ID",
			userID:  "1",
			eventID: "invalid",
			dateStr: "2024-01-15",
			title:   "Meeting",
			wantErr: true,
		},
		{
			name:    "invalid date",
			userID:  "1",
			eventID: "1",
			dateStr: "invalid-date",
			title:   "Meeting",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := validateAndParse(tt.userID, tt.eventID, tt.dateStr, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndParse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
