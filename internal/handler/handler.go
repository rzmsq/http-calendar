package handler

import (
	"encoding/json"
	"errors"
	"http-calendar/internal/models"
	"http-calendar/internal/service"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Result *models.Event `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendSuccess(w http.ResponseWriter, model *models.Event) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(SuccessResponse{Result: model})
	if err != nil {
		log.Printf("Failed encode response: %v\n", err)
	}
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	if err != nil {
		log.Printf("Failed encode response: %v\n", err)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
	}

	uid := r.FormValue("user_id")
	date := r.FormValue("date")
	title := r.FormValue("title")
	description := r.FormValue("description")

	log.Println(uid)
	model, err := service.CreateEvent(uid, date, title, description)
	if errors.Is(err, models.ErrTitleIsRequired) {
		sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, model)
}
