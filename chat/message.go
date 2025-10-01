package chat

import (
	"MessengerChat/db"
	"time"
)

type Message struct {
	ID           int       `json:"id"`
	ChannelID    int       `json:"channel_id,omitempty"`
	SenderID     int       `json:"sender_id"`
	SenderName   string    `json:"sender_name,omitempty"`
	ReceiverID   int       `json:"receiver_id,omitempty"`
	ReceiverName string    `json:"receiver_name,omitempty"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

func SaveMessage(msg Message, database *db.DataBase) (Message, error) {
	var id int
	var createdAt time.Time

	var err error
	if msg.ReceiverID > 0 {
		err = database.Conn.QueryRow(
			`INSERT INTO private_messages(sender_id, receiver_id, content) 
			 VALUES($1,$2,$3) RETURNING id, created_at`,
			msg.SenderID, msg.ReceiverID, msg.Content,
		).Scan(&id, &createdAt)
	} else {
		err = database.Conn.QueryRow(
			`INSERT INTO messages(channel_id, sender_id, content) 
			 VALUES($1,$2,$3) RETURNING id, created_at`,
			msg.ChannelID, msg.SenderID, msg.Content,
		).Scan(&id, &createdAt)
	}
	if err != nil {
		return msg, err
	}

	database.Conn.QueryRow(`SELECT username FROM users WHERE id=$1`, msg.SenderID).Scan(&msg.SenderName)

	if msg.ReceiverID > 0 {
		database.Conn.QueryRow(`SELECT username FROM users WHERE id=$1`, msg.ReceiverID).Scan(&msg.ReceiverName)
	}

	msg.ID = id
	msg.CreatedAt = createdAt
	return msg, nil
}
