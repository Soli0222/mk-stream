package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"mk-stream/internal/config"
	"mk-stream/internal/websocket"
)

const (
	appName    = "mk-stream"
	appVersion = "2.0.0"
	asciiArt   = `
 __  __ _  __      _____ _                            
|  \/  | |/ /     / ____| |                          
| \  / | ' /_____| (___ | |_ _ __ ___  __ _ _ __ ___  
| |\/| |  <______\___ \| __| '__/ _ \/ _` + "`" + ` | '_ ` + "`" + ` _ \ 
| |  | | . \     ____) | |_| | |  __/ (_| | | | | | |
|_|  |_|_|\_\   |_____/ \__|_|  \___|\__,_|_| |_| |_|
                                                      
`
)

func printBanner() {
	fmt.Println(asciiArt)
	fmt.Printf("%s v%s\n\n", appName, appVersion)
}

func main() {
	// バナーの表示
	printBanner()

	// slogの設定
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	logger.Info("Starting application",
		"app", appName,
		"version", appVersion,
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 設定の読み込み
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// WebSocketManagerの作成と接続
	wsManager := websocket.NewManager(cfg.Host, cfg.Token, logger)
	if err := wsManager.Connect(); err != nil {
		logger.Error("WebSocket connection error", "error", err)
		return
	}

	// メッセージリスナーの開始
	wsManager.Listen()

	logger.Info("Connected to WebSocket server. Press Ctrl+C to exit.")

	// シグナルを待つ
	<-interrupt
	logger.Info("Received interrupt signal. Shutting down...")
	wsManager.Disconnect()
	time.Sleep(time.Second)
}
