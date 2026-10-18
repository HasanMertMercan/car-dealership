package application

import (
	"car-dealership/car/internal/application/model"
	"car-dealership/car/internal/application/repository"
	"context"
	"errors"
	"strings"
	"uuid"
)

type UpsertCarInput struct {
	ID    uuid.UUID
	Brand string
	Model string
	Year  int
}

// CarService provides methods to manage cars, including upserting and listing cars, and streaming real-time updates.
type CarService struct {
	listUpdate chan []model.Car // Channel to stream real-time latest state of the car list
}

func NewCarService() *CarService {
	return &CarService{
		listUpdate: make(chan []model.Car, 1),
	}
}

func (s *CarService) UpsertCar(ctx context.Context, input UpsertCarInput) (*model.Car, error) {
	brand := strings.TrimSpace(input.Brand)
	carModel := strings.TrimSpace(input.Model)

	if brand == "" {
		return nil, errors.New("brand is required")
	}

	if carModel == "" {
		return nil, errors.New("model is required")
	}

	// Generate a new UUID if the ID is not provided
	if input.ID == uuid.Nil() {
		input.ID = uuid.NewV4()
	}

	car := model.Car{
		ID:    input.ID,
		Brand: brand,
		Model: carModel,
		Year:  input.Year,
	}
	err := repository.UpsertCar(car)
	if err != nil {
		return nil, err
	}

	// Publish the latest list of cars to the update channel
	select {
	case s.listUpdate <- repository.ListCars():
	default:
	}

	return &car, nil
}

func (s *CarService) ListCars(ctx context.Context) ([]model.Car, error) {
	return repository.ListCars(), nil
}

// ListUpdateChannel returns a channel that streams updates to the list of cars in real-time.
func (s *CarService) ListUpdateChannel() <-chan []model.Car {
	return s.listUpdate
}
