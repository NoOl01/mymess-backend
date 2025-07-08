package web_socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"results/errs"
	"strconv"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
	Send chan Broadcast
}

type Message struct {
	Message string `json:"chat"`
	UserId  int64  `json:"user_id"`
	SendId  int64  `json:"send_id"`
}

type KafkaJson struct {
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
	SendId  int64  `json:"send_id"`
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
			SendId:  strconv.FormatInt(msg.SendId, 10),
		}

		//payload := KafkaJson{
		//	Message: msg.Message,
		//	UserId:  msg.UserId,
		//	SendId:  msg.SendId,
		//}
		//
		//payloadBytes, err := json.Marshal(payload)
		//if err != nil{
		//	fmt.Printf("%s", errs.FailedEncodeToJson)
		//}
		//
		//err = kafka.Producer.Produce(&kf.Message{
		//	TopicPartition: kf.TopicPartition{Topic: strPtr("chat_messages"), Partition: kf.PartitionAny},
		//	Value:          payloadBytes,
		//}, nil)

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

func strPtr(s string) *string {
	return &s
}
