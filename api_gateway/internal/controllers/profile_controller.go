package controllers

import (
	"api_gateway/internal/api_grpc_client/client"
	"api_gateway/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/databasepb"
	"proto/profilepb"
	"results/errs"
	"strconv"
	"strings"
)

type ProfileController struct {
	Client   profilepb.ProfileServiceClient
	DbClient databasepb.DatabaseServiceClient
}

// UpdateNickname
// @Summary Обновить nickname
// @Tags profile
// @Accept       json
// @Produce      json
// @Param        input  body      models.UpdateProfile  true  "Данные для обновления профиля"
// @Security     Token
// @Router       /profile/update_nickname [post]
func (profile *ProfileController) UpdateNickname(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.MissingToken.Error()),
		})
		return
	}

	if !strings.HasPrefix(token, "Bearer") {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidToken.Error()),
		})
		return
	}
	token = strings.ReplaceAll(token, "Bearer ", "")

	body := models.UpdateProfile{}
	if err := c.ShouldBindJSON(&body); err != nil && body.Value != "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.UpdateProfile(profile.Client, body.Value, "nickname", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Result,
		Error:  nil,
	})
}

// UpdateEmail
// @Summary Обновить nickname
// @Tags profile
// @Accept       json
// @Produce      json
// @Param        input  body      models.UpdateProfile  true  "Данные для обновления профиля"
// @Security     Token
// @Router       /profile/update_email [post]
func (profile *ProfileController) UpdateEmail(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.MissingToken.Error()),
		})
		return
	}

	if !strings.HasPrefix(token, "Bearer") {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidToken.Error()),
		})
		return
	}
	token = strings.ReplaceAll(token, "Bearer ", "")

	body := models.UpdateProfile{}
	if err := c.ShouldBindJSON(&body); err != nil && body.Value != "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.UpdateProfile(profile.Client, body.Value, "nickname", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Result,
		Error:  nil,
	})
}

// UpdateBio
// @Summary Обновить bio
// @Tags profile
// @Accept       json
// @Produce      json
// @Param        input  body      models.UpdateProfile  true  "Данные для обновления профиля"
// @Security     Token
// @Router       /profile/update_bio [post]
func (profile *ProfileController) UpdateBio(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.MissingToken.Error()),
		})
		return
	}

	if !strings.HasPrefix(token, "Bearer") {
		c.JSON(http.StatusUnauthorized, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidToken.Error()),
		})
		return
	}
	token = strings.ReplaceAll(token, "Bearer ", "")

	body := models.UpdateProfile{}
	if err := c.ShouldBindJSON(&body); err != nil && body.Value != "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestBody.Error()),
		})
		return
	}

	resp, err := client.UpdateProfile(profile.Client, body.Value, "bio", token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Result,
		Error:  nil,
	})
}

// GetProfileInfo
// @Summary Получение профиля пользователя по id
// @Tags profile
// @Accept       json
// @Produce      json
// @Param id query int true "Username пользователя"
// @Router       /profile/info [get]
func (profile *ProfileController) GetProfileInfo(c *gin.Context) {
	idQuery := c.Query("id")
	if idQuery == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(fmt.Sprintf("%v: id", errs.InvalidRequestQuery)),
		})
		return
	}

	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(errs.InvalidRequestQuery.Error()),
		})
		return
	}

	resp, err := client.GetProfileInfo(profile.DbClient, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.GetBody(),
		Error:  nil,
	})
}
