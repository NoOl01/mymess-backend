package main

import (
	"database/routes"
	"log"
	"net"
	pb "proto/databasepb"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(s, &routes.Server{})

	log.Println("gRPC grpc_server is running on port 50051...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
