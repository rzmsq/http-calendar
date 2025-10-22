# HTTP Calendar

A simple HTTP-based calendar application built with Go.

## Project Overview

This project implements an HTTP calendar service that allows users to manage events and appointments through a RESTful API. It serves as a learning project demonstrating Go web application development practices.

## Task Requirements

Implement an HTTP server for a small calendar of events with the following functionality:

### CRUD Operations
- `POST /create_event` — Create a new event
- `POST /update_event` — Update an existing event
- `POST /delete_event` — Delete an event
- `GET /events_for_day` — Retrieve all events for a specific day
- `GET /events_for_week` — Retrieve events for a week
- `GET /events_for_month` — Retrieve events for a month

### Request Format
- Data for creation/updating is passed in the request body as either URL-form (`application/x-www-form-urlencoded`) or JSON
- Required parameters may include: `user_id`, `date` (YYYY-MM-DD), `event` (text)
- For GET requests, parameters can be passed via query string (e.g., `?user_id=1&date=2023-12-31`)

### Response Format
- Successful execution: JSON format `{"result": "..."}`
- Business logic error: JSON format `{"error": "error description"}`

### HTTP Status Codes
- `200 OK` for successful requests
- `400` for input errors (e.g., incorrect date format)
- `503` for business logic errors (e.g., trying to delete a non-existent event)
- `500` for other errors

### Additional Requirements
- Logging middleware that logs each request (method, URL, time) to stdout or a file
- Server port configurable through environment variables or flags
- Business logic separate from the HTTP layer
- Clean, verified code (vet, lint) without data races
- Unit tests for core business logic functions

## Features

- Calendar event management (create, read, update, delete)
- RESTful API endpoints
- Configuration system
- Persistent storage
- Logging

## Project Structure

- `cmd/server`: Application entrypoint
- `config`: Configuration files
- `internal`: Core application code
    - `config`: Configuration management
    - `handler`: HTTP request handlers
    - `logger`: Logging functionality
    - `models`: Data models
    - `service`: Business logic
    - `storage`: Data persistence
- `logs`: Application logs
