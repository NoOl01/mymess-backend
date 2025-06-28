package router

import (
	"api_gateway/internal/api_storage"
	controllers2 "api_gateway/internal/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"proto/authpb"

	_ "api_gateway/docs"
)

func Router(router *gin.Engine, client authpb.AuthServiceClient) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authController := controllers2.Controller{Client: client}

			auth.POST("/register", authController.RegisterController)
			auth.POST("/login", authController.LoginController)
		}
		if api_storage.Env.DebugMode {
			api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		api.GET("/ping", controllers2.Ping{}.Ping)
	}
}
