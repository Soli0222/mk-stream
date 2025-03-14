package models

// ConnectMessage はWebSocketへの接続メッセージ
type ConnectMessage struct {
	Type string      `json:"type"`
	Body ConnectBody `json:"body"`
}

// ConnectBody は接続メッセージのボディ
type ConnectBody struct {
	Channel string `json:"channel"`
	ID      string `json:"id"`
}

// DisconnectMessage はWebSocketからの切断メッセージ
type DisconnectMessage struct {
	Type string         `json:"type"`
	Body DisconnectBody `json:"body"`
}

// DisconnectBody は切断メッセージのボディ
type DisconnectBody struct {
	ID string `json:"id"`
}

// WebSocketEvent は受信するイベントタイプを定義する
const (
	EventEmojiUpdated = "emojiUpdated"
	EventEmojiAdded   = "emojiAdded"
	EventEmojiDeleted = "emojiDeleted"
)
