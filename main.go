package main

import (
	"log"
	"net/http"

	"MessengerChat/chat"
	"MessengerChat/db"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	connStr := "host=127.0.0.1 port=5432 user=postgres dbname=chatdb sslmode=disable"
	database, err := db.NewDataBase(connStr)
	if err != nil {
		panic(err)
	}

	hub := chat.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			username = "guest"
		}

		userID, err := db.EnsureUser(database, username)
		if err != nil {
			log.Println("EnsureUser failed:", err)
			conn.Close()
			return
		}

		client := chat.NewClient(hub, conn, 1, userID)
		hub.Register <- client

		client.SendHistory(database, 50)

		go client.WritePump()
		go client.ReadPump(func(msg chat.Message) (chat.Message, error) {
			return chat.SaveMessage(msg, database)
		})
	})

	log.Println("WebSocket chat server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
