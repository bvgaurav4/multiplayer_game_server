package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (modify for security)
	},
}
var (
	client  []*websocket.Conn
	clients map[string]*websocket.Conn
)

type Message struct {
	Category string `json:"category"`        // "server", "game", "event"
	Event    string `json:"event,omitempty"` // "player_joined", "game_start"
	Content  string `json:"content"`
	PlayerID string `json:"playerID,omitempty"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	client = append(client, conn)
	fmt.Println("Client connected", client)
	welcome := "hello"
	err = conn.WriteJSON(welcome)
	// Read messages from the client
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			for i := 0; i < len(client); i++ {
				if client[i] == conn {
					fmt.Println("yup found it")
				}
				delete(clients, "somthing")
			}
			break
		}

		fmt.Printf("Received: %s\n %s", msg)
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	port := "8080"
	fmt.Println(client)
	fmt.Println("WebSocket server started on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
