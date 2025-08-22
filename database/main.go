package main

import (
	"database/internal/db_connect"
	"database/internal/grpc_server"
	"database/internal/storage"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	pb "proto/databasepb"
	"results/errs"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	storage.LoadEnv()

	db := db_connect.Connect()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.DbServicePort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(server, &grpc_server.Server{
		Db: db,
	})

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("gRPC auth_server is running on port %s... \n", storage.Env.DbServicePort)
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
