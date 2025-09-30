package main

import (
	"api_gateway/internal/api_grpc_client"
	"api_gateway/internal/api_storage"
	"api_gateway/internal/router"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"proto/authpb"
	"proto/databasepb"
	"proto/profilepb"
	"proto/scyllapb"
	"proto/smtppb"
	"results/errs"
	"syscall"
)

// @title Messenger API
// @version 0.2
// @BasePath /api/v1
// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	api_storage.LoadEnv()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
		AllowHeaders:    []string{"Origin", "Content-type", "Accept", "Authorization"},
	}))

	if !api_storage.Env.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	var conns []*grpc.ClientConn

	addConn := func(conn *grpc.ClientConn) {
		conns = append(conns, conn)
	}

	authClient, authConn, err := api_grpc_client.GrpcClientConnect(
		api_grpc_client.Connect(api_storage.Env.AuthHost, api_storage.Env.AuthPort, authpb.NewAuthServiceClient),
	)
	if err != nil {
		log.Println(err.Error())
	}
	addConn(authConn)

	dbClient, dbConn, err := api_grpc_client.GrpcClientConnect(
		api_grpc_client.Connect(api_storage.Env.DbHost, api_storage.Env.DbPort, databasepb.NewDatabaseServiceClient),
	)
	if err != nil {
		log.Println(err.Error())
	}
	addConn(dbConn)

	smtpClient, smtpConn, err := api_grpc_client.GrpcClientConnect(
		api_grpc_client.Connect(api_storage.Env.DbHost, api_storage.Env.DbPort, smtppb.NewSmtpServiceClient),
	)
	if err != nil {
		log.Println(err.Error())
	}
	addConn(smtpConn)

	profileClient, profileConn, err := api_grpc_client.GrpcClientConnect(
		api_grpc_client.Connect(api_storage.Env.DbHost, api_storage.Env.DbPort, profilepb.NewProfileServiceClient),
	)
	if err != nil {
		log.Println(err.Error())
	}
	addConn(profileConn)

	scyllaClient, scyllaConn, err := api_grpc_client.GrpcClientConnect(
		api_grpc_client.Connect(api_storage.Env.ScyllaServiceHost, api_storage.Env.ScyllaServicePort, scyllapb.NewScyllaServiceClient),
	)
	if err != nil {
		log.Println(err.Error())
	}
	addConn(scyllaConn)

	defer func() {
		for _, conn := range conns {
			if err := conn.Close(); err != nil {
				log.Printf("%v", err)
			}
		}
	}()

	router.Router(r, authClient, dbClient, smtpClient, profileClient, scyllaClient)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- r.Run(fmt.Sprintf("0.0.0.0:%s", api_storage.Env.ApiPort))
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
