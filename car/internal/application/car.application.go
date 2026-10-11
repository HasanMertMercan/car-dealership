package application

import (
	"car-dealership/car/internal/model"
	"car-dealership/car/internal/storage"
	"context"
	"errors"
	"strings"
	"uuid"
)

type AddCarInput struct {
	Brand string
	Model string
	Year  int
}

type CarService struct {
}

func NewCarService() *CarService {
	return &CarService{}
}

func (s *CarService) AddCar(ctx context.Context, input AddCarInput) (*model.Car, error) {
	brand := strings.TrimSpace(input.Brand)
	carModel := strings.TrimSpace(input.Model)

	if brand == "" {
		return nil, errors.New("brand is required")
	}

	if carModel == "" {
		return nil, errors.New("model is required")
	}

	car := model.Car{
		ID:    uuid.NewV4(),
		Brand: brand,
		Model: carModel,
		Year:  input.Year,
	}
	err := storage.AddCar(car)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *CarService) ListCars(ctx context.Context) ([]model.Car, error) {
	return storage.Cars, nil
}
