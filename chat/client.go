package chat

import (
	"MessengerChat/db"
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

func (c *Client) SendHistoryWithPrivate(database *db.DataBase, limit int) {
	rows, err := database.Conn.Query(
		`SELECT id, channel_id, sender_id, content, created_at 
         FROM messages 
         WHERE channel_id=$1 
         ORDER BY created_at DESC 
         LIMIT $2`,
		c.ChannelID, limit,
	)
	if err != nil {
		log.Println("SendHistory query error:", err)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			log.Println("SendHistory scan error:", err)
			continue
		}
		messages = append([]Message{msg}, messages...)
	}

	privRows, err := database.Conn.Query(
		`SELECT id, sender_id, receiver_id, content, created_at
         FROM private_messages
         WHERE sender_id=$1 OR receiver_id=$1
         ORDER BY created_at DESC
         LIMIT $2`,
		c.UserID, limit,
	)
	if err != nil {
		log.Println("SendHistory private query error:", err)
		return
	}
	defer privRows.Close()

	for privRows.Next() {
		var msg Message
		if err := privRows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt); err != nil {
			log.Println("SendHistory private scan error:", err)
			continue
		}
		messages = append([]Message{msg}, messages...)
	}

	for _, msg := range messages {
		c.Send <- msg
	}
}
