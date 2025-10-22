package service

import (
	"errors"
	"http-calendar/internal/models"
	"http-calendar/internal/storage"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02"

func CreateEvent(userID, eventID, dateStr, title, description string) (*models.Event, error) {
	uID, eID, date, err := validateAndParse(userID, eventID, dateStr, title)
	if err != nil || errors.Is(err, models.ErrTitleIsRequired) {
		return nil, err
	}

	event := models.NewEvent(uID, eID, date, title, description)
	err = storage.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func UpdateEvent(userID, eventID, dateStr, title, description string) (*models.Event, error) {
	uID, eID, date, err := validateAndParse(userID, eventID, dateStr, title)
	if err != nil || errors.Is(err, models.ErrTitleIsRequired) {
		return nil, err
	}

	event := models.NewEvent(uID, eID, date, title, description)
	err = storage.UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func DeleteEvent(userID, eventID string) error {
	uID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}
	eID, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		return err
	}
	return storage.DeleteEvent(uID, eID)
}

func GetEventsForDay(userID, dateStr string) ([]models.Event, error) {
	uID, date, err := parseUserIDAndDate(userID, dateStr)
	if err != nil {
		return nil, err
	}
	return storage.GetEventsForDay(uID, date)
}

func GetEventsForWeek(userID, dateStr string) ([]models.Event, error) {
	uID, date, err := parseUserIDAndDate(userID, dateStr)
	if err != nil {
		return nil, err
	}
	return storage.GetEventsForWeek(uID, date)
}

func GetEventsForMonth(userID, dateStr string) ([]models.Event, error) {
	uID, date, err := parseUserIDAndDate(userID, dateStr)
	if err != nil {
		return nil, err
	}
	return storage.GetEventsForMonth(uID, date)
}

func parseUserIDAndDate(userID, dateStr string) (uint64, time.Time, error) {
	uID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return 0, time.Time{}, err
	}

	return uID, date, nil
}

func validateAndParse(userID string, eventID string, dateStr string, title string) (uint64, uint64, time.Time, error) {
	uID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	eID, err := strconv.ParseUint(eventID, 10, 64)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	if title == "" {
		return 0, 0, time.Time{}, models.ErrTitleIsRequired
	}
	return uID, eID, date, nil
}
