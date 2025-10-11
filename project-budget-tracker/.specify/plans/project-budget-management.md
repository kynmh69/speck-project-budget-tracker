# Implementation Plan: プロジェクト予算管理システム

**Branch**: `feature/project-budget-management` | **Date**: 2025-10-11 | **Spec**: [project-budget-management.md](../specs/project-budget-management.md)  
**Input**: Feature specification from `.specify/specs/project-budget-management.md`

## Summary

ソフトウェア開発プロジェクトの工数予実管理と収支管理を行うWebアプリケーションを構築する。フロントエンドはNext.js + shadcn/ui、バックエンドはGolang (Echo + Gorm) で構成し、複数プロジェクトの管理とグラフによる可視化機能を提供する。段階的に実装し、P1機能（プロジェクト管理、工数予実、収支計算）をMVPとして最初に完成させる。

## Technical Context

**Language/Version**:
- Frontend: TypeScript 5.x / Node.js 20.x / Next.js 14.x (App Router)
- Backend: Go 1.21+

**Primary Dependencies**:
- Frontend:
  - Next.js 14 (React 18+, App Router)
  - shadcn/ui (Radix UI + Tailwind CSS)
  - Recharts (グラフ可視化)
  - TanStack Query (React Query v5 - データフェッチング)
  - Zustand (クライアント状態管理)
  - Zod (バリデーション)
  - date-fns (日付処理)
- Backend:
  - Echo v4 (Webフレームワーク)
  - GORM v2 (ORM)
  - JWT-go (認証)
  - validator/v10 (入力バリデーション)
  - golang-migrate (マイグレーション)

**Storage**: PostgreSQL 15+ (開発環境はDocker Compose、本番はRDS等を想定)

**Testing**:
- Frontend: Jest + React Testing Library (単体)、Playwright (E2E)
- Backend: Go testing package + testify (単体)、httptest (APIテスト)

**Target Platform**: 
- Frontend: Webブラウザ (Chrome, Firefox, Safari, Edge 最新2バージョン)
- Backend: Linux サーバー (Docker コンテナ)

**Project Type**: Web application (frontend + backend)

**Performance Goals**:
- API レスポンス: 平均 < 200ms (p95 < 500ms)
- ダッシュボード初回表示: < 2秒 (50プロジェクト)
- グラフ描画: < 3秒 (100データポイント)
- 同時接続ユーザー: 100+ (初期想定)

**Constraints**:
- データベース接続数: 最大20コネクション (初期)
- メモリ使用量: バックエンド < 512MB、フロントエンド初回ロード < 300KB (JS)
- セキュリティ: HTTPS必須、JWT有効期限24時間、XSS/CSRF対策必須

**Scale/Scope**:
- 想定ユーザー数: 10-100名 (組織内利用)
- プロジェクト数: 最大500件 (1組織あたり)
- タスク数: 最大10,000件 (1プロジェクトあたり500件)
- データ保持期間: 5年 (アーカイブ含む)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Simplicity Gates

✅ **Single Responsibility**: 各コンポーネント・サービスは単一の責務を持つ設計
✅ **Minimal Dependencies**: 必要最小限の依存関係（shadcn/uiは実質的にコピー&ペーストでロックインなし）
✅ **Clear Contracts**: REST API設計、明確なスキーマ定義
✅ **Standard Patterns**: Next.js App Router、Echo標準ルーティング、GORM標準パターン

### Potential Violations (要監視)

⚠️ **Frontend + Backend 分離**: モノレポではなく別ディレクトリ構成
- **Why Needed**: フロントエンドとバックエンドで言語・フレームワークが異なるため、明確な分離が必要
- **Mitigation**: 共通のDocker Compose環境、明確なAPI契約、型共有の仕組み（OpenAPI Spec等）

⚠️ **複数の状態管理層**: サーバー状態(TanStack Query) + クライアント状態(Zustand)
- **Why Needed**: サーバーデータのキャッシング・同期とUIステートの管理は異なる性質
- **Mitigation**: 責務を明確に分離、サーバー状態はTanStack Queryに一元化

## Project Structure

### Documentation (this feature)

```
.specify/
├── specs/
│   └── project-budget-management.md      # Feature specification
├── plans/
│   └── project-budget-management.md      # This file (implementation plan)
├── research/
│   └── project-budget-management.md      # Technical research (Phase 0)
├── data-models/
│   └── project-budget-management.md      # Database schema design (Phase 1)
├── contracts/
│   └── project-budget-management/
│       ├── openapi.yaml                  # API contract (Phase 1)
│       └── types.ts                      # Shared TypeScript types
└── tasks/
    └── project-budget-management.md      # Task breakdown (Phase 2)
```

### Source Code (repository root)

```
project-budget-tracker/
├── backend/                              # Golang Echo API
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                   # エントリーポイント
│   ├── internal/
│   │   ├── config/                       # 設定管理
│   │   │   └── config.go
│   │   ├── models/                       # GORMモデル
│   │   │   ├── project.go
│   │   │   ├── task.go
│   │   │   ├── member.go
│   │   │   ├── time_entry.go
│   │   │   └── project_member.go
│   │   ├── repository/                   # データアクセス層
│   │   │   ├── project_repository.go
│   │   │   ├── task_repository.go
│   │   │   └── member_repository.go
│   │   ├── service/                      # ビジネスロジック
│   │   │   ├── project_service.go
│   │   │   ├── task_service.go
│   │   │   ├── budget_service.go
│   │   │   └── auth_service.go
│   │   ├── handler/                      # HTTPハンドラー
│   │   │   ├── project_handler.go
│   │   │   ├── task_handler.go
│   │   │   ├── member_handler.go
│   │   │   ├── budget_handler.go
│   │   │   └── auth_handler.go
│   │   ├── middleware/                   # ミドルウェア
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   ├── dto/                          # データ転送オブジェクト
│   │   │   ├── project_dto.go
│   │   │   ├── task_dto.go
│   │   │   └── response.go
│   │   └── validator/                    # カスタムバリデーター
│   │       └── custom_validators.go
│   ├── migrations/                       # DBマイグレーション
│   │   ├── 000001_create_users.up.sql
│   │   ├── 000001_create_users.down.sql
│   │   ├── 000002_create_projects.up.sql
│   │   └── ...
│   ├── tests/                            # テスト
│   │   ├── integration/
│   │   │   ├── project_test.go
│   │   │   └── task_test.go
│   │   └── unit/
│   │       ├── service/
│   │       └── repository/
│   ├── scripts/                          # ビルド・デプロイスクリプト
│   │   └── seed.go                       # 初期データ投入
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .env.example
│
├── frontend/                             # Next.js + shadcn/ui
│   ├── src/
│   │   ├── app/                          # App Router
│   │   │   ├── (auth)/                  # 認証グループ
│   │   │   │   ├── login/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── register/
│   │   │   │       └── page.tsx
│   │   │   ├── (dashboard)/             # ダッシュボードグループ
│   │   │   │   ├── layout.tsx           # 共通レイアウト
│   │   │   │   ├── page.tsx             # ダッシュボード
│   │   │   │   ├── projects/            # プロジェクト一覧
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/           # プロジェクト詳細
│   │   │   │   │   │   ├── page.tsx
│   │   │   │   │   │   ├── tasks/      # タスク管理
│   │   │   │   │   │   ├── budget/     # 収支管理
│   │   │   │   │   │   └── analytics/  # グラフ表示
│   │   │   │   │   └── new/            # 新規作成
│   │   │   │   │       └── page.tsx
│   │   │   │   └── members/             # メンバー管理
│   │   │   │       └── page.tsx
│   │   │   ├── api/                     # API Routes (プロキシ用)
│   │   │   │   └── [...proxy]/
│   │   │   │       └── route.ts
│   │   │   ├── layout.tsx               # ルートレイアウト
│   │   │   └── globals.css
│   │   ├── components/                   # React コンポーネント
│   │   │   ├── ui/                      # shadcn/ui components
│   │   │   │   ├── button.tsx
│   │   │   │   ├── card.tsx
│   │   │   │   ├── dialog.tsx
│   │   │   │   ├── form.tsx
│   │   │   │   ├── input.tsx
│   │   │   │   ├── select.tsx
│   │   │   │   ├── table.tsx
│   │   │   │   └── ...
│   │   │   ├── projects/                # プロジェクト関連
│   │   │   │   ├── project-list.tsx
│   │   │   │   ├── project-card.tsx
│   │   │   │   ├── project-form.tsx
│   │   │   │   └── project-filters.tsx
│   │   │   ├── tasks/                   # タスク関連
│   │   │   │   ├── task-list.tsx
│   │   │   │   ├── task-form.tsx
│   │   │   │   └── task-item.tsx
│   │   │   ├── budget/                  # 収支関連
│   │   │   │   ├── budget-summary.tsx
│   │   │   │   └── cost-breakdown.tsx
│   │   │   ├── charts/                  # グラフコンポーネント
│   │   │   │   ├── plan-actual-chart.tsx
│   │   │   │   ├── budget-chart.tsx
│   │   │   │   ├── trend-chart.tsx
│   │   │   │   └── pie-chart.tsx
│   │   │   ├── layout/                  # レイアウト
│   │   │   │   ├── header.tsx
│   │   │   │   ├── sidebar.tsx
│   │   │   │   └── navbar.tsx
│   │   │   └── common/                  # 共通コンポーネント
│   │   │       ├── loading.tsx
│   │   │       ├── error-boundary.tsx
│   │   │       └── pagination.tsx
│   │   ├── lib/                         # ユーティリティ
│   │   │   ├── api-client.ts            # API クライアント
│   │   │   ├── utils.ts                 # shadcn/ui utils
│   │   │   ├── auth.ts                  # 認証ヘルパー
│   │   │   └── date-utils.ts            # 日付ユーティリティ
│   │   ├── hooks/                       # カスタムフック
│   │   │   ├── use-projects.ts          # TanStack Query hooks
│   │   │   ├── use-tasks.ts
│   │   │   ├── use-members.ts
│   │   │   └── use-auth.ts
│   │   ├── store/                       # Zustand store
│   │   │   ├── auth-store.ts
│   │   │   └── ui-store.ts
│   │   ├── types/                       # TypeScript型定義
│   │   │   ├── project.ts
│   │   │   ├── task.ts
│   │   │   ├── member.ts
│   │   │   ├── budget.ts
│   │   │   └── api.ts
│   │   └── schemas/                     # Zod スキーマ
│   │       ├── project-schema.ts
│   │       ├── task-schema.ts
│   │       └── member-schema.ts
│   ├── public/                          # 静的ファイル
│   │   └── images/
│   ├── tests/                           # テスト
│   │   ├── e2e/                         # Playwright E2E
│   │   │   ├── projects.spec.ts
│   │   │   └── tasks.spec.ts
│   │   └── unit/                        # Jest 単体テスト
│   │       └── components/
│   ├── .env.local.example
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── Dockerfile
│
├── docker-compose.yml                    # 開発環境
├── docker-compose.prod.yml               # 本番環境
├── .gitignore
└── README.md
```

**Structure Decision**: Web application構成（Option 2）を選択。フロントエンドとバックエンドを明確に分離し、各々が独立してビルド・デプロイ可能な構成とする。Docker Composeで開発環境を統合管理し、本番環境では各コンテナを個別にスケール可能にする。

## Complexity Tracking

*現時点で憲法違反はなし。シンプルな構成を維持*

## Phase 0: Research & Technology Validation

### 0.1 技術スタックの検証

**Goal**: 選定した技術スタックの組み合わせが要件を満たすことを確認

**Tasks**:
1. Next.js 14 App Router + shadcn/ui のセットアップと基本動作確認
2. Echo + GORM のセットアップとCRUD操作の動作確認
3. PostgreSQL との接続確認
4. JWT認証フローの実装パターン確認
5. Recharts によるグラフ描画のパフォーマンス検証
6. TanStack Query のキャッシング戦略確認

**Deliverables**:
- PoC実装（最小限のプロジェクトCRUD）
- パフォーマンスベンチマーク結果
- 技術的制約・推奨事項のドキュメント

### 0.2 開発環境のセットアップ

**Goal**: 開発者が即座に開発を開始できる環境を構築

**Tasks**:
1. Docker Compose による開発環境構築
   - PostgreSQL コンテナ
   - Backend (Go) コンテナ (ホットリロード設定)
   - Frontend (Next.js) コンテナ
2. Makefile によるコマンド統一化
3. 環境変数管理（.env.example）
4. README.md にセットアップ手順を記載

**Deliverables**:
- `docker-compose.yml`
- `Makefile`
- セットアップドキュメント

### 0.3 CI/CD パイプライン基礎設計

**Goal**: 自動テスト・ビルドの基盤を構築

**Tasks**:
1. GitHub Actions ワークフロー設計
   - Backend: lint, test, build
   - Frontend: lint, type-check, test, build
2. Pre-commit hooks 設定（linter実行）
3. ブランチ戦略の確認（feature/*, develop, main）

**Deliverables**:
- `.github/workflows/` ファイル群
- ブランチ戦略ドキュメント

## Phase 1: Architecture & Contracts

### 1.1 データモデル設計

**Goal**: ER図とデータベーススキーマを確定

**Tasks**:
1. エンティティ定義（User, Project, Task, Member, TimeEntry, ProjectMember）
2. リレーションシップ設計
3. インデックス戦略（検索・ソートの最適化）
4. マイグレーションファイル作成

**Deliverables**:
- ER図（Mermaid形式）
- マイグレーションファイル（`.sql`）
- `data-model.md`

### 1.2 API契約設計

**Goal**: フロントエンドとバックエンド間のAPI仕様を確定

**Tasks**:
1. OpenAPI 3.0 スキーマ作成
   - 認証エンドポイント (`POST /api/v1/auth/login`, `/register`)
   - プロジェクトCRUD (`/api/v1/projects`)
   - タスクCRUD (`/api/v1/projects/{id}/tasks`)
   - メンバー管理 (`/api/v1/members`)
   - 収支データ (`/api/v1/projects/{id}/budget`)
   - 統計・集計 (`/api/v1/projects/{id}/analytics`)
2. エラーレスポンス標準化
3. ページネーション設計（cursor vs offset）
4. フィルタ・ソートパラメータ設計

**Deliverables**:
- `openapi.yaml`
- TypeScript型定義（自動生成 or 手動作成）
- API設計ドキュメント

### 1.3 フロントエンド設計

**Goal**: 画面構成とコンポーネント設計を確定

**Tasks**:
1. 画面フロー図作成（ユーザージャーニーマップ）
2. ワイヤーフレーム作成（主要画面）
   - ダッシュボード
   - プロジェクト一覧
   - プロジェクト詳細
   - タスク管理
   - 収支管理
   - グラフ表示
3. コンポーネント階層設計
4. ルーティング設計（Next.js App Router）
5. 状態管理戦略（サーバー状態 vs クライアント状態）

**Deliverables**:
- ワイヤーフレーム（Figma or 手書き）
- コンポーネント構成図
- `frontend-design.md`

### 1.4 認証・認可設計

**Goal**: セキュアな認証フローを設計

**Tasks**:
1. JWT認証フロー設計（Access Token + Refresh Token）
2. パスワードハッシュ化戦略（bcrypt）
3. 権限モデル設計（RBAC: Admin, Manager, Member）
4. セッション管理（Token有効期限、Refresh戦略）

**Deliverables**:
- 認証フロー図
- 認証・認可仕様書
- `auth-design.md`

## Phase 2: MVP Implementation (P1 Features)

### 2.1 Backend: 基盤実装

**Goal**: バックエンドの基本構造とDB接続を実装

**Tasks**:
1. プロジェクト構造セットアップ
2. Echo サーバー初期化（ルーター、ミドルウェア設定）
3. GORM初期化とDB接続
4. マイグレーション実行システム
5. ログ設定（構造化ログ）
6. エラーハンドリング統一
7. CORS設定

**Acceptance Criteria**:
- サーバーが起動し、ヘルスチェックエンドポイント (`/health`) が正常応答
- DBマイグレーションが正常実行される

### 2.2 Backend: 認証機能

**Goal**: ユーザー登録・ログイン機能を実装

**Tasks**:
1. Userモデル実装
2. ユーザー登録エンドポイント (`POST /api/v1/auth/register`)
3. ログインエンドポイント (`POST /api/v1/auth/login`)
4. JWT生成・検証ミドルウェア
5. パスワードハッシュ化（bcrypt）
6. バリデーション実装

**Acceptance Criteria**:
- ユーザー登録が成功し、DBにユーザーが保存される
- ログインが成功し、JWTトークンが返却される
- 不正なトークンでのアクセスが拒否される

### 2.3 Backend: プロジェクト管理API

**Goal**: プロジェクトのCRUD操作を実装 (User Story 1)

**Tasks**:
1. Projectモデル実装
2. プロジェクト作成 (`POST /api/v1/projects`)
3. プロジェクト一覧取得 (`GET /api/v1/projects`)
   - ページネーション
   - ステータスフィルタ
   - 検索機能
4. プロジェクト詳細取得 (`GET /api/v1/projects/{id}`)
5. プロジェクト更新 (`PUT /api/v1/projects/{id}`)
6. プロジェクト削除（論理削除） (`DELETE /api/v1/projects/{id}`)

**Acceptance Criteria**:
- 全てのエンドポイントが正常動作する
- バリデーションが正しく機能する
- 論理削除されたプロジェクトは一覧に表示されない

### 2.4 Backend: タスク管理API

**Goal**: タスクのCRUD操作と工数管理を実装 (User Story 2)

**Tasks**:
1. Taskモデル実装
2. タスク作成 (`POST /api/v1/projects/{projectId}/tasks`)
3. タスク一覧取得 (`GET /api/v1/projects/{projectId}/tasks`)
4. タスク詳細取得 (`GET /api/v1/tasks/{id}`)
5. タスク更新 (`PUT /api/v1/tasks/{id}`)
6. タスク削除 (`DELETE /api/v1/tasks/{id}`)
7. 予実差異計算ロジック（Service層）

**Acceptance Criteria**:
- タスクがプロジェクトに紐づいて管理される
- 予定工数と実績工数が正しく記録される
- 予実差異が自動計算される

### 2.5 Backend: 収支管理API

**Goal**: 売上・コスト・利益の計算を実装 (User Story 3)

**Tasks**:
1. Budgetモデル、Memberモデル、TimeEntryモデル実装
2. 売上登録 (`PUT /api/v1/projects/{id}/budget/revenue`)
3. メンバー単価設定 (`POST /api/v1/members`, `PUT /api/v1/members/{id}`)
4. 工数記録 (`POST /api/v1/time-entries`)
5. コスト自動計算ロジック（工数 × 単価）
6. 利益計算ロジック（売上 - コスト）
7. 収支サマリー取得 (`GET /api/v1/projects/{id}/budget`)

**Acceptance Criteria**:
- 売上が登録できる
- メンバー単価が設定できる
- 工数記録からコストが自動計算される
- 利益・利益率が正しく算出される

### 2.6 Frontend: 基盤実装

**Goal**: Next.jsプロジェクトのセットアップと共通機能実装

**Tasks**:
1. Next.js + TypeScript プロジェクト初期化
2. shadcn/ui セットアップ（必要なコンポーネントインストール）
3. Tailwind CSS 設定
4. レイアウトコンポーネント（Header, Sidebar, Footer）
5. API クライアント実装（axios or fetch wrapper）
6. TanStack Query セットアップ
7. Zustand store セットアップ
8. 共通コンポーネント（Loading, ErrorBoundary, Pagination）

**Acceptance Criteria**:
- Next.jsアプリが起動する
- 基本レイアウトが表示される
- API通信が正常に行える

### 2.7 Frontend: 認証画面

**Goal**: ログイン・登録画面を実装

**Tasks**:
1. ログインページ (`/login`)
2. 登録ページ (`/register`)
3. 認証フォーム実装（Zod バリデーション）
4. JWT保存（Cookie or LocalStorage）
5. 認証状態管理（Zustand）
6. 保護されたルートの実装（middleware）

**Acceptance Criteria**:
- ユーザー登録が正常に完了する
- ログインが成功し、トークンが保存される
- 未認証時は保護されたページにアクセスできない

### 2.8 Frontend: プロジェクト管理画面

**Goal**: プロジェクトの一覧・作成・編集画面を実装 (User Story 1)

**Tasks**:
1. プロジェクト一覧ページ (`/projects`)
   - プロジェクトカード表示
   - 検索・フィルタ機能
   - ページネーション
2. プロジェクト詳細ページ (`/projects/[id]`)
   - 基本情報表示
   - タブナビゲーション（概要、タスク、収支、グラフ）
3. プロジェクト作成ページ (`/projects/new`)
   - フォーム実装（名前、予算、期間、説明）
   - バリデーション
4. プロジェクト編集機能
   - インラインまたはモーダル編集

**Acceptance Criteria**:
- プロジェクト一覧が表示される
- 新規プロジェクトが作成できる
- プロジェクト詳細が表示される
- プロジェクト情報が編集できる

### 2.9 Frontend: タスク管理画面

**Goal**: タスクの一覧・作成・編集と予実表示を実装 (User Story 2)

**Tasks**:
1. タスク一覧表示（プロジェクト詳細内）
   - Table コンポーネント使用
   - 予定工数・実績工数・差異の表示
   - 超過時の警告色表示
2. タスク作成フォーム（Dialog or 専用ページ）
3. タスク編集機能（インライン or Dialog）
4. 予実サマリー表示
   - 合計予定工数
   - 合計実績工数
   - 合計差異
   - 達成率

**Acceptance Criteria**:
- タスク一覧が表形式で表示される
- タスクが作成できる
- 予定工数と実績工数が入力できる
- 予実差異が自動表示される
- 超過時に警告色で表示される

### 2.10 Frontend: 収支管理画面

**Goal**: 売上・コスト・利益の表示を実装 (User Story 3)

**Tasks**:
1. 収支サマリーコンポーネント
   - 売上表示
   - 総コスト表示
   - 利益表示
   - 利益率表示
2. 売上入力フォーム
3. コスト内訳表示
   - メンバー別コスト
   - タスク別コスト
4. 赤字時の警告表示

**Acceptance Criteria**:
- 売上が入力できる
- コストが自動計算され表示される
- 利益・利益率が表示される
- 赤字時に警告表示される

### 2.11 Testing & Quality Assurance

**Goal**: MVP機能の品質を保証

**Tasks**:
1. Backend単体テスト作成
   - Repository層テスト
   - Service層テスト（ビジネスロジック）
2. Backend統合テスト作成
   - API エンドポイントテスト
3. Frontend単体テスト作成
   - コンポーネントテスト（主要なもの）
4. Frontend E2Eテスト作成（Playwright）
   - プロジェクト作成フロー
   - タスク管理フロー
   - 収支確認フロー
5. 手動テスト実施（User Story 1-3のAcceptance Scenarios）

**Acceptance Criteria**:
- テストカバレッジ: Backend 70%+, Frontend 60%+
- 全てのE2Eテストがパス
- 全てのAcceptance Scenariosが満たされる

## Phase 3: P2 Features Implementation

### 3.1 Backend: ダッシュボードAPI

**Goal**: 複数プロジェクトの集計データAPIを実装 (User Story 4)

**Tasks**:
1. ダッシュボードサマリーエンドポイント (`GET /api/v1/dashboard`)
   - 全プロジェクト数
   - 進行中プロジェクト数
   - 総売上・総利益
   - KPI集計（工数達成率、利益率平均など）
2. プロジェクト一覧の拡張フィルタ・ソート
3. 集計クエリのパフォーマンス最適化

**Acceptance Criteria**:
- ダッシュボードデータが2秒以内に返却される（50プロジェクト）
- フィルタ・ソートが正常動作する

### 3.2 Backend: 分析・統計API

**Goal**: グラフ表示用のデータAPIを実装 (User Story 5)

**Tasks**:
1. 予実比較データ (`GET /api/v1/projects/{id}/analytics/plan-actual`)
2. 収支データ (`GET /api/v1/projects/{id}/analytics/budget`)
3. 月次推移データ (`GET /api/v1/projects/{id}/analytics/trends?period=monthly`)
4. タスク別工数割合 (`GET /api/v1/projects/{id}/analytics/task-distribution`)
5. 複数プロジェクト比較データ (`GET /api/v1/analytics/projects-comparison`)

**Acceptance Criteria**:
- 各エンドポイントが適切な形式でデータを返す
- グラフ描画に必要な全データが含まれる

### 3.3 Frontend: ダッシュボード

**Goal**: 複数プロジェクトの横断ビューを実装 (User Story 4)

**Tasks**:
1. ダッシュボードページ (`/dashboard`)
2. KPIカード表示（総プロジェクト数、進行中、総利益など）
3. プロジェクト一覧テーブル（拡張版）
   - ステータス別フィルタ
   - KPIソート
   - クイックアクション（詳細表示、編集）

**Acceptance Criteria**:
- ダッシュボードが2秒以内に表示される
- フィルタ・ソートが直感的に操作できる
- 各プロジェクトの状態が一目で分かる

### 3.4 Frontend: グラフ可視化

**Goal**: 各種グラフコンポーネントを実装 (User Story 5)

**Tasks**:
1. 予実比較グラフ（棒グラフ or 折れ線グラフ）
2. 収支グラフ（棒グラフ）
3. 月次推移グラフ（折れ線グラフ）
4. タスク別工数割合（円グラフ）
5. グラフ共通機能
   - ツールチップ
   - レスポンシブ対応
   - カラーテーマ統一
6. アナリティクスページ (`/projects/[id]/analytics`)

**Acceptance Criteria**:
- 全てのグラフが正しくデータを反映する
- グラフが3秒以内に描画される（100データポイント）
- レスポンシブでモバイルでも見やすい

## Phase 4: P3 Features & Enhancements

### 4.1 Backend: メンバー管理機能拡張

**Goal**: メンバー管理とリソース稼働機能を実装 (User Story 6)

**Tasks**:
1. メンバー詳細API
2. プロジェクトメンバー割り当てAPI
3. メンバー別稼働レポートAPI
4. リソース稼働率計算API

**Acceptance Criteria**:
- メンバー管理が完全に機能する
- 稼働レポートが正確に生成される

### 4.2 Backend: データエクスポートAPI

**Goal**: CSV/PDFエクスポート機能を実装 (User Story 7)

**Tasks**:
1. CSVエクスポートエンドポイント
   - プロジェクトデータ
   - タスクデータ
   - 収支データ
2. PDF生成ライブラリ統合（必要に応じて）

**Acceptance Criteria**:
- CSVが正しいフォーマットで出力される
- 大量データのエクスポートがタイムアウトしない

### 4.3 Frontend: メンバー管理画面

**Goal**: メンバー管理UIを実装 (User Story 6)

**Tasks**:
1. メンバー一覧ページ
2. メンバー登録・編集フォーム
3. プロジェクトメンバー割り当て画面
4. メンバー別稼働レポート表示

**Acceptance Criteria**:
- メンバー管理が直感的に操作できる
- 稼働レポートが見やすく表示される

### 4.4 Frontend: エクスポート機能

**Goal**: データエクスポートUIを実装 (User Story 7)

**Tasks**:
1. エクスポートボタンの実装（各画面）
2. CSV/PDFダウンロード機能
3. エクスポート進行状況の表示

**Acceptance Criteria**:
- ワンクリックでデータがエクスポートされる
- ダウンロードが正常に完了する

## Phase 5: Polish & Production Readiness

### 5.1 パフォーマンス最適化

**Tasks**:
1. フロントエンド
   - コード分割（dynamic import）
   - 画像最適化（Next.js Image）
   - バンドルサイズ削減
2. バックエンド
   - N+1クエリ問題の解消
   - インデックス最適化
   - クエリキャッシング
3. パフォーマンステスト実施

**Acceptance Criteria**:
- Success Criteriaの全パフォーマンス目標を達成

### 5.2 セキュリティ強化

**Tasks**:
1. セキュリティ監査実施
2. 脆弱性スキャン（npm audit, go mod verify）
3. HTTPS設定確認
4. CORS設定最終確認
5. レート制限実装（必要に応じて）
6. CSP設定

**Acceptance Criteria**:
- 既知の脆弱性がゼロ
- セキュリティベストプラクティスに準拠

### 5.3 ドキュメント整備

**Tasks**:
1. API仕様書の最終更新（OpenAPI）
2. ユーザーマニュアル作成
3. 運用マニュアル作成
4. デプロイ手順書
5. トラブルシューティングガイド

**Acceptance Criteria**:
- 全ドキュメントが最新状態
- 新規開発者がドキュメントのみで環境構築可能

### 5.4 本番環境デプロイ準備

**Tasks**:
1. 本番用Dockerイメージ最適化
2. 環境変数管理（Secrets）
3. データベースバックアップ設定
4. ログ集約設定（CloudWatch, Datadog等）
5. モニタリング設定（APM, メトリクス）
6. アラート設定

**Acceptance Criteria**:
- 本番デプロイが自動化されている
- 監視・アラートが機能している

## Milestones & Timeline

### Milestone 1: 環境構築完了 (Week 1)
- Phase 0完了
- 開発環境が全開発者で動作

### Milestone 2: MVP完成 (Week 4-6)
- Phase 2完了
- P1機能（User Story 1-3）が動作
- 基本的なE2Eテストがパス

### Milestone 3: P2機能完成 (Week 8-10)
- Phase 3完了
- ダッシュボードとグラフ機能が動作

### Milestone 4: 全機能完成 (Week 12-14)
- Phase 4完了
- P3機能が動作

### Milestone 5: 本番リリース (Week 16)
- Phase 5完了
- 本番環境にデプロイ
- ユーザー受け入れテスト完了

## Dependencies & Risks

### Critical Dependencies
1. PostgreSQL データベース環境
2. 開発メンバーのGo/TypeScript習熟度
3. デザインリソース（shadcn/uiで軽減）

### Technical Risks
| Risk | Impact | Mitigation |
|------|--------|-----------|
| Echo/GORMの学習曲線 | 中 | Phase 0でPoC作成、ドキュメント整備 |
| Next.js App Routerの複雑性 | 中 | 公式ドキュメント参照、シンプルなルーティング設計 |
| グラフパフォーマンス | 高 | 早期にパフォーマンステスト、必要に応じてサーバーサイド集計 |
| データ量増加時のスケーラビリティ | 中 | 初期からインデックス設計、後続でキャッシング追加 |

### Project Risks
| Risk | Impact | Mitigation |
|------|--------|-----------|
| 要件の曖昧さ | 高 | Questions for Clarificationへの早期回答取得 |
| スコープクリープ | 中 | P1/P2/P3の優先度厳守、機能追加は次イテレーション |
| リソース不足 | 高 | 段階的リリース、MVPを最優先 |

## Testing Strategy

### Unit Testing
- **Backend**: 全Repository、Service層に対してテスト作成（カバレッジ70%+）
- **Frontend**: 重要なコンポーネント、hooks、utilsにテスト作成（カバレッジ60%+）

### Integration Testing
- **Backend**: 全APIエンドポイントに対して統合テスト
- **Frontend**: API通信を含む主要フロー

### E2E Testing
- Playwright使用
- クリティカルユーザーフロー（プロジェクト作成→タスク追加→収支確認）

### Performance Testing
- Apache Bench or k6 によるロードテスト
- 目標: 100同時接続、応答時間p95 < 500ms

### Security Testing
- OWASP ZAP によるセキュリティスキャン
- 依存関係の脆弱性チェック（npm audit, go mod verify）

## Success Metrics

実装の成功は以下の指標で評価：

### Functional Metrics
- ✅ 全User Stories (P1-P3) のAcceptance Criteriaを満たす
- ✅ 24の機能要件（FR-001〜FR-024）を実装

### Performance Metrics
- ✅ API応答時間: 平均 < 200ms, p95 < 500ms
- ✅ ダッシュボード表示: < 2秒 (50プロジェクト)
- ✅ グラフ描画: < 3秒 (100データポイント)

### Quality Metrics
- ✅ テストカバレッジ: Backend 70%+, Frontend 60%+
- ✅ E2Eテスト成功率: 100%
- ✅ 既知の脆弱性: 0件
- ✅ コードレビュー: 全PRで実施

### User Experience Metrics
- ✅ プロジェクト作成時間: < 5分
- ✅ 予実確認: 3クリック以内
- ✅ 収支確認: 2クリック以内
- ✅ 初回利用時の工数入力成功率: 90%+（ユーザーテスト）

## Next Actions

1. **仕様明確化**: Questions for Clarificationへの回答取得
2. **Phase 0開始**: 技術スタック検証とPoC作成
3. **チーム編成**: フロントエンド担当、バックエンド担当の役割分担
4. **キックオフミーティング**: この計画の共有と合意形成
5. **開発環境構築**: 全開発者が環境を立ち上げ
6. **スプリント計画**: 2週間スプリントでのタスク分解

---

**Plan Status**: Draft  
**Last Updated**: 2025-10-11  
**Next Review**: Phase 0完了時
