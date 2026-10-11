package application

import (
	"car-dealership/car/internal/model"
	"car-dealership/car/internal/repository"
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

type CarService struct {
}

func NewCarService() *CarService {
	return &CarService{}
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
	return &car, nil
}

func (s *CarService) ListCars(ctx context.Context) ([]model.Car, error) {
	return repository.ListCars(), nil
}
