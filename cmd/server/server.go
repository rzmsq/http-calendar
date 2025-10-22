package main

import (
	"errors"
	"http-calendar/internal/config"
	"http-calendar/internal/logger"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_event", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /update_event", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /delete_event", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events_for_day", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events_for_week", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /events_for_month", func(w http.ResponseWriter, r *http.Request) {})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: logger.Middleware(mux, cfg.PathLog),
	}

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
