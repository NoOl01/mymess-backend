package controllers

import (
	"api_gateway/internal/api_grpc_client/client"
	"api_gateway/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/authpb"
	"proto/scyllapb"
	"results/errs"
	"strconv"
	"strings"
)

type ChatController struct {
	AuthClient   authpb.AuthServiceClient
	ScyllaClient scyllapb.ScyllaServiceClient
}

// GetChatHistory
// @Summary получение списка чатов
// @Tags chat
// @Accept       json
// @Produce      json
// @Security     Token
// @Router       /chat/get_chats [get]
func (chat *ChatController) GetChatHistory(c *gin.Context) {
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

	getUserIdResp, err := client.GetUserId(chat.AuthClient, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	userId, err := strconv.ParseInt(getUserIdResp.Result, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	resp, err := client.GetChats(chat.ScyllaClient, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResult{
			Result: nil,
			Error:  strPointer(err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, models.BaseResult{
		Result: resp.Chats,
		Error:  nil,
	})
}
