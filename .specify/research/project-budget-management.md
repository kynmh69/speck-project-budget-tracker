# Technical Research: プロジェクト予算管理システム

**Date**: 2025-10-11  
**Phase**: Phase 0 - 環境構築とリサーチ

## 概要

このドキュメントは Phase 0 における技術スタックの検証結果とベストプラクティスをまとめたものです。

## 技術スタック検証結果

### Backend: Go + Echo + GORM

#### ✅ 検証済み項目

1. **Echo v4 セットアップ**
   - ステータス: 成功
   - バージョン: v4.13.4
   - 基本的なルーティング、ミドルウェア設定が正常動作
   - `/health` エンドポイントが正常応答

2. **GORM v2**
   - ステータス: インストール完了
   - バージョン: v1.31.0
   - PostgreSQL ドライバ: v1.6.0

3. **依存関係**
   - ✅ github.com/labstack/echo/v4
   - ✅ gorm.io/gorm
   - ✅ gorm.io/driver/postgres
   - ✅ github.com/golang-jwt/jwt/v5
   - ✅ github.com/go-playground/validator/v10
   - ✅ golang.org/x/crypto (bcrypt)

#### 推奨事項

1. **ディレクトリ構造**
   ```
   backend/
   ├── cmd/server/          # エントリーポイント
   ├── internal/
   │   ├── config/          # 設定管理
   │   ├── models/          # GORMモデル
   │   ├── repository/      # データアクセス層
   │   ├── service/         # ビジネスロジック
   │   ├── handler/         # HTTPハンドラー
   │   ├── middleware/      # ミドルウェア
   │   ├── dto/             # データ転送オブジェクト
   │   └── validator/       # カスタムバリデーター
   ├── migrations/          # DBマイグレーション
   └── tests/               # テスト
   ```

2. **環境変数管理**
   - `.env.example` を提供
   - 本番環境では環境変数またはシークレット管理サービスを使用

3. **エラーハンドリング**
   - Echo のカスタムエラーハンドラーを実装
   - 統一されたエラーレスポンス形式

4. **ロギング**
   - Echo の Logger ミドルウェアを使用
   - 構造化ログ (structured logging) の導入推奨

### Frontend: Next.js 14 + shadcn/ui

#### ✅ 検証済み項目

1. **Next.js 14 (App Router)**
   - ステータス: セットアップ完了
   - TypeScript設定完了
   - App Router構造採用

2. **Tailwind CSS**
   - ステータス: 設定完了
   - shadcn/ui用のカスタムテーマ設定

3. **依存関係（計画）**
   - ✅ next ^14.2.0
   - ✅ react ^18.3.0
   - ✅ typescript ^5.5.0
   - ✅ tailwindcss ^3.4.0
   - 📝 @tanstack/react-query (インストール予定)
   - 📝 zustand (インストール予定)
   - 📝 zod (インストール予定)
   - 📝 recharts (インストール予定)

#### 推奨事項

1. **ディレクトリ構造**
   ```
   frontend/
   ├── app/                 # App Router
   │   ├── (auth)/         # 認証グループ
   │   ├── (dashboard)/    # ダッシュボードグループ
   │   ├── layout.tsx      # ルートレイアウト
   │   └── page.tsx        # ホームページ
   ├── components/
   │   ├── ui/             # shadcn/uiコンポーネント
   │   ├── projects/       # プロジェクト関連
   │   ├── tasks/          # タスク関連
   │   ├── charts/         # グラフコンポーネント
   │   └── layout/         # レイアウトコンポーネント
   ├── lib/                # ユーティリティ
   ├── hooks/              # カスタムフック
   ├── store/              # Zustand store
   ├── types/              # TypeScript型定義
   └── schemas/            # Zodスキーマ
   ```

2. **shadcn/ui セットアップ**
   - コンポーネントは必要に応じて個別インストール
   - `components.json` で設定管理
   - 以下のコンポーネントを優先的に導入:
     - Button, Card, Input, Form
     - Dialog, Table, Select
     - Toast (通知用)

3. **状態管理戦略**
   - **サーバー状態**: TanStack Query (React Query v5)
     - API呼び出し、キャッシング、再取得
     - Optimistic Updates
   - **クライアント状態**: Zustand
     - UI状態、認証状態
     - グローバルな状態管理

4. **APIクライアント**
   - axios または fetch wrapper
   - 認証トークンの自動付与
   - エラーハンドリングの統一

### Database: PostgreSQL

#### 推奨事項

1. **マイグレーション管理**
   - golang-migrate を使用
   - Up/Down マイグレーション
   - バージョン管理

2. **インデックス戦略**
   - 頻繁に検索されるカラムにインデックス
   - 複合インデックスの活用
   - 初期段階から設計

3. **パフォーマンス考慮**
   - コネクションプーリング
   - N+1クエリ問題の回避 (GORM Preload使用)
   - 適切なトランザクション管理

### Docker & Development Environment

#### ✅ 検証済み項目

1. **Docker Compose**
   - ステータス: 設定完了
   - サービス: PostgreSQL, Backend, Frontend
   - ヘルスチェック設定

2. **Makefile**
   - ステータス: 作成完了
   - 主要コマンド:
     - `make dev` - 開発環境起動
     - `make test` - テスト実行
     - `make migrate-up` - マイグレーション実行

#### 推奨事項

1. **ホットリロード**
   - Backend: Air（Go用ホットリロードツール）
   - Frontend: Next.js 組み込みのホットリロード

2. **ボリューム管理**
   - データベースデータの永続化
   - node_modules の適切なマウント

## 認証フロー設計

### JWT認証

#### 推奨実装

1. **トークン戦略**
   - Access Token: 短命 (15分-1時間)
   - Refresh Token: 長命 (7-30日) ※将来実装
   - HttpOnly Cookie または Authorization Header

2. **セキュリティ考慮事項**
   - パスワードのbcryptハッシュ化 (cost: 12)
   - JWT秘密鍵は環境変数管理 (最低32文字)
   - HTTPS通信必須 (本番環境)
   - CSRF対策

3. **実装パターン**
   ```go
   // Backend
   - AuthService: ユーザー登録、ログイン、トークン生成
   - AuthMiddleware: トークン検証、ユーザー情報の抽出
   - AuthHandler: エンドポイント実装
   ```

   ```typescript
   // Frontend
   - useAuth hook: 認証状態管理
   - auth-store: Zustand store
   - api-client: トークンの自動付与
   - middleware: 保護されたルートのガード
   ```

## グラフライブラリ: Recharts

### 選定理由

1. **React統合**: React専用で使いやすい
2. **宣言的API**: コンポーネントベースで直感的
3. **レスポンシブ対応**: 標準でレスポンシブ
4. **カスタマイズ性**: 十分な柔軟性

### 使用予定のグラフタイプ

1. **BarChart**: 予実比較、収支比較
2. **LineChart**: 月次推移
3. **PieChart**: タスク別工数割合
4. **ComposedChart**: 複合グラフ（必要に応じて）

### パフォーマンス考慮

- 100データポイントで3秒以内の描画目標
- 大量データの場合はサーバーサイドで集計
- 遅延ロード、仮想化の活用

## CI/CD戦略

### GitHub Actions

#### 設定済みワークフロー

1. **Backend CI**
   - Lint (golangci-lint)
   - Test (PostgreSQLサービスコンテナ使用)
   - Build
   - カバレッジ測定

2. **Frontend CI**
   - Lint (ESLint)
   - Type Check (TypeScript)
   - Test (Jest)
   - Build

#### 今後の拡張

1. **E2Eテスト** (Playwright)
2. **自動デプロイ** (CD)
3. **セキュリティスキャン**
4. **パフォーマンステスト**

## パフォーマンス目標と計測

### 目標値

- API応答時間: 平均 < 200ms, p95 < 500ms
- ダッシュボード表示: < 2秒 (50プロジェクト)
- グラフ描画: < 3秒 (100データポイント)

### 計測方法

1. **Backend**
   - Echo middleware でレスポンスタイム記録
   - APMツール (将来導入)

2. **Frontend**
   - Web Vitals (LCP, FID, CLS)
   - Lighthouse CI

3. **ロードテスト**
   - k6 または Apache Bench
   - 100同時接続でのテスト

## セキュリティベストプラクティス

### 実装必須項目

1. **認証・認可**
   - ✅ JWT認証
   - ✅ パスワードハッシュ化 (bcrypt)
   - 📝 ロールベースアクセス制御 (RBAC)

2. **入力検証**
   - ✅ バリデーション (go-playground/validator, Zod)
   - ✅ SQLインジェクション対策 (GORM ORM使用)
   - ✅ XSS対策 (React自動エスケープ)

3. **通信**
   - ✅ CORS設定
   - 📝 HTTPS必須 (本番環境)
   - 📝 CSP (Content Security Policy)

4. **レート制限**
   - 📝 API呼び出し制限
   - 📝 ログイン試行制限

## 技術的課題と対策

### 潜在的な課題

1. **N+1クエリ問題**
   - 対策: GORM Preload/Joins の適切な使用
   - 早期発見: クエリログの監視

2. **データ量増加時のパフォーマンス**
   - 対策: ページネーション必須
   - 対策: インデックス最適化
   - 対策: キャッシング戦略

3. **フロントエンドバンドルサイズ**
   - 対策: コード分割 (dynamic import)
   - 対策: Tree shaking
   - 対策: 依存関係の最小化

4. **同時編集時のデータ競合**
   - 対策: 楽観的ロック (updated_at チェック)
   - 対策: トランザクション管理

## 次のステップ (Phase 1)

1. **データモデル設計**
   - ER図作成
   - マイグレーションファイル作成
   - GORMモデル定義

2. **API契約設計**
   - OpenAPI 3.0 スキーマ作成
   - エンドポイント定義
   - リクエスト/レスポンス型定義

3. **フロントエンド詳細設計**
   - ワイヤーフレーム作成
   - コンポーネント設計
   - ルーティング設計

4. **認証設計**
   - JWT認証フロー詳細化
   - 権限モデル設計

## 結論

Phase 0 の環境構築とリサーチは成功裏に完了しました。選定した技術スタック（Go + Echo + GORM、Next.js + shadcn/ui）は要件を満たすことが確認されました。

### 主な成果

- ✅ リポジトリ構造の確立
- ✅ バックエンド基盤セットアップ
- ✅ フロントエンド基盤セットアップ
- ✅ Docker Compose 開発環境
- ✅ CI/CD パイプライン基礎
- ✅ Makefile コマンド統一化

### 準備完了

Phase 1 (アーキテクチャ設計とAPI契約) に進む準備が整いました。

---

**Status**: Complete  
**Next Phase**: Phase 1 - Architecture & Contracts  
**Date Completed**: 2025-10-11
