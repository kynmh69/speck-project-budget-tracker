# プロジェクト予算管理システム

ソフトウェア開発プロジェクトの工数予実管理と収支管理を行うWebアプリケーション

## 概要

このシステムは、複数のソフトウェア開発プロジェクトにおける以下の管理を支援します：

- **工数予実管理**: タスクごとの予定工数と実績工数を記録し、差異を可視化
- **収支管理**: プロジェクトの売上、コスト、利益を自動計算
- **グラフ可視化**: 予実比較、収支推移、タスク別工数などを多様なグラフで表示
- **ダッシュボード**: 複数プロジェクトの状況を横断的に管理

## 技術スタック

### フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **言語**: TypeScript 5.x
- **UIライブラリ**: shadcn/ui (Radix UI + Tailwind CSS)
- **グラフ**: Recharts
- **状態管理**: TanStack Query (サーバー状態) + Zustand (クライアント状態)
- **バリデーション**: Zod

### バックエンド
- **言語**: Go 1.21+
- **Webフレームワーク**: Echo v4
- **ORM**: GORM v2
- **データベース**: PostgreSQL 15+
- **認証**: JWT

## プロジェクト構造

```
.
├── backend/          # Go/Echo API サーバー
├── frontend/         # Next.js フロントエンド
├── docs/             # ドキュメント
├── .github/          # GitHub Actions CI/CD
└── docker-compose.yml # 開発環境
```

## セットアップ手順

### 前提条件

- Docker Desktop
- Go 1.21+ (ローカル開発の場合)
- Node.js 20+ (ローカル開発の場合)
- Make (オプション)

### 1. リポジトリのクローン

```bash
git clone https://github.com/your-org/speck-project-budget-tracker.git
cd speck-project-budget-tracker
```

### 2. 環境変数の設定

```bash
# Backend
cp backend/.env.example backend/.env

# Frontend
cp frontend/.env.local.example frontend/.env.local
```

必要に応じて `.env` ファイルを編集してください。

### 3. Docker Composeで起動

```bash
# 全サービスを起動（DB、Backend、Frontend）
docker-compose up -d

# ログを確認
docker-compose logs -f
```

### 4. アクセス

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **API Health Check**: http://localhost:8080/health

### 5. 停止

```bash
docker-compose down
```

## 開発

### Makefileコマンド（準備中）

```bash
# 開発環境起動
make dev

# テスト実行
make test

# ビルド
make build

# マイグレーション実行
make migrate-up

# マイグレーションロールバック
make migrate-down

# データシード投入
make seed
```

### バックエンド開発

```bash
cd backend

# 依存関係インストール
go mod download

# 開発サーバー起動（ホットリロード）
go run cmd/server/main.go

# テスト実行
go test ./...

# マイグレーション作成
migrate create -ext sql -dir migrations -seq <migration_name>
```

### フロントエンド開発

```bash
cd frontend

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# テスト実行
npm test

# E2Eテスト実行
npm run test:e2e
```

## API仕様

API仕様はOpenAPI 3.0形式で定義されています。

- 仕様書: `project-budget-tracker/.specify/contracts/project-budget-management/openapi.yaml`（準備中）
- Swagger UI: http://localhost:8080/swagger (準備中)

## テスト

### バックエンド

```bash
cd backend

# 単体テスト
go test ./internal/...

# 統合テスト
go test ./tests/integration/...

# カバレッジ
go test -cover ./...
```

### フロントエンド

```bash
cd frontend

# 単体テスト
npm test

# E2Eテスト
npm run test:e2e

# カバレッジ
npm test -- --coverage
```

## デプロイ

デプロイ手順は `docs/deployment.md` を参照してください（準備中）。

## ドキュメント

- [仕様書](project-budget-tracker/.specify/specs/project-budget-management.md)
- [実装計画](project-budget-tracker/.specify/plans/project-budget-management.md)
- [タスクリスト](project-budget-tracker/.specify/tasks/project-budget-management.md)

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 貢献

コントリビューションガイドラインは [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください（準備中）。

## コミット規約

このプロジェクトは [Conventional Commits](https://www.conventionalcommits.org/) に従います。
詳細は `project-budget-tracker/.github/instructions/commit.instructions.md` を参照してください。

## サポート

問題が発生した場合は、GitHubのIssuesで報告してください。