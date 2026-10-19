package application

import (
	"car-dealership/car/internal/application/model"
	"context"
	"errors"
	"strings"
	"uuid"
)

// CarRepository defines the interface for car data storage operations.
// It is defined within the application layer since this is the Go idiomatic way to define interfaces close to where they are used.
type CarRepository interface {
	UpsertCar(ctx context.Context, car model.Car) error
	ListCars(ctx context.Context) ([]model.Car, error)
}

type UpsertCarInput struct {
	ID    uuid.UUID
	Brand string
	Model string
	Year  int
}

// CarService provides methods to manage cars, including upserting and listing cars, and streaming real-time updates.
type CarService struct {
	listUpdate chan []model.Car // Channel to stream real-time latest state of the car list
	repository CarRepository
}

func NewCarService(repository CarRepository) *CarService {
	return &CarService{
		listUpdate: make(chan []model.Car, 1),
		repository: repository,
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
	err := s.repository.UpsertCar(ctx, car)
	if err != nil {
		return nil, err
	}

	latestCars, err := s.repository.ListCars(ctx)
	if err != nil {
		return nil, err
	}

	// Publish the latest list of cars to the update channel
	select {
	case s.listUpdate <- latestCars:
	default:
	}

	return &car, nil
}

func (s *CarService) ListCars(ctx context.Context) ([]model.Car, error) {
	latestCars, err := s.repository.ListCars(ctx)
	if err != nil {
		return nil, err
	}
	return latestCars, nil
}

// ListUpdateChannel returns a channel that streams updates to the list of cars in real-time.
func (s *CarService) ListUpdateChannel() <-chan []model.Car {
	return s.listUpdate
}
