package misskey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Client はMisskey APIとの通信を担当する
type Client struct {
	Host   string
	Token  string
	Logger *slog.Logger
}

// NewClient は新しいMisskey APIクライアントを作成する
func NewClient(host, token string, logger *slog.Logger) *Client {
	return &Client{
		Host:   host,
		Token:  token,
		Logger: logger,
	}
}

// PostNote はテキストをMisskeyに投稿する
func (c *Client) PostNote(text string, visibility string, localOnly bool) error {
	endpoint := fmt.Sprintf("https://%s/api/notes/create", c.Host)

	// JSONデータの作成
	jsonData := map[string]interface{}{
		"i":         c.Token,
		"text":      text,
		"localOnly": localOnly,
	}

	if visibility != "" {
		jsonData["visibility"] = visibility
	}

	// JSONデータをバイトに変換
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("JSON data encoding error: %w", err)
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}

	// リクエストヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントの作成
	client := &http.Client{}

	// リクエストの送信
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request sending error: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスの確認
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	c.Logger.Info("Misskey API response", "status", resp.Status)
	return nil
}
