package repository

import (
	"car-dealership/car/internal/application/model"
)

type CarRepository struct {
	Cars []model.Car
}

func NewCarRepository() *CarRepository {
	return &CarRepository{
		Cars: TempCarStorage,
	}
}

func (r *CarRepository) UpsertCar(car model.Car) error {

	// Check if the car already exists in the storage by ID and update it if found
	for i, existingCar := range TempCarStorage {
		if existingCar.ID == car.ID {
			TempCarStorage[i] = car
			return nil
		}
	}

	// If the car was not found in the existing storage, append the new car to the in-memory storage
	TempCarStorage = append(TempCarStorage, car)
	return nil
}

func (r *CarRepository) ListCars() []model.Car {
	return TempCarStorage
}
