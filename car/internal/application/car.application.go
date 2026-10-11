package application

import (
	"car-dealership/car/internal/model"
	"context"
	"errors"
	"strings"
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

func (s *CarService) AddCar(ctx context.Context, input AddCarInput) (model.Car, error) {
	brand := strings.TrimSpace(input.Brand)
	carModel := strings.TrimSpace(input.Model)

	if brand == "" {
		return model.Car{}, errors.New("brand is required")
	}

	if carModel == "" {
		return model.Car{}, errors.New("model is required")
	}

	car := model.Car{
		Brand: brand,
		Model: carModel,
		Year:  input.Year,
	}

	return car, nil
}
