package chat

type Hub struct {
	Channels   map[int]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Channels:   make(map[int]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Channels[client.ChannelID] == nil {
				h.Channels[client.ChannelID] = make(map[*Client]bool)
			}
			h.Channels[client.ChannelID][client] = true

		case client := <-h.Unregister:
			if _, ok := h.Channels[client.ChannelID][client]; ok {
				delete(h.Channels[client.ChannelID], client)
				close(client.Send)
			}

		case msg := <-h.Broadcast:
			if msg.ReceiverID > 0 {
				for _, clients := range h.Channels {
					for c := range clients {
						if c.UserID == msg.ReceiverID || c.UserID == msg.SenderID {
							select {
							case c.Send <- msg:
							default:
								close(c.Send)
								delete(clients, c)
							}
						}
					}
				}
			} else { // канал
				for c := range h.Channels[msg.ChannelID] {
					select {
					case c.Send <- msg:
					default:
						close(c.Send)
						delete(h.Channels[msg.ChannelID], c)
					}
				}
			}
		}
	}
}
