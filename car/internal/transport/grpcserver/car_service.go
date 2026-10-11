package grpcserver

import (
	"context"
	"fmt"
	"uuid"

	carpb "car-dealership/car/gen"
	"car-dealership/car/internal/application"
)

type Server struct {
	carpb.UnimplementedCarServiceServer
	application *application.CarService
}

func NewServer(application *application.CarService) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) UpsertCar(
	ctx context.Context,
	request *carpb.UpsertCarRequest,
) (*carpb.UpsertCarResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		id = uuid.NewV4()
	}

	car, err := s.application.UpsertCar(
		ctx,
		application.UpsertCarInput{
			ID:    id,
			Brand: request.GetBrand(),
			Model: request.GetModel(),
			Year:  int(request.GetYear()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &carpb.UpsertCarResponse{
		Message: fmt.Sprintf(
			"Car added: %s %s",
			car.Brand,
			car.Model,
		),
		Car: &carpb.Car{
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		},
	}, nil
}

func (s *Server) ListCars(ctx context.Context, request *carpb.ListCarsRequest) (*carpb.ListCarsResponse, error) {
	cars, err := s.application.ListCars(ctx)
	if err != nil {
		return nil, err
	}

	// Convert the list of cars to the protobuf representation
	var carProtos []*carpb.Car
	for _, car := range cars {
		carProtos = append(carProtos, &carpb.Car{
			Id:    car.ID.String(),
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		})
	}

	return &carpb.ListCarsResponse{
		Cars: carProtos,
	}, nil
}
