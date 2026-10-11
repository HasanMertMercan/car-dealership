package storage

import (
	"car-dealership/car/internal/model"
)

func AddCar(car model.Car) error {
	// Append the new car to the in-memory storage
	Cars = append(Cars, car)
	return nil
}
