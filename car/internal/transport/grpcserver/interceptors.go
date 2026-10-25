package grpcserver

import (
	carpb "car-dealership/car/gen"
	"context"
	"log"
	"time"
	"uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Implement logging logic here
	start := time.Now()
	log.Printf("gRPC request started: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("gRPC request completed: %s, Duration: %s, Error: %v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func InvalidUUIDInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Implement UUID validation logic here
	// For example, you can check if the request contains a valid UUID and return an error if it doesn't
	switch r := req.(type) {
	case *carpb.UpsertCarRequest:
		id, err := uuid.Parse(r.GetId())
		if id.String() != "" && err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid UUID: %s", r.Id)
		}
	}
	return handler(ctx, req)
}
