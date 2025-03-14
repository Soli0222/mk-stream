package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"mk-stream/internal/handlers"
	"mk-stream/internal/models"
)

// Manager はWebSocket接続を管理する
type Manager struct {
	host      string
	token     string
	conn      *websocket.Conn
	connectID string
	done      chan struct{}
	logger    *slog.Logger
}

// NewManager は新しいManagerを作成する
func NewManager(host, token string, logger *slog.Logger) *Manager {
	return &Manager{
		host:      host,
		token:     token,
		connectID: uuid.New().String(),
		done:      make(chan struct{}),
		logger:    logger,
	}
}

// Connect はWebSocket接続を確立する
func (m *Manager) Connect() error {
	conn, err := connectWebSocket(m.host)
	if err != nil {
		return err
	}
	m.conn = conn

	// Send connect message
	connectData := models.ConnectMessage{
		Type: "connect",
		Body: models.ConnectBody{
			Channel: "main",
			ID:      m.connectID,
		},
	}
	sendData, err := json.Marshal(connectData)
	if err != nil {
		return fmt.Errorf("JSON data encoding error: %w", err)
	}

	err = m.conn.WriteMessage(websocket.TextMessage, sendData)
	if err != nil {
		return fmt.Errorf("message sending error: %w", err)
	}

	return nil
}

// Listen はメッセージの受信を開始する
func (m *Manager) Listen() {
	go func() {
		for {
			select {
			case <-m.done:
				return
			default:
				_, message, err := m.conn.ReadMessage()
				if err != nil {
					m.logger.Error("Message reading error", "error", err)
					// If an error occurs, attempt to reconnect
					for i := 0; i < 5; i++ { // 最大5回再試行
						m.logger.Info("Attempting to reconnect", "attempt", i+1, "max", 5)
						if err := m.Connect(); err == nil {
							m.logger.Info("Successfully reconnected")
							break
						}
						time.Sleep(5 * time.Second) // 5秒待ってから再試行
					}
					continue
				}
				m.logger.Info("Received message", "message", string(message))
				m.handleMessage(message)
			}
		}
	}()
}

// Disconnect はWebSocket接続を切断する
func (m *Manager) Disconnect() {
	close(m.done)

	// Send disconnect message
	disconnectData := models.DisconnectMessage{
		Type: "disconnect",
		Body: models.DisconnectBody{
			ID: m.connectID,
		},
	}
	sendData, err := json.Marshal(disconnectData)
	if err != nil {
		m.logger.Error("JSON data encoding error", "error", err)
		return
	}

	err = m.conn.WriteMessage(websocket.TextMessage, sendData)
	if err != nil {
		m.logger.Error("Message sending error", "error", err)
	}

	err = m.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		m.logger.Error("Close message sending error", "error", err)
	}

	m.conn.Close()
}

// connectWebSocket はWebSocket接続を確立する
func connectWebSocket(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/streaming"}
	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial(u.String(), nil)
	return conn, err
}

// handleMessage はWebSocketメッセージを処理する
func (m *Manager) handleMessage(jsonData []byte) {
	var message map[string]interface{}
	if err := json.Unmarshal(jsonData, &message); err != nil {
		m.logger.Error("Error decoding JSON", "error", err)
		return
	}

	// Check the "type" field
	messageType, ok := message["type"].(string)
	if !ok {
		m.logger.Error("Invalid message format", "reason", "type field is missing or not a string")
		return
	}

	// Check the "body" field
	body, ok := message["body"].(map[string]interface{})
	if !ok {
		m.logger.Error("Invalid message format", "reason", "body field is missing or not a map")
		return
	}

	// Get the necessary fields
	switch messageType {
	case models.EventEmojiUpdated:
		m.handleEmojiUpdated(body)
	case models.EventEmojiAdded:
		m.handleEmojiAdded(body)
	case models.EventEmojiDeleted:
		m.handleEmojiDeleted(body)
	default:
		m.logger.Info("Unknown message type", "type", messageType)
	}
}

// handleEmojiUpdated はemoji更新メッセージを処理する
func (m *Manager) handleEmojiUpdated(body map[string]interface{}) {
	emojiData, ok := body["emojis"].([]interface{})
	if !ok {
		m.logger.Error("Invalid emojiUpdated message format", "reason", "emojis field is missing or not an array")
		return
	}

	m.logger.Info("Received emojiUpdated message")
	handlers.HandleEmojiUpdated(emojiData, m.logger)
}

// handleEmojiAdded はemoji追加メッセージを処理する
func (m *Manager) handleEmojiAdded(body map[string]interface{}) {
	emojiData, ok := body["emoji"].(map[string]interface{})
	if !ok {
		m.logger.Error("Invalid emojiAdded message format", "reason", "emoji field is missing or not a map")
		return
	}

	m.logger.Info("Received emojiAdded message")
	handlers.HandleEmojiAdded(emojiData, m.logger)
}

// handleEmojiDeleted はemoji削除メッセージを処理する
func (m *Manager) handleEmojiDeleted(body map[string]interface{}) {
	emojiData, ok := body["emojis"].([]interface{})
	if !ok {
		m.logger.Error("Invalid emojiDeleted message format", "reason", "emojis field is missing or not an array")
		return
	}

	m.logger.Info("Received emojiDeleted message")
	handlers.HandleEmojiDeleted(emojiData, m.logger)
}
