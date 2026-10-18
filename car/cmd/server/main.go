package main

import (
	"fmt"
	"net"
	"net/http"

	carpb "car-dealership/car/gen"

	"google.golang.org/grpc"

	"car-dealership/car/internal/application"
	"car-dealership/car/internal/transport/grpcserver"
	"car-dealership/car/internal/transport/restserver"
)

type carServiceServer struct {
	carpb.UnimplementedCarServiceServer
	application *application.CarService
}

func main() {
	// Initialize the application
	carApplication := application.NewCarService()

	// Initialize the gRPC server handler with the application instance
	carGRPCHandler := grpcserver.NewServer(carApplication)

	// Initialize the REST server handler with the application instance
	carRESTHandler := restserver.NewHandler(carApplication)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the CarService server with the gRPC server
	carpb.RegisterCarServiceServer(
		grpcServer,
		carGRPCHandler,
	)

	// Register REST routes with the HTTP server mux
	mux := http.NewServeMux()
	carRESTHandler.RegisterRoutes(mux)

	// Create the HTTP server for the REST API
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("HTTP server for REST API listening on :8080")

	// Start the HTTP server for the REST API in a separate goroutine to allow the gRPC server to run concurrently
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			// Handle the error.
			fmt.Printf("failed to serve HTTP: %v\n", err)
		}
	}()

	// Create a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}
	fmt.Println("gRPC server for GRPC API listening on :50051")

	// Start serving requests on the listener
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
