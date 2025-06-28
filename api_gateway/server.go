package main

import (
	"api_gateway/internal/api_grpc_client"
	"api_gateway/internal/api_storage"
	"api_gateway/internal/router"
	"errs"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Messenger API
// @version 0.1
// @host localhost:8080
// @BasePath /api/v1
func main() {
	api_storage.LoadEnv()

	r := gin.Default()
	r.Use(cors.Default())

	if !api_storage.Env.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	client, err := api_grpc_client.GrpcClientConnect()
	if err != nil {
		fmt.Println(err)
		return
	}

	router.Router(r, client)

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
