package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	grpcclientconnect "profile/internal/grpc_client"
	"profile/internal/grpc_server"
	"profile/internal/storage"
	pb "proto/profilepb"
	"results/errs"
	"syscall"
)

func main() {
	storage.LoadEnv()

	dbClient, dbConn, err := grpcclientconnect.GrpcDatabaseClientConnect()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%s: %v", errs.GrpcClientCloseFailed, err)
		}
	}(dbConn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", storage.Env.ProfilePort))
	if err != nil {
		log.Fatalf("%s: %v", errs.FailedListen, err)
	}

	server := grpc.NewServer()
	pb.RegisterProfileServiceServer(server, &grpc_server.Server{
		Client: dbClient,
	})

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC auth_server is running on port %s... \n", storage.Env.ProfilePort)
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
