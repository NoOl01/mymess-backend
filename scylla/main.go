package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "proto/scyllapb"
	"results/errs"
	"scylla_db/internal/db"
	"scylla_db/internal/grpc_server"
	"scylla_db/internal/storage"
	"syscall"
)

func main() {
	storage.LoadEnv()
	db.InitScylla([]string{fmt.Sprintf("%s", storage.Env.ScyllaHost)}, storage.Env.ScyllaSpaceKey)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.ScyllaServicePort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterScyllaServiceServer(server, &grpc_server.Server{})

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC auth_server is running on port %s... \n", storage.Env.ScyllaServicePort)
		serverErr <- server.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("%s: %v", errs.ServerError, err)
	case sig := <-quit:
		log.Printf("Received signal: %v. Shutting down.", sig)
		server.GracefulStop()
		return
	}
}
