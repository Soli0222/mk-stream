package handlers

import (
	"log/slog"
	"strconv"

	"mk-stream/internal/misskey"
)

// HandleEmojiAdded は絵文字が追加された時の処理を行う
func HandleEmojiAdded(emojiData map[string]interface{}, logger *slog.Logger) {
	category := getStringValue(emojiData, "category")
	license := getStringValue(emojiData, "license")
	localOnly := getBoolStringValue(emojiData, "localOnly")
	name := getStringValue(emojiData, "name")

	text := formatEmojiAddedText(name, category, license, localOnly)

	err := misskey.PostLocalNote(text, logger)
	if err != nil {
		logger.Error("Failed to post note", "error", err)
	}
}

// HandleEmojiUpdated は絵文字が更新された時の処理を行う
func HandleEmojiUpdated(emojiData []interface{}, logger *slog.Logger) {
	if len(emojiData) == 0 {
		logger.Warn("No emoji data found")
		return
	}

	emojiMap, ok := emojiData[0].(map[string]interface{})
	if !ok {
		logger.Error("Invalid emoji data format")
		return
	}

	category := getStringValue(emojiMap, "category")
	license := getStringValue(emojiMap, "license")
	localOnly := getBoolStringValue(emojiMap, "localOnly")
	name := getStringValue(emojiMap, "name")

	text := formatEmojiUpdatedText(name, category, license, localOnly)

	err := misskey.PostLocalNote(text, logger)
	if err != nil {
		logger.Error("Failed to post note", "error", err)
	}
}

// HandleEmojiDeleted は絵文字が削除された時の処理を行う
func HandleEmojiDeleted(emojiData []interface{}, logger *slog.Logger) {
	if len(emojiData) == 0 {
		logger.Warn("No emoji data found")
		return
	}

	emojiMap, ok := emojiData[0].(map[string]interface{})
	if !ok {
		logger.Error("Invalid emoji data format")
		return
	}

	category := getStringValue(emojiMap, "category")
	name := getStringValue(emojiMap, "name")

	text := formatEmojiDeletedText(name, category)

	err := misskey.PostLocalNote(text, logger)
	if err != nil {
		logger.Error("Failed to post note", "error", err)
	}
}

// formatEmojiAddedText は絵文字追加の通知テキストをフォーマットする
func formatEmojiAddedText(name, category, license, localOnly string) string {
	return "<center>$[sparkle 🎉] カスタム絵文字が追加されました $[sparkle 🎉]\n\n:" + name +
		":</center>\n\n名前: `" + name + "` \nカテゴリー: `" + category +
		"`\nライセンス: `" + license + "`\nローカルのみ: `" + localOnly + "`"
}

// formatEmojiUpdatedText は絵文字更新の通知テキストをフォーマットする
func formatEmojiUpdatedText(name, category, license, localOnly string) string {
	return "<center>$[jelly 🔄] カスタム絵文字が更新されました $[jelly 🔄]\n\n:" + name +
		":</center>\n\n名前: `" + name + "` \nカテゴリー: `" + category +
		"`\nライセンス: `" + license + "`\nローカルのみ: `" + localOnly + "`"
}

// formatEmojiDeletedText は絵文字削除の通知テキストをフォーマットする
func formatEmojiDeletedText(name, category string) string {
	return "<center>$[spin.y 🗑️] カスタム絵文字が削除されました $[spin.y 🗑️]\n\n</center>\n\n名前: `" +
		name + "` \nカテゴリー: `" + category + "`"
}

// getStringValue はmapから文字列値を安全に取得する
func getStringValue(data map[string]interface{}, key string) string {
	value, ok := data[key]
	if !ok || value == nil {
		return "nil"
	}

	strValue, ok := value.(string)
	if !ok {
		return "nil"
	}

	return strValue
}

// getBoolStringValue はmapからブール値を文字列として安全に取得する
func getBoolStringValue(data map[string]interface{}, key string) string {
	value, ok := data[key]
	if !ok || value == nil {
		return "nil"
	}

	boolValue, ok := value.(bool)
	if !ok {
		return "nil"
	}

	return strconv.FormatBool(boolValue)
}
