package web_socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"results/errs"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
	Send chan Broadcast
}

type Message struct {
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.UnRegister <- c
		err := c.Conn.Close()
		if err != nil {
			fmt.Printf("%s: %s", errs.WSClientCloseFailed, err)
			return
		}
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("%s", errs.WsDecodeJsonFailed)
			continue
		}

		broadcast := Broadcast{
			Message: []byte(msg.Message),
			SendId:  msg.UserId,
		}

		hub.Broadcast <- broadcast
	}
}

func (c *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			fmt.Printf("%s: %s", errs.WSClientCloseFailed, err)
		}
	}(c.Conn)

	for msg := range c.Send {

		if err := c.Conn.WriteMessage(websocket.TextMessage, msg.Message); err != nil {
			break
		}
	}
}
