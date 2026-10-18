package repository

import (
	"car-dealership/car/internal/application/model"
)

func UpsertCar(car model.Car) error {

	// Check if the car already exists in the storage by ID and update it if found
	for i, existingCar := range Cars {
		if existingCar.ID == car.ID {
			Cars[i] = car
			return nil
		}
	}

	// If the car was not found in the existing storage, append the new car to the in-memory storage
	Cars = append(Cars, car)
	return nil
}

func ListCars() []model.Car {
	return Cars
}
