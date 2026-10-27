package grpcserver

import (
	"context"
	"fmt"
	"uuid"

	carpb "car-dealership/car/gen"
	"car-dealership/car/internal/application"
	"car-dealership/car/internal/application/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	carpb.UnimplementedCarServiceServer
	application *application.CarServiceHandler
}

func NewServer(application *application.CarServiceHandler) *Server {
	return &Server{
		application: application,
	}
}

func (s *Server) UpsertCar(
	ctx context.Context,
	request *carpb.UpsertCarRequest,
) (*carpb.UpsertCarResponse, error) {
	pathID := request.GetId()
	id, err := uuid.Parse(pathID)
	if err != nil && pathID != "" {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID")
	}
	if pathID == "" {
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
			Id:    car.ID.String(),
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		},
	}, nil
}

func (s *Server) ListCars(ctx context.Context, request *carpb.ListCarsRequest) (*carpb.ListCarsResponse, error) {
	filter := &model.CarFilter{}
	if request.Search != nil {
		filter.Search = request.Search
	}
	if request.Year != nil {
		year := int(request.GetYear())
		filter.Year = &year
	}

	pagination := &model.Pagination{}
	if request.PageSize <= 0 || request.PageSize > 100 {
		pagination.PageSize = 20 // Default page size if invalid
	} else {
		pagination.PageSize = request.PageSize
	}

	pageToken := request.GetPageToken()
	if pageToken != "" {
		pt, err := uuid.Parse(pageToken)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid page token")
		}
		pagination.PageToken = pt
	}

	cars, nextPageToken, err := s.application.ListCars(ctx, filter, pagination)
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

	token := ""
	if nextPageToken != uuid.Nil() {
		token = nextPageToken.String()
	}

	return &carpb.ListCarsResponse{
		Cars:          carProtos,
		NextPageToken: token,
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
