package misskey

import (
	"log/slog"
	"os"
)

var defaultClient *Client

// InitDefaultClient はデフォルトのMisskey APIクライアントを初期化する
func InitDefaultClient(logger *slog.Logger) {
	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")

	defaultClient = NewClient(host, token, logger)
}

// PostLocalNote はローカル向けノートを投稿する
func PostLocalNote(text string, logger *slog.Logger) error {
	if defaultClient == nil {
		InitDefaultClient(logger)
	}

	return defaultClient.PostNote(text, "public", true)
}

// PostHomeNote はホーム向けノートを投稿する
func PostHomeNote(text string, logger *slog.Logger) error {
	if defaultClient == nil {
		InitDefaultClient(logger)
	}

	return defaultClient.PostNote(text, "home", false)
}

// PostFollowersNote はフォロワー向けノートを投稿する
func PostFollowersNote(text string, logger *slog.Logger) error {
	if defaultClient == nil {
		InitDefaultClient(logger)
	}

	return defaultClient.PostNote(text, "followers", false)
}
