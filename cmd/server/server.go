package main

import (
	"context"
	"errors"
	"http-calendar/internal/config"
	"http-calendar/internal/handler"
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
	run(cfg)
}

func run(cfg *config.Config) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_event", handler.CreateHandler)
	mux.HandleFunc("POST /update_event", handler.UpdateHandler)
	mux.HandleFunc("POST /delete_event", handler.DeleteHandler)
	mux.HandleFunc("GET /events_for_day", handler.GetEventsForDayHandler)
	mux.HandleFunc("GET /events_for_week", handler.GetEventsForWeekHandler)
	mux.HandleFunc("GET /events_for_month", handler.GetEventsForMonthHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: logger.Middleware(mux, cfg.PathLog),
	}

	serverError := make(chan error, 1)
	log.Printf("starting http server on port %s\n", httpServer.Addr)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverError:
		log.Fatalf("http server error: %v\n", err)
	case sig := <-stop:
		log.Printf("stopping http server by signal %v\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("http server shutdown error: %s\n", err)
	}
	log.Println("http server shutdown complete")
}
