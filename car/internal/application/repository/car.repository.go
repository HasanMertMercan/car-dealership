package repository

import (
	"car-dealership/car/internal/application/model"
	"context"
)

type InMemoryCarRepository struct {
	Cars []model.Car
}

func NewCarRepository() *InMemoryCarRepository {
	return &InMemoryCarRepository{
		Cars: make([]model.Car, 0),
	}
}

func (r *InMemoryCarRepository) UpsertCar(ctx context.Context, car model.Car) error {

	// Check if the car already exists in the storage by ID and update it if found
	for i, existingCar := range r.Cars {
		if existingCar.ID == car.ID {
			r.Cars[i] = car
			return nil
		}
	}

	// If the car was not found in the existing storage, append the new car to the in-memory storage
	r.Cars = append(r.Cars, car)
	return nil
}

func (r *InMemoryCarRepository) ListCars(ctx context.Context) ([]model.Car, error) {
	return r.Cars, nil
}
