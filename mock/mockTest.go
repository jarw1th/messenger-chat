package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

type Message struct {
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func main() {
	username := "guest"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws?username="+username, nil)
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
			ChannelID: 1,
			Content:   text,
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
