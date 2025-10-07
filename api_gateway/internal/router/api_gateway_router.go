package router

import (
	"api_gateway/internal/api_storage"
	"api_gateway/internal/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"proto/authpb"
	"proto/databasepb"
	"proto/profilepb"
	"proto/scyllapb"
	"proto/smtppb"

	_ "api_gateway/docs"
)

func Router(router *gin.Engine,
	client authpb.AuthServiceClient,
	dbClient databasepb.DatabaseServiceClient,
	smtpClient smtppb.SmtpServiceClient,
	profileClient profilepb.ProfileServiceClient,
	scyllaClient scyllapb.ScyllaServiceClient) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authController := controllers.AuthController{Client: client}

			auth.POST("/register", authController.RegisterController)
			auth.POST("/login", authController.LoginController)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/send_otp", authController.SendOtpCode)
			auth.POST("/reset_password", authController.ResetPassword)
			auth.POST("/update_password", authController.UpdatePassword)
			auth.GET("/my_profile", authController.MyProfile)
		}
		profile := api.Group("/profile")
		{
			profileController := controllers.ProfileController{Client: profileClient, DbClient: dbClient}

			profile.POST("/update_nickname", profileController.UpdateNickname)
			profile.POST("/update_email", profileController.UpdateNickname)
			profile.POST("/update_bio", profileController.UpdateBio)
			profile.GET("/info", profileController.GetProfileInfo)
		}
		search := api.Group("/search")
		{
			searchController := controllers.SearchController{Client: dbClient}

			search.GET("/profiles", searchController.SearchUser)
			search.GET("/profile_by_id", searchController.SearchById)
		}
		chat := api.Group("/chat")
		{
			chatController := controllers.ChatController{AuthClient: client, ScyllaClient: scyllaClient}

			chat.GET("/get_chats", chatController.GetChatHistory)
		}
		if api_storage.Env.DebugMode {
			api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
	ping := controllers.Ping{AuthClient: client, DbClient: dbClient, SmtpClient: smtpClient, ProfileClient: profileClient, ScyllaClient: scyllaClient}
	router.GET("/api/v1/ping", ping.Ping)

	if api_storage.Env.DebugMode {
		router.GET("/ws", ping.WebSocket)
	}
}
