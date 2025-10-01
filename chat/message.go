package chat

import (
	"MessengerChat/db"
	"log"
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

	if msg.ReceiverID > 0 {
		err := database.Conn.QueryRow(
			`INSERT INTO private_messages(sender_id, receiver_id, content) 
			 VALUES($1,$2,$3) RETURNING id, created_at`,
			msg.SenderID, msg.ReceiverID, msg.Content,
		).Scan(&id, &createdAt)
		if err != nil {
			return msg, err
		}
	} else {
		err := database.Conn.QueryRow(
			`INSERT INTO messages(channel_id, sender_id, content) 
			 VALUES($1,$2,$3) RETURNING id, created_at`,
			msg.ChannelID, msg.SenderID, msg.Content,
		).Scan(&id, &createdAt)
		if err != nil {
			return msg, err
		}
	}

	// Получаем имя отправителя
	err := database.Conn.QueryRow(`SELECT username FROM users WHERE id=$1`, msg.SenderID).Scan(&msg.SenderName)
	if err != nil {
		log.Println("Cannot fetch sender name:", err)
	}

	// Получаем имя получателя, если приватное сообщение
	if msg.ReceiverID > 0 {
		err = database.Conn.QueryRow(`SELECT username FROM users WHERE id=$1`, msg.ReceiverID).Scan(&msg.ReceiverName)
		if err != nil {
			log.Println("Cannot fetch receiver name:", err)
		}
	}

	msg.ID = id
	msg.CreatedAt = createdAt
	return msg, nil
}
