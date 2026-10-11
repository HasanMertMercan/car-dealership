package main

import (
	"fmt"
	"net"

	pb "car-dealership/car/gen/v1"

	"google.golang.org/grpc"
)

type carServiceServer struct {
	pb.UnimplementedCarServiceServer
}

func main() {
	// Create a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}
	fmt.Println("gRPC server listening on :50051")

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	// Register the CarService server with the gRPC server
	pb.RegisterCarServiceServer(grpcServer, &carServiceServer{})

	// Start serving requests on the listener
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
