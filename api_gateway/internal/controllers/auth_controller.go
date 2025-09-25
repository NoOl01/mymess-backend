package controllers

import (
	"api_gateway/internal/api_grpc_client/client"
	"api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/authpb"
	"results/errs"
)

type AuthController struct {
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
func (auth *AuthController) RegisterController(c *gin.Context) {
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	resp, err := client.RegisterRequest(auth.Client, body.Nickname, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
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
func (auth *AuthController) LoginController(c *gin.Context) {
	var body models.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	if (body.Username == nil || body.Email == nil) && body.Password == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.LoginRequest(auth.Client, *body.Username, *body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: gin.H{
			"access_token": resp.AccessToken},
		Error: nil,
	})
}

// RefreshToken
// @Summary Refresh
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.Refresh  true  "Обновление токена"
// @Router /auth/refresh [post]
func (auth *AuthController) RefreshToken(c *gin.Context) {
	var body models.Refresh

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	if body.AccessToken == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.SmtpCodeOrEmailMissing.Error()),
		})
		return
	}

	resp, err := client.RefreshRequest(auth.Client, body.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: gin.H{
			"access_token": resp.AccessToken},
		Error: nil,
	})
}

// SendOtpCode
// @Summary Отправить код на почту
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.SendOtp  true  "Данные для отправки кода на почту"
// @Router /auth/send_otp [post]
func (auth *AuthController) SendOtpCode(c *gin.Context) {
	var body models.SendOtp

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: errs.InvalidRequestBody.Error(),
		})
		return
	}

	if body.Email == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: errs.InvalidRequestBody.Error(),
		})
		return
	}

	resp, err := client.SendOtp(auth.Client, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Result,
	})
}

// ResetPassword
// @Summary Подтвердить OTP код
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.ResetPassword  true  "Данные для проверки кода"
// @Router /auth/reset_password [post]
func (auth *AuthController) ResetPassword(c *gin.Context) {
	var body models.ResetPassword

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	if body.Email == "" || body.Code == 0 {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.ResetPassword(auth.Client, body.Email, body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: gin.H{
			"status":      resp.Result,
			"reset_token": resp.ResetToken,
		},
		Error: nil,
	})
}

// UpdatePassword
// @Summary Обновление пароля
// @Tags auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.UpdatePassword  true  "Данные для обновления пароля"
// @Router /auth/update_password [post]
func (auth *AuthController) UpdatePassword(c *gin.Context) {
	var body models.UpdatePassword

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	if body.Email == "" || body.Password == "" || body.ResetToken == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: errs.InvalidRequestBody.Error(),
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.UpdatePassword(auth.Client, body.Email, body.Password, body.ResetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: err.Error(),
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Result,
		Error:  nil,
	})
}

// MyProfile
// @Summary Обновление пароля
// @Tags auth
// @Accept       json
// @Produce      json
// @Security     Token
// @Router /auth/my_profile [get]
func (auth *AuthController) MyProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.MissingToken.Error()),
		})
		return
	}

	resp, err := client.MyProfile(auth.Client, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp,
		Error:  nil,
	})
}
