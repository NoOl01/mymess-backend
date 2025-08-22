package web_socket

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Broadcast
	Register   chan *Client
	UnRegister chan *Client
}

type Broadcast struct {
	Message []byte
	ChatId  int64
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Broadcast),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.UnRegister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
