package restserver

import (
	"car-dealership/car/internal/application"
	"car-dealership/car/internal/application/model"
	"encoding/json"
	"fmt"
	"net/http"
	"uuid"
)

type Handler struct {
	application *application.CarService
}

func NewHandler(application *application.CarService) *Handler {
	return &Handler{
		application: application,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/cars", h.UpsertCar)
	mux.HandleFunc("PUT /v1/cars/{id}", h.UpsertCar)
}

type upsertCarRequest struct {
	ID    string `json:"id,omitempty"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type upsertCarResponse struct {
	Message string    `json:"message"`
	Car     model.Car `json:"car"`
}

// UpsertCar handles HTTP requests to upsert a car. (path: /cars or /cars/{id}, method: POST or PUT)
func (h *Handler) UpsertCar(w http.ResponseWriter, r *http.Request) {
	var reqBody upsertCarRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Printf("failed to decode request body: %v\n", r.Body)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.NewV4()
	if reqBody.ID != "" {
		parsedID, err := uuid.Parse(reqBody.ID)
		if err != nil {
			http.Error(w, "invalid car ID", http.StatusBadRequest)
			return
		}

		id = parsedID
	}

	car := application.UpsertCarInput{
		ID:    id,
		Brand: reqBody.Brand,
		Model: reqBody.Model,
		Year:  reqBody.Year,
	}

	addedCar, err := h.application.UpsertCar(r.Context(), car)
	if err != nil {
		http.Error(w, "failed to upsert car", http.StatusInternalServerError)
		return
	}

	respBody := &upsertCarResponse{
		Message: "car upserted successfully",
		Car:     *addedCar,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respBody)
}
