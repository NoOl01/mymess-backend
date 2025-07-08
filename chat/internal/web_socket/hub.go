package web_socket

type Hub struct {
	Clients     map[*Client]bool
	ClientsById map[string]*Client
	Broadcast   chan Broadcast
	Register    chan *Client
	UnRegister  chan *Client
}

type Broadcast struct {
	Message []byte
	SendId  string
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		ClientsById: make(map[string]*Client),
		Broadcast:   make(chan Broadcast),
		Register:    make(chan *Client),
		UnRegister:  make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			h.ClientsById[client.Id] = client
		case client := <-h.UnRegister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientsById, client.Id)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			if client, ok := h.ClientsById[message.SendId]; ok {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					delete(h.ClientsById, client.Id)
				}
			}
		}
	}
}
