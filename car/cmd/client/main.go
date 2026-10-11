package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "car-dealership/car/gen/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new CarService client using the connection.
	carServiceClient := pb.NewCarServiceClient(conn)

	// Contact the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add a new car to the server (Temporarily named "Toyota Corolla")
	r, err := carServiceClient.AddCar(ctx, &pb.CreateCar{
		Brand: "Toyota",
		Model: "Corolla",
	})
	if err != nil {
		log.Fatalf("could not add car: %v", err)
	}

	// Print the response from the server after adding the new car.
	fmt.Printf("Cars: %v\n", r.GetMessage())
}
