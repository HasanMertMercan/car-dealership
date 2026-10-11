package repository

import (
	"car-dealership/car/internal/model"
	"car-dealership/car/internal/storage"
)

func UpsertCar(car model.Car) error {

	// Check if the car already exists in the storage by ID and update it if found
	for i, existingCar := range storage.Cars {
		if existingCar.ID == car.ID {
			storage.Cars[i] = car
			return nil
		}
	}

	// If the car was not found in the existing storage, append the new car to the in-memory storage
	storage.Cars = append(storage.Cars, car)
	return nil
}

func ListCars() []model.Car {
	return storage.Cars
}
