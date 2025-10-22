package main

import (
	"context"
	"errors"
	"http-calendar/internal/config"
	"http-calendar/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	serverError := make(chan error, 1)
	log.Printf("starting http server on port %s\n", cfg.Port)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 3)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	select {
	case err := <-serverError:
		log.Fatalf("http server error: %s\n", err)
	case <-stop:
		log.Println("stopping http server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("http server shutdown error: %s\n", err)
	}
	log.Println("http server shutdown complete")
}
