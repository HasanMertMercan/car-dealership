package grpcserver

import (
	"context"
	"fmt"

	carv1 "car-dealership/car/gen/v1"
	"car-dealership/car/internal/application"
)

type Server struct {
	carv1.UnimplementedCarServiceServer
	application *application.CarService
}

func NewServer(application *application.CarService) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) AddCar(
	ctx context.Context,
	request *carv1.AddCarRequest,
) (*carv1.AddCarResponse, error) {
	car, err := s.application.AddCar(
		ctx,
		application.AddCarInput{
			Brand: request.GetBrand(),
			Model: request.GetModel(),
			Year:  int(request.GetYear()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &carv1.AddCarResponse{
		Message: fmt.Sprintf(
			"Car added: %s %s",
			car.Brand,
			car.Model,
		),
		Car: &carv1.Car{
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		},
	}, nil
}

func (s *Server) ListCars(ctx context.Context, request *carv1.ListCarsRequest) (*carv1.ListCarsResponse, error) {
	cars, err := s.application.ListCars(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the list of cars to the protobuf representation
	var carProtos []*carv1.Car
	for _, car := range cars {
		carProtos = append(carProtos, &carv1.Car{
			Id:    car.ID.String(),
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		})
	}

	return &carv1.ListCarsResponse{
		Cars: carProtos,
	}, nil
}
