package router

import (
	"api_gateway/internal/api_storage"
	"api_gateway/internal/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"proto/authpb"
	"proto/databasepb"
	"proto/smtppb"

	_ "api_gateway/docs"
)

func Router(router *gin.Engine,
	client authpb.AuthServiceClient,
	dbClient databasepb.DatabaseServiceClient,
	smtpClient smtppb.SmtpServiceClient) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authController := controllers.Controller{Client: client}

			auth.POST("/register", authController.RegisterController)
			auth.POST("/login", authController.LoginController)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/send_otp", authController.SendOtpCode)
			auth.POST("/reset_password", authController.ResetPassword)
			auth.POST("/update_password", authController.UpdatePassword)
		}
		if api_storage.Env.DebugMode {
			api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
	ping := controllers.Ping{AuthClient: client, DbClient: dbClient, SmtpClient: smtpClient}
	router.GET("/api/v1/ping", ping.Ping)

	if api_storage.Env.DebugMode {
		router.GET("/ws", ping.WebSocket)
	}
}
