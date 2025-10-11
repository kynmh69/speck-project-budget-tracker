.PHONY: help dev up down build test clean migrate-up migrate-down migrate-create seed logs

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## 開発環境を起動
	docker-compose up

up: ## バックグラウンドで全サービスを起動
	docker-compose up -d

down: ## 全サービスを停止
	docker-compose down

build: ## 全サービスをビルド
	docker-compose build

rebuild: ## 全サービスを再ビルドして起動
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

test-backend: ## バックエンドのテストを実行
	cd backend && go test ./...

test-frontend: ## フロントエンドのテストを実行
	cd frontend && npm test

test: test-backend test-frontend ## 全テストを実行

clean: ## Docker リソースをクリーンアップ
	docker-compose down -v
	docker system prune -f

migrate-up: ## データベースマイグレーションを実行
	docker-compose exec backend migrate -path ./migrations -database "postgres://postgres:postgres@db:5432/project_budget_tracker?sslmode=disable" up

migrate-down: ## データベースマイグレーションをロールバック
	docker-compose exec backend migrate -path ./migrations -database "postgres://postgres:postgres@db:5432/project_budget_tracker?sslmode=disable" down

migrate-create: ## 新しいマイグレーションファイルを作成 (make migrate-create name=create_users)
	@read -p "Migration name: " name; \
	cd backend && migrate create -ext sql -dir migrations -seq $$name

seed: ## 初期データを投入
	docker-compose exec backend go run scripts/seed.go

logs: ## ログを表示
	docker-compose logs -f

logs-backend: ## バックエンドのログを表示
	docker-compose logs -f backend

logs-frontend: ## フロントエンドのログを表示
	docker-compose logs -f frontend

logs-db: ## データベースのログを表示
	docker-compose logs -f db

shell-backend: ## バックエンドコンテナにシェルで入る
	docker-compose exec backend sh

shell-frontend: ## フロントエンドコンテナにシェルで入る
	docker-compose exec frontend sh

shell-db: ## データベースコンテナにシェルで入る
	docker-compose exec db psql -U postgres -d project_budget_tracker

install-backend: ## バックエンドの依存関係をインストール
	cd backend && go mod download

install-frontend: ## フロントエンドの依存関係をインストール
	cd frontend && npm install

install: install-backend install-frontend ## 全依存関係をインストール

lint-backend: ## バックエンドのlintを実行
	cd backend && go fmt ./...
	cd backend && go vet ./...

lint-frontend: ## フロントエンドのlintを実行
	cd frontend && npm run lint

lint: lint-backend lint-frontend ## 全lintを実行
