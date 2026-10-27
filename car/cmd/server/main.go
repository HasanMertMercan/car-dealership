package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	carpb "car-dealership/car/gen"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"car-dealership/car/internal/application"
	"car-dealership/car/internal/application/repository"
	"car-dealership/car/internal/transport/grpcserver"
	"car-dealership/car/internal/transport/restserver"
)

func main() {
	fmt.Println("Starting Car Service Server...")
	databaseURL := os.Getenv("CAR_DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("CAR_DATABASE_URL is not set")
		return
	}

	startupContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(startupContext, databaseURL)
	if err != nil {
		fmt.Printf("failed to create pgx pool: %v\n", err)
		return
	}
	defer pool.Close()

	if err := pool.Ping(startupContext); err != nil {
		fmt.Printf("failed to ping database: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to the database")

	// 1. Initialize the application
	// carRepository := inmemoryrepository.NewCarRepository(pool) // Legacy in-memory repository for testing
	carRepository := repository.NewCarRepository(pool)
	carApplication := application.NewHandler(carRepository)

	// 2.a. Initialize the gRPC server handler with the application instance
	carGRPCHandler := grpcserver.NewServer(carApplication)

	// Create a new gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor( // Chain multiple unary interceptors for gRPC requests
			grpcserver.LoggingInterceptor,
			grpcserver.InvalidUUIDInterceptor,
		), // Add logging and UUID validation interceptors for gRPC requests (implemented in interceptors.go)
	)

	// Register the CarService server with the gRPC server
	carpb.RegisterCarServiceServer(
		grpcServer,
		carGRPCHandler,
	)

	// 2.b. Initialize the REST server handler with the application instance
	carRESTHandler := restserver.NewHandler(carApplication)

	// Register REST routes with the HTTP server mux
	mux := http.NewServeMux()
	carRESTHandler.RegisterRoutes(mux)

	// Create the HTTP server for the REST API
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("HTTP server for REST API listening on :8080")

	// 3.a. Start the HTTP server for the REST API in a separate goroutine to allow the gRPC server to run concurrently
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

	// 3.b. Start serving requests on the listener
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
