package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

type Message struct {
	ChannelID  int    `json:"channel_id"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	Content    string `json:"content"`
}

func main() {
	username := "guest"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	channelName := "General"
	if len(os.Args) > 2 {
		channelName = os.Args[2]
	}

	receiverID := 0
	if len(os.Args) > 3 {
		if r, err := strconv.Atoi(os.Args[3]); err == nil {
			receiverID = r
		}
	}

	url := "ws://localhost:8080/ws?username=" + username + "&channel_name=" + channelName
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Получено: %s", message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		msg := Message{
			ChannelID:  0,
			ReceiverID: receiverID,
			Content:    text,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("json marshal:", err)
			continue
		}
		err = c.WriteMessage(websocket.TextMessage, jsonMsg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("scanner error:", err)
	}
}
