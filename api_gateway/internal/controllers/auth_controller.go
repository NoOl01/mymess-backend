package controllers

import (
	"api_gateway/internal/api_grpc_client"
	"api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/authpb"
	"results/errs"
)

type Controller struct {
	Client authpb.AuthServiceClient
}

func strPointer(s string) *string {
	return &s
}

// RegisterController
// @Summary Регистрация
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.Register  true  "Данные регистрации"
// @Router /auth/register [post]
func (auth *Controller) RegisterController(c *gin.Context) {
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	resp, err := api_grpc_client.RegisterRequest(auth.Client, body.Username, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.AuthResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResult{
		Result: gin.H{
			"access_token": resp.AccessToken},
		Error: nil,
	})
}

// LoginController
// @Summary Авторизация
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.Login  true  "Данные для входа"
// @Router /auth/login [post]
func (auth *Controller) LoginController(c *gin.Context) {
	var body models.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	if (body.Username == nil || body.Email == nil) && body.Password == "" {
		c.JSON(http.StatusBadRequest, models.AuthResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := api_grpc_client.LoginRequest(auth.Client, *body.Username, *body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.AuthResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResult{
		Result: gin.H{
			"access_token": resp.AccessToken},
		Error: nil,
	})
}
