package repository

import (
	"car-dealership/car/internal/model"
	"car-dealership/car/internal/storage"
)

func AddCar(car model.Car) error {
	// Append the new car to the in-memory storage
	storage.Cars = append(storage.Cars, car)
	return nil
}

func ListCars() []model.Car {
	return storage.Cars
}
