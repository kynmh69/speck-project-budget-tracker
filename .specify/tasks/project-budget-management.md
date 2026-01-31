# Tasks: プロジェクト予算管理システム

**Input**: Design documents from `.specify/specs/` and `.specify/plans/`  
**Prerequisites**: plan.md, spec.md, data-model.md (Phase 1で作成), contracts/ (Phase 1で作成)

**Organization**: タスクはユーザーストーリーごとにグループ化し、各ストーリーが独立して実装・テスト可能な構成とする。

## Format: `[ID] [P?] [Story] Description`
- **[P]**: 並列実行可能（異なるファイル、依存関係なし）
- **[Story]**: タスクが属するユーザーストーリー（例: US1, US2, US3）
- ファイルパスを具体的に記載

## Path Conventions
- **Web app**: `backend/` (Golang), `frontend/` (Next.js)
- Backend: `backend/internal/`, `backend/cmd/`, `backend/migrations/`
- Frontend: `frontend/src/app/`, `frontend/src/components/`, `frontend/src/lib/`

---

## Phase 0: 環境構築とリサーチ (Week 1)

**Purpose**: 開発環境の構築と技術スタックの検証

### Setup

- [ ] T001 リポジトリのルートディレクトリ構造作成 (`backend/`, `frontend/`, `.github/`, `docs/`)
- [ ] T002 [P] `.gitignore` の作成（Go、Node.js、環境変数、IDE設定）
- [ ] T003 [P] `README.md` の初期版作成（プロジェクト概要、セットアップ手順）

### Backend Setup

- [ ] T004 `backend/go.mod` 初期化（Go 1.21+, Echo, GORM, JWT-go, validator等）
- [ ] T005 [P] `backend/cmd/server/main.go` エントリーポイント作成
- [ ] T006 [P] `backend/internal/config/config.go` 環境変数管理実装
- [ ] T007 [P] `backend/.env.example` 作成（DB接続、JWT秘密鍵、ポート等）
- [ ] T008 Echo サーバー基本設定（`backend/internal/router/router.go`）
- [ ] T009 [P] ミドルウェア設定（`backend/internal/middleware/logger.go`, `cors.go`）
- [ ] T010 [P] ヘルスチェックエンドポイント（`GET /health`）実装

### Frontend Setup

- [ ] T011 Next.js 14 プロジェクト初期化（`frontend/`, TypeScript, App Router）
- [ ] T012 [P] shadcn/ui セットアップ（`components.json`, Tailwind設定）
- [ ] T013 [P] `frontend/.env.local.example` 作成（API URL等）
- [ ] T014 [P] 基本コンポーネントインストール（Button, Card, Input, Form, Dialog, Table）
- [ ] T015 [P] グローバルレイアウト作成（`frontend/src/app/layout.tsx`）
- [ ] T016 [P] Tailwind CSS カスタムテーマ設定（色、フォント）

### Docker & Development Environment

- [ ] T017 PostgreSQL用 `docker-compose.yml` 作成
- [ ] T018 [P] Backend用 `Dockerfile` 作成（マルチステージビルド）
- [ ] T019 [P] Frontend用 `Dockerfile` 作成
- [ ] T020 `Makefile` 作成（起動、ビルド、テスト、マイグレーション等のコマンド）
- [ ] T021 開発環境動作確認（`docker-compose up`でDB、Backend、Frontendが起動）

### CI/CD Setup

- [ ] T022 [P] GitHub Actions: Backend CI（`.github/workflows/backend-ci.yml`）
- [ ] T023 [P] GitHub Actions: Frontend CI（`.github/workflows/frontend-ci.yml`）
- [ ] T024 [P] Pre-commit hooks設定（Husky or lefthook）

### Technical Validation

- [ ] T025 Echo + GORM でシンプルなCRUD PoC作成（テスト用エンティティ）
- [ ] T026 [P] PostgreSQL接続確認とマイグレーション実行テスト
- [ ] T027 [P] JWT認証フローのPoC実装
- [ ] T028 [P] Next.js + TanStack Query でAPI呼び出しPoC
- [ ] T029 Recharts でグラフ描画パフォーマンステスト（100データポイント）

**Checkpoint**: 開発環境が完全に動作し、全開発者が `docker-compose up` でローカル開発可能

---

## Phase 1: アーキテクチャ設計とAPI契約 (Week 2-3)

**Purpose**: データモデル、API仕様、画面設計の確定

### Data Model Design

- [ ] T030 ER図作成（`.specify/data-models/project-budget-management.md`）
- [ ] T031 DBマイグレーション: Users テーブル（`backend/migrations/000001_create_users.up.sql`）
- [ ] T032 [P] DBマイグレーション: Projects テーブル（`backend/migrations/000002_create_projects.up.sql`）
- [ ] T033 [P] DBマイグレーション: Tasks テーブル（`backend/migrations/000003_create_tasks.up.sql`）
- [ ] T034 [P] DBマイグレーション: Members テーブル（`backend/migrations/000004_create_members.up.sql`）
- [ ] T035 [P] DBマイグレーション: TimeEntries テーブル（`backend/migrations/000005_create_time_entries.up.sql`）
- [ ] T036 DBマイグレーション: ProjectMembers テーブル（`backend/migrations/000006_create_project_members.up.sql`）
- [ ] T037 [P] インデックス追加マイグレーション（検索・ソート最適化）

### API Contract Design

- [ ] T038 OpenAPI 3.0 スキーマ作成（`.specify/contracts/project-budget-management/openapi.yaml`）
- [ ] T039 [P] 認証エンドポイント定義（`/api/v1/auth/login`, `/register`）
- [ ] T040 [P] プロジェクトCRUDエンドポイント定義（`/api/v1/projects`）
- [ ] T041 [P] タスクCRUDエンドポイント定義（`/api/v1/projects/{id}/tasks`）
- [ ] T042 [P] メンバー管理エンドポイント定義（`/api/v1/members`）
- [ ] T043 [P] 収支データエンドポイント定義（`/api/v1/projects/{id}/budget`）
- [ ] T044 [P] 分析データエンドポイント定義（`/api/v1/projects/{id}/analytics`）
- [ ] T045 エラーレスポンス標準化とページネーション設計

### Frontend Design

- [ ] T046 画面フロー図作成（ユーザージャーニーマップ）
- [ ] T047 ワイヤーフレーム作成（ダッシュボード、プロジェクト一覧、詳細、タスク、収支、グラフ）
- [ ] T048 TypeScript型定義作成（`frontend/src/types/`）- OpenAPIから生成 or 手動
- [ ] T049 ルーティング設計ドキュメント（App Router構造）

### Authentication Design

- [ ] T050 JWT認証フロー設計ドキュメント（Access/Refresh Token）
- [ ] T051 権限モデル設計（RBAC: Admin, Manager, Member）

**Checkpoint**: 全てのAPI仕様とデータモデルが確定し、実装準備完了

---

## Phase 2: Foundational（基盤実装） (Week 3-4)

**Purpose**: 全ユーザーストーリーの実装前に必要な基盤機能

**⚠️ CRITICAL**: このフェーズが完了するまで、ユーザーストーリーの実装は開始不可

### Backend Foundation

- [ ] T052 GORMモデル: User（`backend/internal/models/user.go`）
- [ ] T053 [P] GORMモデル: Project（`backend/internal/models/project.go`）
- [ ] T054 [P] GORMモデル: Task（`backend/internal/models/task.go`）
- [ ] T055 [P] GORMモデル: Member（`backend/internal/models/member.go`）
- [ ] T056 [P] GORMモデル: TimeEntry（`backend/internal/models/time_entry.go`）
- [ ] T057 [P] GORMモデル: ProjectMember（`backend/internal/models/project_member.go`）
- [ ] T058 DB接続とマイグレーション実行（`backend/internal/database/database.go`）
- [ ] T059 エラーハンドリング統一（`backend/internal/errors/errors.go`）
- [ ] T060 [P] レスポンスDTO基本構造（`backend/internal/dto/response.go`）
- [ ] T061 [P] カスタムバリデーター（`backend/internal/validator/custom_validators.go`）

### Backend Authentication

- [ ] T062 AuthService実装（`backend/internal/service/auth_service.go`）- 登録、ログイン、JWT生成
- [ ] T063 パスワードハッシュ化ユーティリティ（bcrypt）
- [ ] T064 JWT生成・検証ミドルウェア（`backend/internal/middleware/auth.go`）
- [ ] T065 AuthHandler実装（`backend/internal/handler/auth_handler.go`）
- [ ] T066 認証エンドポイント登録（`POST /api/v1/auth/register`, `/login`）
- [ ] T067 [P] 認証統合テスト（`backend/tests/integration/auth_test.go`）

### Frontend Foundation

- [ ] T068 APIクライアント実装（`frontend/src/lib/api-client.ts`）- axios/fetch wrapper
- [ ] T069 TanStack Query セットアップ（`frontend/src/lib/query-client.ts`）
- [ ] T070 Zustand認証ストア（`frontend/src/store/auth-store.ts`）
- [ ] T071 [P] 認証ヘルパー関数（`frontend/src/lib/auth.ts`）- Token管理
- [ ] T072 [P] 共通UIコンポーネント: Loading（`frontend/src/components/common/loading.tsx`）
- [ ] T073 [P] 共通UIコンポーネント: ErrorBoundary（`frontend/src/components/common/error-boundary.tsx`）
- [ ] T074 [P] 共通UIコンポーネント: Pagination（`frontend/src/components/common/pagination.tsx`）
- [ ] T075 レイアウトコンポーネント: Header（`frontend/src/components/layout/header.tsx`）
- [ ] T076 [P] レイアウトコンポーネント: Sidebar（`frontend/src/components/layout/sidebar.tsx`）
- [ ] T077 ダッシュボードレイアウト（`frontend/src/app/(dashboard)/layout.tsx`）
- [ ] T078 認証ガードミドルウェア（Next.js middleware.ts）

### Frontend Authentication

- [ ] T079 Zodスキーマ: 認証（`frontend/src/schemas/auth-schema.ts`）
- [ ] T080 ログインページ（`frontend/src/app/(auth)/login/page.tsx`）
- [ ] T081 [P] 登録ページ（`frontend/src/app/(auth)/register/page.tsx`）
- [ ] T082 認証フック（`frontend/src/hooks/use-auth.ts`）- TanStack Query
- [ ] T083 ログイン・登録フォームコンポーネント実装

**Checkpoint**: 認証が完全に動作し、保護されたルートが機能。ユーザーストーリー実装開始可能。

---

## Phase 3: User Story 1 - プロジェクト管理 (Priority: P1) 🎯 MVP (Week 4-5)

**Goal**: プロジェクトの作成・編集・一覧表示機能を提供

**Independent Test**: プロジェクト一覧画面で新規プロジェクトを作成し、詳細画面で編集できることを確認

### Backend: Project CRUD

- [ ] T084 [US1] ProjectRepository実装（`backend/internal/repository/project_repository.go`）
- [ ] T085 [US1] ProjectService実装（`backend/internal/service/project_service.go`）
- [ ] T086 [US1] ProjectDTO実装（`backend/internal/dto/project_dto.go`）
- [ ] T087 [US1] ProjectHandler実装（`backend/internal/handler/project_handler.go`）
- [ ] T088 [US1] プロジェクト作成エンドポイント（`POST /api/v1/projects`）
- [ ] T089 [US1] プロジェクト一覧取得（`GET /api/v1/projects`）- ページネーション、フィルタ、検索
- [ ] T090 [US1] プロジェクト詳細取得（`GET /api/v1/projects/{id}`）
- [ ] T091 [US1] プロジェクト更新（`PUT /api/v1/projects/{id}`）
- [ ] T092 [US1] プロジェクト削除/論理削除（`DELETE /api/v1/projects/{id}`）
- [ ] T093 [P] [US1] ProjectService単体テスト（`backend/tests/unit/service/project_service_test.go`）
- [ ] T094 [P] [US1] Project API統合テスト（`backend/tests/integration/project_test.go`）

### Frontend: Project Management UI

- [ ] T095 [US1] TypeScript型定義: Project（`frontend/src/types/project.ts`）
- [ ] T096 [US1] Zodスキーマ: Project（`frontend/src/schemas/project-schema.ts`）
- [ ] T097 [US1] プロジェクトフック（`frontend/src/hooks/use-projects.ts`）- TanStack Query
- [ ] T098 [US1] プロジェクトカードコンポーネント（`frontend/src/components/projects/project-card.tsx`）
- [ ] T099 [US1] プロジェクトフォームコンポーネント（`frontend/src/components/projects/project-form.tsx`）
- [ ] T100 [US1] プロジェクトフィルターコンポーネント（`frontend/src/components/projects/project-filters.tsx`）
- [ ] T101 [US1] プロジェクト一覧ページ（`frontend/src/app/(dashboard)/projects/page.tsx`）
- [ ] T102 [US1] プロジェクト作成ページ（`frontend/src/app/(dashboard)/projects/new/page.tsx`）
- [ ] T103 [US1] プロジェクト詳細ページ（`frontend/src/app/(dashboard)/projects/[id]/page.tsx`）
- [ ] T104 [US1] プロジェクト編集機能（Dialog or inline編集）
- [ ] T105 [P] [US1] プロジェクト管理E2Eテスト（`frontend/tests/e2e/projects.spec.ts`）

**Checkpoint**: プロジェクトのCRUD操作が完全に機能し、独立してテスト可能

---

## Phase 4: User Story 2 - 工数予実管理 (Priority: P1) 🎯 MVP (Week 5-6)

**Goal**: タスクの予定工数・実績工数の入力と予実差異の表示

**Independent Test**: プロジェクト内でタスクを作成し、予定・実績工数を入力して差異を確認

### Backend: Task Management

- [X] T106 [US2] TaskRepository実装（`backend/internal/repository/task_repository.go`）
- [X] T107 [US2] TaskService実装（`backend/internal/service/task_service.go`）- 予実差異計算含む
- [X] T108 [US2] TaskDTO実装（`backend/internal/dto/task_dto.go`）
- [X] T109 [US2] TaskHandler実装（`backend/internal/handler/task_handler.go`）
- [X] T110 [US2] タスク作成エンドポイント（`POST /api/v1/projects/{projectId}/tasks`）
- [X] T111 [US2] タスク一覧取得（`GET /api/v1/projects/{projectId}/tasks`）
- [X] T112 [US2] タスク詳細取得（`GET /api/v1/tasks/{id}`）
- [X] T113 [US2] タスク更新（`PUT /api/v1/tasks/{id}`）- 予定・実績工数更新
- [X] T114 [US2] タスク削除（`DELETE /api/v1/tasks/{id}`）
- [X] T115 [US2] プロジェクト予実サマリー取得（`GET /api/v1/projects/{id}/summary`）
- [X] T116 [P] [US2] TaskService単体テスト（`backend/tests/unit/service/task_service_test.go`）
- [X] T117 [P] [US2] Task API統合テスト（`backend/tests/integration/task_test.go`）

### Frontend: Task Management UI

- [X] T118 [US2] TypeScript型定義: Task（`frontend/src/types/task.ts`）
- [X] T119 [US2] Zodスキーマ: Task（`frontend/src/schemas/task-schema.ts`）
- [X] T120 [US2] タスクフック（`frontend/src/hooks/use-tasks.ts`）
- [X] T121 [US2] タスクリストコンポーネント（`frontend/src/components/tasks/task-list.tsx`）- Table使用
- [X] T122 [US2] タスクアイテムコンポーネント（`frontend/src/components/tasks/task-item.tsx`）- 予実差異表示
- [X] T123 [US2] タスクフォームコンポーネント（`frontend/src/components/tasks/task-form.tsx`）
- [X] T124 [US2] 予実サマリーコンポーネント（`frontend/src/components/tasks/plan-actual-summary.tsx`）
- [X] T125 [US2] タスク管理ページ（`frontend/src/app/(dashboard)/projects/[id]/tasks/page.tsx`）
- [X] T126 [US2] 工数超過時の警告色表示ロジック実装
- [ ] T127 [P] [US2] タスク管理E2Eテスト（`frontend/tests/e2e/tasks.spec.ts`）

**Checkpoint**: タスクの予実管理が完全に機能し、予実差異が正確に表示される

---

## Phase 5: User Story 3 - 収支管理 (Priority: P1) 🎯 MVP (Week 6-7)

**Goal**: 売上・コスト・利益の管理と自動計算

**Independent Test**: プロジェクトに売上を入力し、メンバー単価と工数からコストが自動計算され、利益が表示される

### Backend: Budget Management

- [X] T128 [US3] MemberRepository実装（`backend/internal/repository/member_repository.go`）
- [X] T129 [US3] TimeEntryRepository実装（`backend/internal/repository/time_entry_repository.go`）
- [X] T130 [US3] BudgetService実装（`backend/internal/service/budget_service.go`）- コスト・利益計算
- [X] T131 [US3] MemberService実装（`backend/internal/service/member_service.go`）
- [X] T132 [US3] BudgetDTO実装（`backend/internal/dto/budget_dto.go`）
- [X] T133 [US3] MemberDTO実装（`backend/internal/dto/member_dto.go`）
- [X] T134 [US3] BudgetHandler実装（`backend/internal/handler/budget_handler.go`）
- [X] T135 [US3] MemberHandler実装（`backend/internal/handler/member_handler.go`）
- [X] T136 [US3] 売上登録エンドポイント（`PUT /api/v1/projects/{id}/budget/revenue`）
- [X] T137 [US3] メンバー作成エンドポイント（`POST /api/v1/members`）
- [X] T138 [US3] メンバー単価設定（`PUT /api/v1/members/{id}`）
- [X] T139 [US3] 工数記録エンドポイント（`POST /api/v1/time-entries`）
- [X] T140 [US3] 収支サマリー取得（`GET /api/v1/projects/{id}/budget`）- コスト・利益計算
- [ ] T141 [P] [US3] BudgetService単体テスト（`backend/tests/unit/service/budget_service_test.go`）
- [ ] T142 [P] [US3] Budget API統合テスト（`backend/tests/integration/budget_test.go`）

### Frontend: Budget Management UI

- [X] T143 [US3] TypeScript型定義: Budget, Member（`frontend/src/types/budget.ts`, `member.ts`）
- [X] T144 [US3] Zodスキーマ: Member（`frontend/src/schemas/member-schema.ts`）
- [X] T145 [US3] 収支フック（`frontend/src/hooks/use-budget.ts`）
- [X] T146 [US3] メンバーフック（`frontend/src/hooks/use-members.ts`）
- [X] T147 [US3] 収支サマリーコンポーネント（`frontend/src/components/budget/budget-summary.tsx`）
- [X] T148 [US3] コスト内訳コンポーネント（`frontend/src/components/budget/cost-breakdown.tsx`）
- [X] T149 [US3] 売上入力フォーム実装
- [X] T150 [US3] 収支管理ページ（`frontend/src/app/(dashboard)/projects/[id]/budget/page.tsx`）
- [X] T151 [US3] 赤字時の警告表示実装
- [ ] T152 [P] [US3] 収支管理E2Eテスト（`frontend/tests/e2e/budget.spec.ts`）

**Checkpoint**: 収支管理が完全に機能し、利益が自動計算される。MVP完成！

---

## Phase 6: MVP統合テストとリファクタリング (Week 7)

**Purpose**: MVP（P1機能）の品質保証と最適化

- [ ] T153 全APIエンドポイントの統合テスト実行
- [ ] T154 全E2Eテストシナリオ実行
- [ ] T155 [P] パフォーマンステスト（APIレスポンスタイム）
- [ ] T156 [P] セキュリティスキャン（npm audit, go mod verify）
- [ ] T157 コードレビューとリファクタリング
- [ ] T158 [P] ドキュメント更新（README、API仕様）
- [ ] T159 [P] 初期データシード作成（`backend/scripts/seed.go`）

**Checkpoint**: MVP品質基準を満たし、デプロイ準備完了

---

## Phase 7: User Story 4 - ダッシュボード (Priority: P2) (Week 8-9)

**Goal**: 複数プロジェクトの横断的な管理とKPI表示

**Independent Test**: 複数プロジェクトを作成し、ダッシュボードで全体サマリーを確認

### Backend: Dashboard API

- [ ] T160 [US4] DashboardService実装（`backend/internal/service/dashboard_service.go`）
- [ ] T161 [US4] DashboardHandler実装（`backend/internal/handler/dashboard_handler.go`）
- [ ] T162 [US4] ダッシュボードサマリーエンドポイント（`GET /api/v1/dashboard`）
- [ ] T163 [US4] プロジェクト一覧の拡張フィルタ実装（ステータス、期間、利益率）
- [ ] T164 [US4] KPI集計クエリ最適化（インデックス確認）
- [ ] T165 [P] [US4] Dashboard API統合テスト

### Frontend: Dashboard UI

- [ ] T166 [US4] ダッシュボードフック（`frontend/src/hooks/use-dashboard.ts`）
- [ ] T167 [US4] KPIカードコンポーネント（`frontend/src/components/dashboard/kpi-card.tsx`）
- [ ] T168 [US4] プロジェクトテーブルコンポーネント（拡張版、ソート・フィルタ機能）
- [ ] T169 [US4] ダッシュボードページ（`frontend/src/app/(dashboard)/page.tsx`）
- [ ] T170 [P] [US4] ダッシュボードE2Eテスト

**Checkpoint**: ダッシュボードが動作し、複数プロジェクトの管理が容易に

---

## Phase 8: User Story 5 - グラフ可視化 (Priority: P2) (Week 9-10)

**Goal**: 予実、収支、推移をグラフで視覚化

**Independent Test**: プロジェクトデータからグラフが正しく生成され、インタラクティブに操作可能

### Backend: Analytics API

- [ ] T171 [US5] AnalyticsService実装（`backend/internal/service/analytics_service.go`）
- [ ] T172 [US5] AnalyticsHandler実装（`backend/internal/handler/analytics_handler.go`）
- [ ] T173 [US5] 予実比較データAPI（`GET /api/v1/projects/{id}/analytics/plan-actual`）
- [ ] T174 [US5] 収支データAPI（`GET /api/v1/projects/{id}/analytics/budget`）
- [ ] T175 [US5] 月次推移データAPI（`GET /api/v1/projects/{id}/analytics/trends`）
- [ ] T176 [US5] タスク別工数割合API（`GET /api/v1/projects/{id}/analytics/task-distribution`）
- [ ] T177 [US5] 複数プロジェクト比較API（`GET /api/v1/analytics/projects-comparison`）
- [ ] T178 [P] [US5] Analytics API統合テスト

### Frontend: Charts & Analytics UI

- [ ] T179 [US5] Rechartsセットアップと共通スタイル設定
- [ ] T180 [US5] 予実比較グラフコンポーネント（`frontend/src/components/charts/plan-actual-chart.tsx`）
- [ ] T181 [US5] 収支グラフコンポーネント（`frontend/src/components/charts/budget-chart.tsx`）
- [ ] T182 [US5] 推移グラフコンポーネント（`frontend/src/components/charts/trend-chart.tsx`）
- [ ] T183 [US5] 円グラフコンポーネント（`frontend/src/components/charts/pie-chart.tsx`）
- [ ] T184 [US5] グラフ共通機能（ツールチップ、レスポンシブ、カラーテーマ）
- [ ] T185 [US5] アナリティクスページ（`frontend/src/app/(dashboard)/projects/[id]/analytics/page.tsx`）
- [ ] T186 [US5] グラフパフォーマンステスト（100データポイント）
- [ ] T187 [P] [US5] グラフ表示E2Eテスト

**Checkpoint**: 全グラフが正常に表示され、データが正確に反映される

---

## Phase 9: User Story 6 - メンバー管理 (Priority: P3) (Week 11-12)

**Goal**: メンバーの詳細管理とリソース稼働レポート

**Independent Test**: メンバーを登録し、プロジェクトに割り当て、稼働レポートを確認

### Backend: Member Management Extension

- [ ] T188 [US6] MemberService拡張（稼働率計算、レポート生成）
- [ ] T189 [US6] メンバー詳細取得API（`GET /api/v1/members/{id}`）
- [ ] T190 [US6] プロジェクトメンバー割り当てAPI（`POST /api/v1/projects/{id}/members`）
- [ ] T191 [US6] メンバー別稼働レポートAPI（`GET /api/v1/members/{id}/workload`）
- [ ] T192 [US6] リソース稼働率API（`GET /api/v1/analytics/resource-utilization`）
- [ ] T193 [P] [US6] Member拡張機能統合テスト

### Frontend: Member Management UI

- [ ] T194 [US6] メンバー一覧ページ（`frontend/src/app/(dashboard)/members/page.tsx`）
- [ ] T195 [US6] メンバー詳細ページ（`frontend/src/app/(dashboard)/members/[id]/page.tsx`）
- [ ] T196 [US6] メンバーフォームコンポーネント（登録・編集）
- [ ] T197 [US6] プロジェクトメンバー割り当てコンポーネント
- [ ] T198 [US6] メンバー稼働レポートコンポーネント
- [ ] T199 [P] [US6] メンバー管理E2Eテスト

**Checkpoint**: メンバー管理が完全に機能し、稼働状況が可視化される

---

## Phase 10: User Story 7 - データエクスポート (Priority: P3) (Week 12-13)

**Goal**: プロジェクトデータのCSV/PDFエクスポート

**Independent Test**: データをエクスポートし、ファイルに正しく出力される

### Backend: Export API

- [ ] T200 [US7] ExportService実装（`backend/internal/service/export_service.go`）- CSV生成
- [ ] T201 [US7] ExportHandler実装（`backend/internal/handler/export_handler.go`）
- [ ] T202 [US7] プロジェクトCSVエクスポートAPI（`GET /api/v1/projects/{id}/export/csv`）
- [ ] T203 [US7] タスクCSVエクスポートAPI（`GET /api/v1/projects/{id}/tasks/export/csv`）
- [ ] T204 [US7] 収支CSVエクスポートAPI（`GET /api/v1/projects/{id}/budget/export/csv`）
- [ ] T205 [US7] 全プロジェクトサマリーCSVエクスポート（`GET /api/v1/projects/export/csv`）
- [ ] T206 [P] [US7] Export API統合テスト

### Frontend: Export UI

- [ ] T207 [US7] エクスポートボタンコンポーネント（各画面に配置）
- [ ] T208 [US7] CSV/PDFダウンロード機能実装
- [ ] T209 [US7] エクスポート進行状況表示
- [ ] T210 [P] [US7] エクスポート機能E2Eテスト

**Checkpoint**: データエクスポートが正常に動作し、全機能が完成

---

## Phase 11: パフォーマンス最適化とセキュリティ (Week 14-15)

**Purpose**: 本番環境準備とシステム全体の最適化

### Performance Optimization

- [ ] T211 フロントエンドコード分割（dynamic import）
- [ ] T212 [P] Next.js Image最適化
- [ ] T213 [P] バンドルサイズ分析と削減
- [ ] T214 バックエンドN+1クエリ解消
- [ ] T215 [P] データベースインデックス最適化
- [ ] T216 [P] APIレスポンスキャッシング実装
- [ ] T217 パフォーマンステスト実施（k6 or Apache Bench）
- [ ] T218 Success Criteria達成確認（API < 200ms, Dashboard < 2s, Graph < 3s）

### Security Hardening

- [ ] T219 セキュリティ監査実施（OWASP ZAP）
- [ ] T220 [P] 脆弱性スキャン（npm audit, go mod verify）
- [ ] T221 [P] HTTPS設定確認（本番環境）
- [ ] T222 [P] CORS設定最終確認
- [ ] T223 [P] CSP (Content Security Policy) 設定
- [ ] T224 レート制限実装（必要に応じて）
- [ ] T225 [P] 環境変数とシークレット管理確認

### Production Readiness

- [ ] T226 本番用Dockerイメージ最適化（マルチステージビルド）
- [ ] T227 [P] 本番用docker-compose.prod.yml作成
- [ ] T228 [P] データベースバックアップ設定
- [ ] T229 [P] ログ集約設定（CloudWatch, Datadog等）
- [ ] T230 [P] モニタリング・APM設定
- [ ] T231 [P] アラート設定（エラー率、レスポンスタイム）
- [ ] T232 デプロイ手順書作成

**Checkpoint**: 本番環境にデプロイ可能な状態

---

## Phase 12: ドキュメント整備とリリース準備 (Week 15-16)

**Purpose**: ドキュメント完成とリリース

- [ ] T233 API仕様書最終更新（OpenAPI Spec）
- [ ] T234 [P] ユーザーマニュアル作成（日本語）
- [ ] T235 [P] 運用マニュアル作成（バックアップ、復旧手順等）
- [ ] T236 [P] デプロイ手順書更新
- [ ] T237 [P] トラブルシューティングガイド作成
- [ ] T238 [P] 開発者向けオンボーディングドキュメント
- [ ] T239 README.md最終更新
- [ ] T240 [P] CHANGELOG.md作成
- [ ] T241 全E2Eテスト最終実行
- [ ] T242 ユーザー受け入れテスト（UAT）実施
- [ ] T243 本番環境デプロイ
- [ ] T244 本番環境動作確認
- [ ] T245 リリースアナウンス

**Checkpoint**: 本番リリース完了！

---

## Dependencies & Execution Order

### Phase Dependencies

1. **Phase 0 (環境構築)**: 依存なし、即開始可能
2. **Phase 1 (設計)**: Phase 0完了後
3. **Phase 2 (基盤)**: Phase 1完了後 - **全ユーザーストーリーをブロック**
4. **Phase 3-5 (P1機能/MVP)**: Phase 2完了後、順次実装
5. **Phase 6 (MVP統合)**: Phase 3-5完了後
6. **Phase 7-8 (P2機能)**: Phase 6完了後、並列実装可能
7. **Phase 9-10 (P3機能)**: Phase 6完了後、並列実装可能
8. **Phase 11-12 (最適化・リリース)**: 全機能完了後

### User Story Dependencies

- **US1 (プロジェクト管理)**: Phase 2完了後、依存なし
- **US2 (工数予実)**: Phase 2完了後、US1完了推奨（プロジェクト必須）
- **US3 (収支管理)**: Phase 2完了後、US1・US2完了推奨
- **US4 (ダッシュボード)**: US1-3完了後（MVPベース）
- **US5 (グラフ)**: US1-3完了後（データ必須）
- **US6 (メンバー管理)**: Phase 2完了後、独立実装可能
- **US7 (エクスポート)**: 各機能完了後

### Parallel Opportunities

- Phase 0: T002, T003, T005-T007, T009-T010, T012-T016, T018-T019, T022-T024 並列可能
- Phase 1: T032-T037（マイグレーション）、T039-T044（API定義）並列可能
- Phase 2: T053-T057（モデル）、T072-T074（共通UI）並列可能
- Backend開発者とFrontend開発者は並列作業可能（API契約確定後）
- 異なるユーザーストーリーは並列実装可能（Phase 2完了後）

---

## Implementation Strategy

### 推奨アプローチ: MVP First

1. **Week 1**: Phase 0完了 → 開発環境動作
2. **Week 2-3**: Phase 1完了 → 設計確定
3. **Week 3-4**: Phase 2完了 → 基盤完成（重要マイルストーン）
4. **Week 4-7**: Phase 3-5完了 → **MVP完成**（最重要マイルストーン）
5. **Week 7**: Phase 6完了 → MVP品質保証、デモ・フィードバック取得
6. **Week 8-10**: Phase 7-8完了 → P2機能追加
7. **Week 11-13**: Phase 9-10完了 → P3機能追加
8. **Week 14-16**: Phase 11-12完了 → 本番リリース

### チーム分担例（2-3名）

**1名体制**:
- Phase順に実装（Backend → Frontend）

**2名体制**:
- Developer A: Backend担当（Go/Echo/GORM）
- Developer B: Frontend担当（Next.js/TypeScript）
- Phase 2完了後、API契約ベースで並列開発

**3名体制**:
- Developer A: Backend + 基盤
- Developer B: Frontend + UI/UX
- Developer C: テスト + DevOps + ドキュメント

---

## Success Metrics（再掲）

### 実装完了の定義

✅ **機能**: 全24の機能要件（FR-001〜FR-024）実装完了  
✅ **ユーザーストーリー**: P1-P3の全Acceptance Scenarios満たす  
✅ **テスト**: Backend 70%+, Frontend 60%+ カバレッジ、全E2Eテストパス  
✅ **パフォーマンス**: API < 200ms, Dashboard < 2s, Graph < 3s  
✅ **セキュリティ**: 既知脆弱性0件、OWASP基準準拠  
✅ **ドキュメント**: 完全な運用・開発ドキュメント  
✅ **デプロイ**: 本番環境で正常動作

---

## Notes

- タスクは具体的なファイルパスを含む
- [P]タスクは並列実行可能（異なるファイル、依存なし）
- [US#]ラベルでユーザーストーリーとの紐付けを明示
- 各フェーズのCheckpointで進捗確認と品質検証
- MVP（Phase 3-6）を最優先で完成させる
- コミットは各タスクまたは論理的なグループ単位で実施
- 不明点はQuestions for Clarificationを参照し、早期に解決

**Status**: Draft  
**Last Updated**: 2025-10-11  
**Total Tasks**: 245
