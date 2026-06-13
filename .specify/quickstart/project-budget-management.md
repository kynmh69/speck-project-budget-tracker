# Quickstart & Validation Guide: プロジェクト予算管理システム

**Feature**: project-budget-management
**Date**: 2026-06-12
**Plan**: [plans/project-budget-management.md](../plans/project-budget-management.md)
**Spec**: [specs/project-budget-management.md](../specs/project-budget-management.md)

MVP（User Story 1〜3 / Phase 3〜5実装分）が end-to-end で動作することを検証するためのガイド。
API契約の詳細は [contracts/project-budget-management/openapi.yaml](../contracts/project-budget-management/openapi.yaml)、
スキーマは [data-models/project-budget-management.md](../data-models/project-budget-management.md) を参照。

## 前提条件

- Docker Desktop（PostgreSQL 15 / Backend / Frontend をコンテナで起動）
- Make（任意。なければ `docker-compose` コマンドを直接実行）
- ローカル実行する場合: Go 1.21+、Node.js 20+

## セットアップ

```bash
# 環境変数の準備
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local

# 全サービス起動（db: 5432, backend: 8080, frontend: 3000）
make up          # = docker-compose up -d

# マイグレーション適用
make migrate-up

# 起動確認
curl http://localhost:8080/health
# => {"status":"ok","message":"Project Budget Tracker API is running"}
```

> **Note**: 初期データシード（`make seed` / `backend/scripts/seed.go`）は Phase 6 タスク T159 で作成予定。現時点では下記のAPI経由でデータを投入する。

## 自動テスト

```bash
# Backend: 単体 + 統合テスト（tests/unit/service, tests/integration）
make test-backend        # = cd backend && go test ./...

# Backend: lint
make lint-backend        # = go fmt + go vet

# Frontend: lint / 型チェック
make lint-frontend       # = npm run lint
cd frontend && npm run type-check
```

> **Note**: Frontend の単体テスト（Jest）・E2Eテスト（Playwright）は未セットアップ（Phase 6 タスク T154 ほかで導入予定）。現時点で `make test-frontend` は失敗する。

**期待結果**: `go test ./...` が全パス（budget / project / task の統合テストと service 単体テストを含む）。

## APIスモークテスト（curl）

### 1. 認証（前提: ユーザー登録 → ログイン → トークン取得）

```bash
# 登録
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Demo User","email":"demo@example.com","password":"password123"}'

# ログインしてJWTを取得
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"demo@example.com","password":"password123"}' | jq -r '.token // .data.token')
```

**期待結果**: 登録・ログインとも 2xx。`/api/v1/projects` 等にトークンなしでアクセスすると 401。

### 2. User Story 1 — プロジェクト管理

```bash
# 作成
curl -s -X POST http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"検証プロジェクト","description":"quickstart検証用","status":"active"}'

# 一覧（ページネーション・フィルタ・検索）
curl -s "http://localhost:8080/api/v1/projects?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 詳細 / 更新 / 削除（論理削除）
# GET/PUT/DELETE /api/v1/projects/{id}
```

**期待結果**: 作成したプロジェクトが一覧に表示され、削除後は一覧から消える（論理削除）。

### 3. User Story 2 — 工数予実管理

```bash
# タスク作成（予定工数つき）
curl -s -X POST http://localhost:8080/api/v1/projects/$PROJECT_ID/tasks \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"設計タスク","planned_hours":40}'

# 実績工数を更新 → PUT /api/v1/tasks/{id} {"actual_hours":48}

# 予実サマリー
curl -s http://localhost:8080/api/v1/projects/$PROJECT_ID/summary \
  -H "Authorization: Bearer $TOKEN"
```

**期待結果**: サマリーに合計予定工数・合計実績工数・差異が含まれ、実績 48 > 予定 40 の超過が差異として算出される。

### 4. User Story 3 — 収支管理

```bash
# メンバー登録（時間単価つき） → POST /api/v1/members {"name":"...","hourly_rate":5000}
# プロジェクトへ割り当て   → POST /api/v1/projects/{id}/members
# 工数記録               → POST /api/v1/time-entries
# 売上登録               → PUT  /api/v1/projects/{id}/budget/revenue {"revenue":1000000}

# 収支サマリー
curl -s http://localhost:8080/api/v1/projects/$PROJECT_ID/budget \
  -H "Authorization: Bearer $TOKEN"
```

**期待結果**: コスト = Σ(工数 × 単価) が自動計算され、利益 = 売上 − コスト、利益率が返却される。売上 < コストのとき赤字として扱われる。

## UI検証シナリオ（http://localhost:3000）

1. `/register` → `/login` でログインできる（未認証で保護ページに入れないこと）
2. `/projects` でプロジェクト作成 → 一覧・検索・フィルタが動作する
3. プロジェクト詳細のタスクタブで予定・実績工数を入力 → 差異が表示され、超過分が警告色になる
4. 収支タブで売上を入力 → コスト・利益・利益率が表示され、赤字時に警告表示される

## パフォーマンス確認（Phase 6 / T155）

- API応答: 平均 < 200ms、p95 < 500ms（k6 / Apache Bench 等）
- セキュリティ: `cd frontend && npm audit`、`cd backend && go mod verify`（T156）

## トラブルシューティング

- ログ確認: `make logs`（個別: `make logs-backend` / `make logs-frontend` / `make logs-db`）
- DBに直接入る: `make shell-db`
- 環境再構築: `make rebuild`、ボリュームごと削除は `make clean`
