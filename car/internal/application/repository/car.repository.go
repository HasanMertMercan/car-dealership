package repository

import (
	"car-dealership/car/internal/application/model"
	"context"
)

type InMemoryCarRepository struct {
	cars []model.Car
}

func NewCarRepository() *InMemoryCarRepository {
	return &InMemoryCarRepository{
		cars: make([]model.Car, 0),
	}
}

func (r *InMemoryCarRepository) UpsertCar(ctx context.Context, car model.Car) error {

	// Check if the car already exists in the storage by ID and update it if found
	for i, existingCar := range r.cars {
		if existingCar.ID == car.ID {
			r.cars[i] = car
			return nil
		}
	}

	// If the car was not found in the existing storage, append the new car to the in-memory storage
	r.cars = append(r.cars, car)
	return nil
}

func (r *InMemoryCarRepository) ListCars(ctx context.Context) ([]model.Car, error) {
	return r.cars, nil
}
