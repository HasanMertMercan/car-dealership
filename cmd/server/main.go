package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	fmt.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
