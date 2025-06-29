package controllers

import (
	"api_gateway/internal/models"
	"github.com/gin-gonic/gin"
)

type Ping struct{}

// Ping
// @Summary Проверка доступности сервера и версии приложения
// @Description Публичный запрос, который клиент отправляет при запуске. Возвращает статус сервера и актуальную версию приложения.
// @Tags ping
// @Produce json
// @Router /auth/ping [get]
func (p *Ping) Ping(c *gin.Context) {

	c.JSON(200, models.PingResult{
		ApiGateway: "ok",
		Auth:       "ok",
		Database:   "ok",
		Cache:      "WIP",
		Message:    "WIP",
		Search:     "WIP",
	})
}
