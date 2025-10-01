package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan Message
	UserID    int
	ChannelID int
}

func NewClient(hub *Hub, conn *websocket.Conn, channelID int, userID int) *Client {
	return &Client{
		Hub:       hub,
		Conn:      conn,
		Send:      make(chan Message, 256),
		UserID:    userID,
		ChannelID: channelID,
	}
}

func (c *Client) ReadPump(saveFunc func(Message) (Message, error)) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		msg.SenderID = c.UserID
		savedMsg, err := saveFunc(msg)
		if err != nil {
			log.Println("SaveMessage error:", err)
			continue
		}
		c.Hub.Broadcast <- savedMsg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
