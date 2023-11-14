package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"mk-stream/handlers"
)

type ConnectMessage struct {
	Type string      `json:"type"`
	Body ConnectBody `json:"body"`
}

type ConnectBody struct {
	Channel string `json:"channel"`
	ID      string `json:"id"`
}

type DisconnectMessage struct {
	Type string         `json:"type"`
	Body DisconnectBody `json:"body"`
}

type DisconnectBody struct {
	ID string `json:"id"`
}

func connectWebSocket(host, token string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/streaming", RawQuery: "i=" + token}
	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	return conn, err
}

func WebSocketMessage(jsonData []byte) {
	// Parse JSON data using map[string]interface{}
	var message map[string]interface{}
	if err := json.Unmarshal(jsonData, &message); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Check the "type" field
	messageType, ok := message["type"].(string)
	if !ok {
		fmt.Println("Invalid message format: type field is missing or not a string")
		return
	}

	// Check the "body" field
	body, ok := message["body"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid message format: body field is missing or not a map")
		return
	}
	fmt.Println(body)

	// Get the necessary fields
	switch messageType {
	case "emojiUpdated":
		emojiData, ok := body["emojis"].([]interface{})
		if !ok {
			fmt.Println("Invalid emojiUpdated message format: emojis field is missing or not an array")
			return
		}

		// Perform processing for emojis here
		fmt.Println("Received emojiUpdated message:", emojiData)
		handlers.NoteEmojiUpdated(emojiData)

	case "emojiAdded":
		emojiData, ok := body["emoji"].(map[string]interface{})
		if !ok {
			fmt.Println("Invalid emojiAdded message format: emojis field is missing or not an array")
			return
		}

		// Perform processing for emojis here
		fmt.Println("Received emojiAdded message:", emojiData)
		handlers.NoteEmojiAdded(emojiData)

	case "emojiDeleted":
		emojiData, ok := body["emojis"].([]interface{})
		if !ok {
			fmt.Println("Invalid emojiDeleted message format: emojis field is missing or not an array")
			return
		}

		// Perform processing for emojis here
		fmt.Println("Received emojiDeleted message:", emojiData)
		handlers.NoteEmojiDeleted(emojiData)

	default:
		fmt.Println("Unknown message type:", messageType)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")

	connectID := uuid.New().String()

	// Establish WebSocket connection
	conn, err := connectWebSocket(host, token)
	if err != nil {
		fmt.Println("WebSocket connection error:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to asynchronously receive messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Message reading error:", err)
				// If an error occurs, attempt to reconnect
				conn, err = connectWebSocket(host, token)
				if err != nil {
					fmt.Println("Reconnection error:", err)
					return
				}
				continue
			}
			fmt.Printf("Received message: %s\n", message)
			WebSocketMessage(message)
		}
	}()

	// Send JSON data to the WebSocket
	connectData := ConnectMessage{
		Type: "connect",
		Body: ConnectBody{
			Channel: "main",
			ID:      connectID,
		},
	}
	sendData, err := json.Marshal(connectData)
	if err != nil {
		fmt.Println("JSON data encoding error:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, sendData)
	if err != nil {
		fmt.Println("Message sending error:", err)
		return
	}

	// Send messages interactively
	for {
		select {
		case <-interrupt:
			// Send JSON data to the WebSocket (disconnect)
			disconnectData := DisconnectMessage{
				Type: "disconnect",
				Body: DisconnectBody{
					ID: connectID,
				},
			}
			sendData, err := json.Marshal(disconnectData)
			if err != nil {
				fmt.Println("JSON data encoding error:", err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, sendData)
			if err != nil {
				fmt.Println("Message sending error:", err)
				return
			}

			fmt.Println("Exiting...")
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Close message sending error:", err)
			}
			time.Sleep(time.Second)
			return
		}
	}
}
