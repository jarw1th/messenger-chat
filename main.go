package main

import (
	"log"
	"net/http"

	"MessengerChat/chat"
	"MessengerChat/db"
	"MessengerChat/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	connStr := utils.GetDBConnStr()
	database, err := db.NewDataBase(connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	hub := chat.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, hub, database)
	})

	log.Println("WebSocket chat server started on :8080")
	port := utils.GetPort()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, hub *chat.Hub, database *db.DataBase) {
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

	channelName := r.URL.Query().Get("channel_name")
	if channelName == "" {
		channelName = "General"
	}
	channelID, err := db.EnsureChannel(database, channelName)
	if err != nil {
		log.Println("EnsureChannel failed:", err)
		conn.Close()
		return
	}

	client := chat.NewClient(hub, conn, channelID, userID)
	hub.Register <- client

	showAll := (userID == 0)
	client.SendHistoryWithPrivate(database, 50, showAll)

	go client.WritePump()
	go client.ReadPump(func(msg chat.Message) (chat.Message, error) {
		return chat.SaveMessage(msg, database)
	})
}
