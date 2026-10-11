package main

import (
	"fmt"
	"net"

	carv1 "car-dealership/car/gen/v1"

	"google.golang.org/grpc"

	"car-dealership/car/internal/application"
	"car-dealership/car/internal/transport/grpcserver"
)

type carServiceServer struct {
	carv1.UnimplementedCarServiceServer
	application *application.CarService
}

func main() {
	// Create a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}
	fmt.Println("gRPC server listening on :50051")

	// Initialize the application and gRPC server handler
	carApplication := application.NewCarService()
	carHandler := grpcserver.NewServer(carApplication)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the CarService server with the gRPC server
	carv1.RegisterCarServiceServer(
		grpcServer,
		carHandler,
	)

	// Start serving requests on the listener
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
