package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "proto/smtppb"
	"results/errs"
	"smtp/internal/grpc_server"
	"smtp/internal/storage"
	"syscall"
)

func main() {
	storage.LoadEnv()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.SmtpServicePort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterSmtpServiceServer(server, &grpc_server.Server{})

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC auth_server is running on port %s... \n", storage.Env.SmtpServicePort)
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
