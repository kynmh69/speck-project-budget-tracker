# Frontend Design: プロジェクト予算管理システム

**Date**: 2025-10-11  
**Phase**: Phase 1 - Frontend Design

## 画面構成

### ルーティング設計（Next.js App Router）

```
app/
├── (auth)/                    # 認証グループ
│   ├── login/
│   │   └── page.tsx          # ログイン画面
│   └── register/
│       └── page.tsx          # ユーザー登録画面
│
├── (dashboard)/               # ダッシュボードグループ（認証必須）
│   ├── layout.tsx            # 共通レイアウト（Header + Sidebar）
│   ├── page.tsx              # ダッシュボード
│   ├── projects/
│   │   ├── page.tsx          # プロジェクト一覧
│   │   ├── new/
│   │   │   └── page.tsx      # プロジェクト新規作成
│   │   └── [id]/
│   │       ├── page.tsx      # プロジェクト詳細（概要タブ）
│   │       ├── tasks/
│   │       │   └── page.tsx  # タスク管理タブ
│   │       ├── budget/
│   │       │   └── page.tsx  # 収支管理タブ
│   │       └── analytics/
│   │           └── page.tsx  # グラフ分析タブ
│   └── members/
│       ├── page.tsx          # メンバー一覧
│       ├── new/
│       │   └── page.tsx      # メンバー新規作成
│       └── [id]/
│           └── page.tsx      # メンバー詳細
│
├── layout.tsx                # ルートレイアウト
└── page.tsx                  # ランディングページ
```

## ワイヤーフレーム

### 1. ログイン画面 (`/login`)

```
┌─────────────────────────────────────────────┐
│                                             │
│     プロジェクト予算管理システム              │
│                                             │
│  ┌───────────────────────────────────────┐  │
│  │  メールアドレス                       │  │
│  │  [________________________]          │  │
│  │                                     │  │
│  │  パスワード                          │  │
│  │  [________________________]          │  │
│  │                                     │  │
│  │  [  ログイン  ]                     │  │
│  │                                     │  │
│  │  アカウントをお持ちでない方           │  │
│  │  [ 新規登録 ]                       │  │
│  └───────────────────────────────────────┘  │
│                                             │
└─────────────────────────────────────────────┘
```

### 2. ダッシュボード (`/dashboard`)

```
┌─────────────────────────────────────────────────────┐
│ Header: Logo | ダッシュボード | 通知 | ユーザー  │
├───────┬─────────────────────────────────────────────┤
│       │  KPIカード                                  │
│       │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐      │
│       │  │総PJ数│ │進行中│ │総利益│ │利益率│      │
│       │  │  10  │ │  5   │ │500万│ │ 25% │      │
│ Side  │  └──────┘ └──────┘ └──────┘ └──────┘      │
│ bar   │                                             │
│       │  プロジェクト一覧                            │
│ ・PJ   │  ┌───────────────────────────────────────┐ │
│ ・メン │  │名前 │ ステータス │ 利益率 │ 更新日 │操作││
│ ・設定 │  ├───────────────────────────────────────┤ │
│       │  │PJ A │ 進行中    │  30%  │2025-10-10│詳細││
│       │  │PJ B │ 進行中    │  20%  │2025-10-09│詳細││
│       │  │PJ C │ 完了      │  35%  │2025-10-05│詳細││
│       │  └───────────────────────────────────────┘ │
└───────┴─────────────────────────────────────────────┘
```

### 3. プロジェクト一覧 (`/projects`)

```
┌─────────────────────────────────────────────────────┐
│ Header                                              │
├───────┬─────────────────────────────────────────────┤
│       │  プロジェクト一覧                            │
│       │                                             │
│       │  [ 新規作成 ]     [検索: ____] [フィルタ▼] │
│       │                                             │
│ Side  │  ┌────────────────────────────────────────┐│
│ bar   │  │ プロジェクトA                          ││
│       │  │ ステータス: 進行中                     ││
│       │  │ 予算: 500万円 | 利益率: 25%           ││
│       │  │ [ 詳細 ]                              ││
│       │  └────────────────────────────────────────┘│
│       │  ┌────────────────────────────────────────┐│
│       │  │ プロジェクトB                          ││
│       │  │ ステータス: 計画中                     ││
│       │  │ 予算: 300万円 | 利益率: --            ││
│       │  │ [ 詳細 ]                              ││
│       │  └────────────────────────────────────────┘│
│       │                                             │
│       │  ページネーション: ◀ 1 2 3 ▶               │
└───────┴─────────────────────────────────────────────┘
```

### 4. プロジェクト詳細 - タスク管理タブ (`/projects/[id]/tasks`)

```
┌─────────────────────────────────────────────────────┐
│ Header: プロジェクトA                                │
├───────┬─────────────────────────────────────────────┤
│       │  [ 概要 | タスク | 収支 | グラフ ]           │
│       │                                             │
│       │  タスク管理              [ + 新規タスク ]   │
│       │                                             │
│ Side  │  予実サマリー                               │
│ bar   │  ┌────────────────────────────────────────┐│
│       │  │ 予定工数:  100h                        ││
│       │  │ 実績工数:   80h                        ││
│       │  │ 差異:      -20h (達成率: 80%)          ││
│       │  └────────────────────────────────────────┘│
│       │                                             │
│       │  タスク一覧                                 │
│       │  ┌─────────────────────────────────────┐  │
│       │  │タスク名│予定h│実績h│差異│ステータス│操作││
│       │  ├─────────────────────────────────────┤  │
│       │  │機能A開発│ 30 │ 35 │+5 │完了      │編集││
│       │  │機能B開発│ 40 │ 25 │-15│進行中    │編集││
│       │  │テスト   │ 30 │ 20 │-10│進行中    │編集││
│       │  └─────────────────────────────────────┘  │
└───────┴─────────────────────────────────────────────┘
```

### 5. プロジェクト詳細 - 収支管理タブ (`/projects/[id]/budget`)

```
┌─────────────────────────────────────────────────────┐
│ Header: プロジェクトA                                │
├───────┬─────────────────────────────────────────────┤
│       │  [ 概要 | タスク | 収支 | グラフ ]           │
│       │                                             │
│       │  収支管理                                   │
│       │                                             │
│ Side  │  ┌────────────────────────────────────────┐│
│ bar   │  │ 収支サマリー                           ││
│       │  │                                        ││
│       │  │ 売上:      5,000,000円                ││
│       │  │ 総コスト:  3,750,000円                ││
│       │  │ ────────────────────                  ││
│       │  │ 利益:      1,250,000円  [編集]       ││
│       │  │ 利益率:    25.0%                      ││
│       │  │                                        ││
│       │  └────────────────────────────────────────┘│
│       │                                             │
│       │  コスト内訳（メンバー別）                    │
│       │  ┌────────────────────────────────────┐  │
│       │  │メンバー名 │工数 │単価  │コスト      │  │
│       │  ├────────────────────────────────────┤  │
│       │  │山田太郎   │100h│5,000│ 500,000円  │  │
│       │  │佐藤花子   │ 80h│6,000│ 480,000円  │  │
│       │  │鈴木一郎   │120h│4,500│ 540,000円  │  │
│       │  └────────────────────────────────────┘  │
└───────┴─────────────────────────────────────────────┘
```

### 6. プロジェクト詳細 - グラフ分析タブ (`/projects/[id]/analytics`)

```
┌─────────────────────────────────────────────────────┐
│ Header: プロジェクトA                                │
├───────┬─────────────────────────────────────────────┤
│       │  [ 概要 | タスク | 収支 | グラフ ]           │
│       │                                             │
│       │  分析・グラフ                               │
│       │                                             │
│ Side  │  予実比較                                   │
│ bar   │  ┌────────────────────────────────────────┐│
│       │  │                    ┃                   ││
│       │  │   予定 ■  実績 ■   ┃                   ││
│       │  │      ▆▆    ▆▆      ┃    ▆▆            ││
│       │  │    ▆▆▆▆  ▆▆▆▆      ┃  ▆▆▆▆            ││
│       │  │  ▆▆▆▆▆▆▆▆▆▆▆▆    ▆▆┃▆▆▆▆▆▆            ││
│       │  │  機能A  機能B     テスト                ││
│       │  └────────────────────────────────────────┘│
│       │                                             │
│       │  タスク別工数割合                            │
│       │  ┌────────────────────────────────────────┐│
│       │  │         ╱╲                             ││
│       │  │       ╱    ╲       機能A: 35%         ││
│       │  │     ╱ 機能A ╲      機能B: 40%         ││
│       │  │    ╱─────────╲     テスト: 25%        ││
│       │  │   │  テスト   │機能B                   ││
│       │  │    ╲_________╱                         ││
│       │  └────────────────────────────────────────┘│
└───────┴─────────────────────────────────────────────┘
```

## コンポーネント設計

### レイアウトコンポーネント

1. **Header** (`components/layout/header.tsx`)
   - ロゴ
   - ナビゲーション
   - 通知アイコン
   - ユーザーメニュー

2. **Sidebar** (`components/layout/sidebar.tsx`)
   - メインナビゲーション
   - 折りたたみ可能

3. **DashboardLayout** (`app/(dashboard)/layout.tsx`)
   - Header + Sidebar + コンテンツエリア

### UIコンポーネント（shadcn/ui）

優先的に導入するコンポーネント：

1. **Form系**
   - Button
   - Input
   - Select
   - Textarea
   - DatePicker
   - Form (react-hook-form統合)

2. **データ表示系**
   - Table
   - Card
   - Badge
   - Progress

3. **フィードバック系**
   - Dialog/Modal
   - Toast (通知)
   - Alert
   - Loading Spinner

4. **ナビゲーション系**
   - Tabs
   - Dropdown Menu
   - Pagination

### 機能別コンポーネント

#### プロジェクト関連 (`components/projects/`)

- `project-list.tsx` - プロジェクト一覧
- `project-card.tsx` - プロジェクトカード
- `project-form.tsx` - プロジェクト作成・編集フォーム
- `project-filters.tsx` - フィルタ・検索UI

#### タスク関連 (`components/tasks/`)

- `task-list.tsx` - タスク一覧テーブル
- `task-item.tsx` - タスク行（予実表示）
- `task-form.tsx` - タスク作成・編集フォーム
- `plan-actual-summary.tsx` - 予実サマリー表示

#### 収支関連 (`components/budget/`)

- `budget-summary.tsx` - 収支サマリーカード
- `cost-breakdown.tsx` - コスト内訳テーブル
- `revenue-form.tsx` - 売上入力フォーム

#### グラフ関連 (`components/charts/`)

- `plan-actual-chart.tsx` - 予実比較グラフ（棒グラフ）
- `budget-chart.tsx` - 収支グラフ（棒グラフ）
- `trend-chart.tsx` - 推移グラフ（折れ線グラフ）
- `pie-chart.tsx` - 円グラフ（タスク別工数割合）

#### ダッシュボード関連 (`components/dashboard/`)

- `kpi-card.tsx` - KPI表示カード
- `recent-projects.tsx` - 最近のプロジェクト一覧
- `stats-overview.tsx` - 統計概要

#### 共通コンポーネント (`components/common/`)

- `loading.tsx` - ローディングスピナー
- `error-boundary.tsx` - エラー表示
- `pagination.tsx` - ページネーション
- `empty-state.tsx` - データなし状態の表示

## 状態管理戦略

### TanStack Query（サーバー状態）

APIデータの取得、キャッシング、同期を管理：

```typescript
// hooks/use-projects.ts
export function useProjects(params: ProjectListParams) {
  return useQuery({
    queryKey: ['projects', params],
    queryFn: () => projectApi.getProjects(params),
  });
}

export function useProject(id: string) {
  return useQuery({
    queryKey: ['projects', id],
    queryFn: () => projectApi.getProject(id),
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: projectApi.createProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}
```

### Zustand（クライアント状態）

UI状態、認証状態などのグローバル状態：

```typescript
// store/auth-store.ts
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,
  login: async (email, password) => {
    const { user, token } = await authApi.login(email, password);
    set({ user, token, isAuthenticated: true });
  },
  logout: () => {
    set({ user: null, token: null, isAuthenticated: false });
  },
}));

// store/ui-store.ts
interface UIState {
  sidebarOpen: boolean;
  toggleSidebar: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  sidebarOpen: true,
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
}));
```

## バリデーション（Zod）

フォーム入力のバリデーション：

```typescript
// schemas/project-schema.ts
export const projectSchema = z.object({
  name: z.string().min(1, '名前は必須です').max(200),
  description: z.string().optional(),
  status: z.enum(['planning', 'in_progress', 'completed', 'on_hold']),
  budget_amount: z.number().min(0).optional(),
  start_date: z.string().optional(),
  end_date: z.string().optional(),
});

export type ProjectFormData = z.infer<typeof projectSchema>;

// schemas/task-schema.ts
export const taskSchema = z.object({
  name: z.string().min(1, 'タスク名は必須です'),
  description: z.string().optional(),
  planned_hours: z.number().min(0, '0以上の値を入力してください'),
  assigned_to: z.string().uuid().optional(),
  status: z.enum(['todo', 'in_progress', 'completed', 'blocked']),
});
```

## APIクライアント

```typescript
// lib/api-client.ts
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// トークン自動付与
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// エラーハンドリング
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 認証エラー時は自動ログアウト
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

## レスポンシブ対応

- **デスクトップファースト**: 主にPC利用を想定
- **ブレークポイント**:
  - Desktop: 1024px以上
  - Tablet: 768px - 1023px
  - Mobile: 767px以下（基本機能のみ対応）

## アクセシビリティ

- セマンティックHTML
- キーボードナビゲーション対応
- ARIA属性の適切な使用
- カラーコントラスト（WCAG AA準拠）

## パフォーマンス最適化

1. **コード分割**
   - 動的インポート（`dynamic(() => import())` ）
   - ページ単位での分割

2. **画像最適化**
   - Next.js Image コンポーネント使用
   - WebP形式

3. **キャッシング**
   - TanStack Query による自動キャッシング
   - staleTime, cacheTime の適切な設定

4. **遅延ロード**
   - グラフコンポーネントの遅延ロード
   - スクロール時のデータロード

---

**Status**: Draft  
**Next**: 認証設計、実装開始  
**Date**: 2025-10-11
