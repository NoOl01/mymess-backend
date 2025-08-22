package controllers

import (
	"api_gateway/internal/api_grpc_client/client"
	"api_gateway/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/databasepb"
	"results/errs"
)

type SearchController struct {
	Client databasepb.DatabaseServiceClient
}

// SearchUser
// @Summary Поиск пользователей по username
// @Tags search
// @Accept       json
// @Produce      json
// @Param username query string true "Username пользователя"
// @Router       /search/profiles [get]
func (search *SearchController) SearchUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.BaseResult{
			Result: nil,
			Error:  strPointer(fmt.Sprintf("%v: username", errs.InvalidRequestQuery)),
		})
		return
	}

	resp, err := client.FindUser(search.Client, username)
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
