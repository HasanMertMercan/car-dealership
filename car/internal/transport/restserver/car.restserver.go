package restserver

import (
	"car-dealership/car/internal/application"
	"car-dealership/car/internal/application/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"uuid"
)

// Handler is the HTTP handler for the car REST API. It wraps the application layer.
type Handler struct {
	application *application.CarService
}

// NewHandler creates a new instance of the car REST API handler with the given application layer.
func NewHandler(application *application.CarService) *Handler {
	return &Handler{
		application: application,
	}
}

// RegisterRoutes registers the HTTP routes for the car REST API with the given HTTP server mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/cars", h.UpsertCar)
	mux.HandleFunc("PUT /v1/cars/{id}", h.UpsertCar)
	mux.HandleFunc("GET /v1/cars", h.ListCars)
}

// upsertCarRequest represents the request payload for upserting a car.
type upsertCarRequest struct {
	ID    string `json:"id,omitempty"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

// upsertCarResponse represents the response payload for upserting a car.
type upsertCarResponse struct {
	Message string    `json:"message"`
	Car     model.Car `json:"car"`
}

// UpsertCar handles HTTP requests to upsert a car. (path: /v1/cars or /v1/cars/{id}, method: POST or PUT)
func (h *Handler) UpsertCar(w http.ResponseWriter, r *http.Request) {
	var reqBody upsertCarRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Printf("failed to decode request body: %v\n", r.Body)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	idValue := r.PathValue("id")
	if idValue == "" {
		idValue = reqBody.ID
	}

	id := uuid.NewV4()
	if idValue != "" {
		parsedID, err := uuid.Parse(idValue)
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

func (h *Handler) ListCars(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	searchParam := queryParams.Get("search")
	yearParam := queryParams.Get("year")
	filter := &model.CarFilter{}
	if searchParam != "" {
		filter.Search = &searchParam
	}
	if yearParam != "" {
		year, err := strconv.Atoi(yearParam)
		if err == nil {
			filter.Year = &year
		} else {
			fmt.Printf("invalid year parameter: %v\n", err)
		}
	}

	cars, nextPageToken, err := h.application.ListCars(r.Context(), filter, &model.Pagination{
		PageSize:  20,
		PageToken: uuid.Nil(), // Temporarily using Nil as the page token since pagination is not yet implemented
	})
	if err != nil {
		fmt.Printf("failed to list cars: %v\n", err)
		http.Error(w, "failed to list cars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Cars          []model.Car `json:"cars"`
		NextPageToken string      `json:"next_page_token"`
	}{
		Cars:          cars,
		NextPageToken: nextPageToken.String(),
	})
}
