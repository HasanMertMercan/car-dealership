package grpcserver

import (
	"context"
	"fmt"
	"uuid"

	carpb "car-dealership/car/gen"
	"car-dealership/car/internal/application"

	"google.golang.org/grpc"
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

// This method streams updates to the list of cars to all subscribers in real-time.
// It only returns error state, otherwise it keeps streaming updates indefinitely.
func (s *Server) StreamCarListUpdates(request *carpb.StreamCarListUpdatesRequest, stream grpc.ServerStreamingServer[carpb.StreamCarListUpdatesResponse]) error {

	// Subscribe to the application's list update channel.
	appUpdates := s.application.ListUpdateChannel()
	for { // Infinite loop to continuously stream updates
		select {
		case cars := <-appUpdates: // Receive the latest list of cars from the application's update channel
			protoCars := make([]*carpb.Car, 0)
			for _, car := range cars { // Convert each car to its protobuf representation
				protoCar := &carpb.Car{
					Id:    car.ID.String(),
					Brand: car.Brand,
					Model: car.Model,
					Year:  int32(car.Year),
				}
				protoCars = append(protoCars, protoCar)
			}
			// Send the converted list of cars to the client via the stream
			stream.Send(&carpb.StreamCarListUpdatesResponse{
				Cars: protoCars,
			})
		// Continue to listen for further updates
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}
