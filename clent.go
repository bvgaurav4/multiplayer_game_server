package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://localhost:8080/play"

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	// Start a goroutine to read messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			fmt.Println("Server:", string(message))
		}
	}()

	// Read input from the user and send messages
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		scanner.Scan()
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Closing connection...")
			break
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
