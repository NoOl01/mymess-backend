package main

import (
	grpcclientconnect "auth/internal/grpc_client"
	"auth/internal/grpc_server"
	"auth/internal/storage"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	pb "proto/authpb"
	"results/errs"
	"syscall"
)

func main() {
	storage.LoadEnv()

	dbClient, dbConn, err := grpcclientconnect.GrpcDatabaseClientConnect()
	if err != nil {
		fmt.Println(err.Error())
	}
	smtpClient, smtpConn, err := grpcclientconnect.GrpcSmtpClientConnect()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%s: %v", errs.GrpcClientCloseFailed, err)
		}
	}(dbConn)

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%s: %v", errs.GrpcClientCloseFailed, err)
		}
	}(smtpConn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.AuthPort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, &grpc_server.Server{
		Client:     dbClient,
		SmtpClient: smtpClient,
	})

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
		server.GracefulStop()
		return
	}
}
