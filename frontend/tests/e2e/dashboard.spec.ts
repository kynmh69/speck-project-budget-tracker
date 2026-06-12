import { test, expect, Page } from '@playwright/test';

// APIをルートモックしたダッシュボードのE2Eテスト (T170)
// バックエンド・DBなしでフロントエンドの動作を検証する

const user = {
  id: '11111111-1111-1111-1111-111111111111',
  email: 'demo@example.com',
  name: 'デモユーザー',
  role: 'member',
};

const now = '2026-06-01T00:00:00Z';

const projects = [
  {
    id: 'aaaaaaaa-0000-0000-0000-000000000001',
    user_id: user.id,
    name: 'ECサイトリニューアル',
    status: 'in_progress',
    start_date: '2026-04-01',
    end_date: '2026-08-31',
    budget_amount: 5000000,
    created_at: now,
    updated_at: now,
  },
  {
    id: 'aaaaaaaa-0000-0000-0000-000000000002',
    user_id: user.id,
    name: '社内システム保守',
    status: 'completed',
    created_at: now,
    updated_at: now,
  },
  {
    id: 'aaaaaaaa-0000-0000-0000-000000000003',
    user_id: user.id,
    name: 'モバイルアプリ開発',
    status: 'planning',
    created_at: now,
    updated_at: now,
  },
];

const dashboard = {
  total_projects: 3,
  active_projects: 1,
  completed_projects: 1,
  total_revenue: 3000000,
  total_profit: 1400000,
  average_profit_rate: 45.0,
  recent_projects: projects,
};

async function mockApi(page: Page) {
  await page.route('**/api/v1/auth/me', (route) =>
    route.fulfill({ json: { success: true, data: user } })
  );
  await page.route('**/api/v1/dashboard', (route) =>
    route.fulfill({ json: { success: true, data: dashboard } })
  );
  await page.route('**/api/v1/projects**', (route) => {
    const url = new URL(route.request().url());
    const status = url.searchParams.get('status');
    const filtered = status ? projects.filter((p) => p.status === status) : projects;
    return route.fulfill({
      json: {
        success: true,
        data: {
          projects: filtered,
          pagination: { page: 1, per_page: 10, total: filtered.length, total_pages: 1 },
        },
      },
    });
  });
}

// Next.jsミドルウェアはCookie、APIクライアントとZustandストアは
// localStorageを参照するため両方を設定する
async function authenticate(page: Page) {
  await page.context().addCookies([
    { name: 'token', value: 'e2e-test-token', url: 'http://localhost:3000' },
  ]);
  await page.addInitScript(
    ([u]) => {
      localStorage.setItem('token', 'e2e-test-token');
      localStorage.setItem(
        'auth-storage',
        JSON.stringify({ state: { user: u, isAuthenticated: true }, version: 0 })
      );
    },
    [user]
  );
}

test.describe('ダッシュボード', () => {
  test('KPIカードにサマリーが表示される', async ({ page }) => {
    await authenticate(page);
    await mockApi(page);

    await page.goto('/dashboard');

    await expect(page.getByRole('heading', { name: 'ダッシュボード' })).toBeVisible();
    await expect(page.getByText('ようこそ、デモユーザーさん')).toBeVisible();

    await expect(page.getByText('総プロジェクト数')).toBeVisible();
    await expect(page.getByText('総売上')).toBeVisible();
    await expect(page.getByText('総利益')).toBeVisible();
    await expect(page.getByText('平均利益率')).toBeVisible();
    await expect(page.getByText('45.0%')).toBeVisible();
    // 通貨表示（記号はロケール依存のため数値部分のみ検証）
    await expect(page.getByText(/3,000,000/)).toBeVisible();
    await expect(page.getByText(/1,400,000/)).toBeVisible();
  });

  test('プロジェクト一覧テーブルが表示される', async ({ page }) => {
    await authenticate(page);
    await mockApi(page);

    await page.goto('/dashboard');

    await expect(page.getByRole('heading', { name: 'プロジェクト一覧' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'ECサイトリニューアル' })).toBeVisible();
    await expect(page.getByRole('link', { name: '社内システム保守' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'モバイルアプリ開発' })).toBeVisible();
    await expect(page.getByText('全3件')).toBeVisible();
  });

  test('ステータスフィルタで絞り込める', async ({ page }) => {
    await authenticate(page);
    await mockApi(page);

    await page.goto('/dashboard');
    await expect(page.getByRole('link', { name: 'ECサイトリニューアル' })).toBeVisible();

    // ステータスフィルタで「進行中」を選択
    await page.getByRole('combobox').click();
    await page.getByRole('option', { name: '進行中' }).click();

    await expect(page.getByRole('link', { name: 'ECサイトリニューアル' })).toBeVisible();
    await expect(page.getByRole('link', { name: '社内システム保守' })).not.toBeVisible();
    await expect(page.getByText('全1件')).toBeVisible();
  });

  test('未認証アクセスはログインページへリダイレクトされる', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveURL(/\/login/);
  });
});
