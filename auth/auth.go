package main

import (
	grpcclientconnect "auth/internal/grpc_client"
	"auth/internal/storage"
	"database/routes"
	"errs"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "proto/databasepb"
	"syscall"
)

func main() {
	storage.LoadEnv()

	client, conn, err := grpcclientconnect.GrpcClientConnect()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%s: %v", errs.GrpcClientCloseFailed, err)
		}
	}(conn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.AuthPort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(server, &routes.Server{})

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC auth_server is running on port %s... \n", storage.Env.AuthPort)
		serverErr <- server.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("%s: %v", errs.ServerError, err)
	case sig := <-quit:
		log.Printf("Received signal: %v. Shutting down.", sig)
		return
	}
}
