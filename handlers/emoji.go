package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func NoteEmojiAdded(emojiData map[string]interface{}) {
	category := getStringValue(emojiData, "category")
	license := getStringValue(emojiData, "license")
	localOnly := getBoolStringValue(emojiData, "localOnly")
	name := getStringValue(emojiData, "name")

	text := "<center>$[sparkle 🎉] カスタム絵文字が追加されました $[sparkle 🎉]\n\n:" + name + ":</center>\n名前: `" + name + "` \n\nカテゴリー: `" + category + "`\nライセンス: `" + license + "`\nローカルのみ: `" + localOnly + "`"
	note(text)
}

func NoteEmojiUpdated(emojiData []interface{}) {
	var emojiMap map[string]interface{}
	for _, emoji := range emojiData {
		emojiMap = emoji.(map[string]interface{})
	}

	category := getStringValue(emojiMap, "category")
	license := getStringValue(emojiMap, "license")
	localOnly := getBoolStringValue(emojiMap, "localOnly")
	name := getStringValue(emojiMap, "name")

	text := "<center>$[jelly 🔄] カスタム絵文字が更新されました $[jelly 🔄]\n\n:" + name + ":</center>\n名前: `" + name + "` \n\nカテゴリー: `" + category + "`\nライセンス: `" + license + "`\nローカルのみ: `" + localOnly + "`"
	note(text)
}

func NoteEmojiDeleted(emojiData []interface{}) {
	var emojiMap map[string]interface{}
	for _, emoji := range emojiData {
		emojiMap = emoji.(map[string]interface{})
	}

	category := getStringValue(emojiMap, "category")
	name := getStringValue(emojiMap, "name")

	text := "<center>$[spin.y 🗑️] カスタム絵文字が削除されました $[spin.y 🗑️]\n\n</center>\n名前: `" + name + "` \n\nカテゴリー: `" + category + "`"
	note(text)
}

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

func note(text string) {
	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")

	endpoint := "https://" + host + "/api/notes/create"

	// JSONデータの作成
	jsonData := map[string]interface{}{
		"i":         token,
		"text":      text,
		"localOnly": true,
	}

	// JSONデータをバイトに変換
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("JSONデータの変換エラー:", err)
		return
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("リクエストの作成エラー:", err)
		return
	}

	// リクエストヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントの作成
	client := &http.Client{}

	// リクエストの送信
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("リクエストの送信エラー:", err)
		return
	}
	defer resp.Body.Close()

	// レスポンスの表示
	fmt.Println("レスポンスコード:", resp.Status)
}
