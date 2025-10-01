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
	ChannelID int    `json:"channel_id"`
	Content   string `json:"content"`
}

func main() {
	username := "guest"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	channelID := 1 // default
	if len(os.Args) > 2 {
		if ch, err := strconv.Atoi(os.Args[2]); err == nil {
			channelID = ch
		}
	}

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws?username="+username+"&channel_id="+strconv.Itoa(channelID), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// Горутина для чтения входящих сообщений
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
			ChannelID: channelID,
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
