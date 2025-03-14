# MK-Stream

![mk-stream](https://img.shields.io/badge/version-2.0.0-blue)

```
 __  __ _  __      _____ _                            
|  \/  | |/ /     / ____| |                          
| \  / | ' /_____| (___ | |_ _ __ ___  __ _ _ __ ___  
| |\/| |  <______\___ \| __| '__/ _ \/ _` | '_ ` _ \ 
| |  | | . \     ____) | |_| | |  __/ (_| | | | | | |
|_|  |_|_|\_\   |_____/ \__|_|  \___|\__,_|_| |_| |_|
```

MK-Streamは、Misskeyサーバーのストリーミング APIを監視し、特定のイベントが発生した際に自動的にノートを投稿するためのアプリケーションです。サーバー管理者向けのツールとして、サーバーの状態変化やユーザーアクティビティなどの情報を自動的に共有することができます。

## 機能

現在は以下の機能が実装されています：

- **カスタム絵文字の監視**
  - 絵文字が追加されたとき → 通知ノートを投稿
  - 絵文字が更新されたとき → 通知ノートを投稿
  - 絵文字が削除されたとき → 通知ノートを投稿

## 動作環境

- Go 1.24.1以上

## セットアップ

### ローカル環境での実行

1. リポジトリをクローンします：
   ```bash
   git clone https://github.com/Soli0222/mk-stream.git
   cd mk-stream
   ```

2. 依存パッケージをインストールします：
   ```bash
   go mod download
   ```

3. 環境変数を設定します：
   ```bash
   cp .env.example .env
   ```
   
   .envファイルを編集して、以下の設定を行います：
   ```
   HOST=your-misskey-server.tld  # MisskeyサーバーのホストURL（例：mi.example.com）
   TOKEN=your_token_here         # MisskeyのAPIトークン
   ```

4. アプリケーションを実行します：
   ```bash
   go run main.go
   ```

### Docker を使った実行

1. .envファイルを作成します：
   ```bash
   cp .env.example .env
   ```
   
   内容を適切に編集します。

2. Docker Compose でアプリケーションを起動します：
   ```bash
   docker-compose up -d
   ```

### Kubernetes でのデプロイ（Helm）

1. Helm チャートをカスタマイズします：
   ```bash
   cd helm
   ```
   
   values.yaml を必要に応じて編集します。

2. シークレットを設定します：
   ```
   # secret.yaml を編集して API トークンを設定
   ```

3. Helm チャートをデプロイします：
   ```bash
   helm install mk-stream . -f values.yaml
   ```

## ビルド方法

### ローカルビルド

```bash
go build -o mk-stream main.go
```

### Docker イメージのビルド

```bash
docker build -t mk-stream:latest .
```

## 設定オプション

| 環境変数 | 説明 | デフォルト値 |
|----------|------|-------------|
| `HOST` | MisskeyサーバーのホストURL | なし（必須） |
| `TOKEN` | MisskeyのAPIトークン | なし（必須） |

## 開発

### プロジェクト構造

```
mk-stream/
├── main.go             # エントリーポイント
├── internal/
│   ├── config/         # 設定関連の処理
│   ├── handlers/       # イベントハンドラ
│   ├── misskey/        # Misskey API 通信
│   ├── models/         # データモデル
│   └── websocket/      # WebSocket 通信
├── helm/               # Kubernetes デプロイ用の Helm チャート
└── docker-compose.yml  # Docker Compose 設定
```

### 貢献方法

1. このリポジトリをフォークします
2. 機能追加やバグ修正を行う新しいブランチを作成します
3. 変更をコミットします
4. リポジトリにプッシュします
5. プルリクエストを作成します

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。詳細は LICENSE ファイルをご覧ください。

## 作者

- [@Soli0222](https://github.com/Soli0222)

---

バグ報告や機能リクエストは、GitHubの[Issues](https://github.com/Soli0222/mk-stream/issues)にてお願いします。