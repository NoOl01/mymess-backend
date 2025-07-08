package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"message/internal/jwt"
	"message/internal/storage"
	"message/internal/web_socket"
	"net/http"
	"results/errs"
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWS(hub *web_socket.Hub, w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, errs.MissingToken.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := jwt.ValidateJwt(tokenString)
	if err != nil {
		http.Error(w, errs.InvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("%s: %s\n", errs.WsUpgradeFailed, err)
		return
	}

	client := &web_socket.Client{
		Id:   userID,
		Conn: conn,
		Send: make(chan web_socket.Broadcast, 256),
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump(hub)
}

func main() {
	storage.LoadEnv()

	hub := web_socket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	fmt.Printf("Listening on :%s \n", storage.Env.MessagePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", storage.Env.MessagePort), nil); err != nil {
		log.Fatal(errs.WsListenAndServeFailed)
	}
}
