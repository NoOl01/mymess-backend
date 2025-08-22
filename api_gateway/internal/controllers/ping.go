package controllers

import (
	"api_gateway/internal/api_grpc_client"
	"github.com/gin-gonic/gin"
	"net/http"
	"proto/authpb"
	"proto/databasepb"
	"proto/profilepb"
	"proto/scyllapb"
	"proto/smtppb"
	"sync"
)

type Ping struct {
	AuthClient    authpb.AuthServiceClient
	DbClient      databasepb.DatabaseServiceClient
	SmtpClient    smtppb.SmtpServiceClient
	ProfileClient profilepb.ProfileServiceClient
	ScyllaClient  scyllapb.ScyllaServiceClient
}

type Result struct {
	Auth     chan string
	Database chan string
	Smtp     chan string
	Profile  chan string
	Scylla   chan string
}

// Ping
// @Summary Проверка доступности сервера и версии приложения
// @Description Публичный запрос, который клиент отправляет при запуске. Возвращает статус сервисов.
// @Tags ping
// @Produce json
// @Router /ping [get]
func (p *Ping) Ping(c *gin.Context) {

	wg := sync.WaitGroup{}
	wg.Add(5)

	res := Result{
		Auth:     make(chan string, 1),
		Database: make(chan string, 1),
		Smtp:     make(chan string, 1),
	}

	go func() {
		defer wg.Done()
		resp, err := api_grpc_client.AuthPing(p.AuthClient)
		if err != nil {
			res.Auth <- err.Error()
			return
		}
		res.Auth <- resp.Result
	}()

	go func() {
		defer wg.Done()
		resp, err := api_grpc_client.DatabasePing(p.DbClient)
		if err != nil {
			res.Database <- err.Error()
			return
		}
		res.Database <- resp.Result
	}()

	go func() {
		defer wg.Done()
		resp, err := api_grpc_client.SmtpPing(p.SmtpClient)
		if err != nil {
			res.Smtp <- err.Error()
			return
		}
		res.Smtp <- resp.Result
	}()

	go func() {
		defer wg.Done()
		resp, err := api_grpc_client.ProfilePing(p.ProfileClient)
		if err != nil {
			res.Smtp <- err.Error()
			return
		}
		res.Profile <- resp.Result
	}()

	go func() {
		defer wg.Done()
		resp, err := api_grpc_client.ScyllaPing(p.ScyllaClient)
		if err != nil {
			res.Smtp <- err.Error()
			return
		}
		res.Scylla <- resp.Result
	}()

	wg.Wait()

	close(res.Auth)
	close(res.Database)
	close(res.Smtp)

	c.JSON(http.StatusOK, gin.H{
		"api_gateway": "ok",
		"auth":        <-res.Auth,
		"database":    <-res.Database,
		"smtp":        <-res.Smtp,
		"profile":     <-res.Profile,
		"scylla":      <-res.Scylla,
	})
}

// WebSocket
// @Summary WebSocket connect
// @Tags web_socket
// @Produce plain
// @Param token query string true "JWT access token"
// @Success 101 {string} string "Switching Protocols"
// @Router /ws [get]
func (p *Ping) WebSocket(c *gin.Context) {
	c.JSON(http.StatusSwitchingProtocols, gin.H{
		"result": "this is an example of connecting to a web socket",
	})
}
