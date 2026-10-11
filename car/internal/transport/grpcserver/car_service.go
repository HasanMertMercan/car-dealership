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
