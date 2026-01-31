# Authentication & Authorization Design

**Date**: 2025-10-11  
**Phase**: Phase 1 - Authentication Design

## 認証方式

### JWT (JSON Web Token)

シンプルでステートレスな認証方式としてJWTを採用。

## 認証フロー

### 1. ユーザー登録フロー

```
Frontend                Backend              Database
   │                       │                     │
   │  POST /auth/register  │                     │
   ├──────────────────────>│                     │
   │  {email, password,    │                     │
   │   name}                │                     │
   │                       │  バリデーション      │
   │                       │                     │
   │                       │  パスワードハッシュ化│
   │                       │  (bcrypt, cost=12)  │
   │                       │                     │
   │                       │  INSERT users       │
   │                       ├────────────────────>│
   │                       │                     │
   │                       │  ユーザー作成完了    │
   │                       │<────────────────────┤
   │                       │                     │
   │                       │  JWT生成            │
   │                       │  (payload: user_id, │
   │                       │   email, role)      │
   │                       │                     │
   │  {token, user}        │                     │
   │<──────────────────────┤                     │
   │                       │                     │
   │  トークン保存          │                     │
   │  (localStorage)       │                     │
   │                       │                     │
```

### 2. ログインフロー

```
Frontend                Backend              Database
   │                       │                     │
   │  POST /auth/login     │                     │
   ├──────────────────────>│                     │
   │  {email, password}    │                     │
   │                       │                     │
   │                       │  SELECT user        │
   │                       │  WHERE email=?      │
   │                       ├────────────────────>│
   │                       │                     │
   │                       │  ユーザー情報取得    │
   │                       │<────────────────────┤
   │                       │                     │
   │                       │  パスワード検証      │
   │                       │  (bcrypt.Compare)   │
   │                       │                     │
   │                       │  JWT生成            │
   │                       │                     │
   │  {token, user}        │                     │
   │<──────────────────────┤                     │
   │                       │                     │
   │  トークン保存          │                     │
   │  (localStorage)       │                     │
   │                       │                     │
```

### 3. 保護されたAPIアクセスフロー

```
Frontend                Backend              Database
   │                       │                     │
   │  GET /api/v1/projects │                     │
   │  Authorization:       │                     │
   │  Bearer <token>       │                     │
   ├──────────────────────>│                     │
   │                       │                     │
   │                       │  トークン検証        │
   │                       │  (JWT signature)    │
   │                       │                     │
   │                       │  ユーザー情報抽出    │
   │                       │  (from JWT payload) │
   │                       │                     │
   │                       │  権限チェック        │
   │                       │  (role-based)       │
   │                       │                     │
   │                       │  SELECT projects    │
   │                       │  WHERE user_id=?    │
   │                       ├────────────────────>│
   │                       │                     │
   │                       │  プロジェクト一覧    │
   │                       │<────────────────────┤
   │                       │                     │
   │  {data, pagination}   │                     │
   │<──────────────────────┤                     │
   │                       │                     │
```

### 4. ログアウトフロー

```
Frontend                Backend
   │                       │
   │  ログアウト操作        │
   │                       │
   │  トークン削除          │
   │  (localStorage)       │
   │                       │
   │  リダイレクト          │
   │  -> /login            │
   │                       │
```

## JWT構成

### Payload

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "manager",
  "iat": 1234567890,
  "exp": 1234657890
}
```

### 設定

- **アルゴリズム**: HS256 (HMAC-SHA256)
- **有効期限**: 24時間（86400秒）
- **秘密鍵**: 環境変数で管理（最低32文字）

## 実装詳細

### Backend (Go)

#### 1. パスワードハッシュ化

```go
// internal/service/auth_service.go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### 2. JWT生成

```go
// internal/service/auth_service.go
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(user *models.User, secret string) (string, error) {
    claims := &Claims{
        UserID: user.ID.String(),
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

#### 3. JWT検証ミドルウェア

```go
// internal/middleware/auth.go
func AuthMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
            }
            
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            
            token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
                return []byte(secret), nil
            })
            
            if err != nil || !token.Valid {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
            }
            
            claims, ok := token.Claims.(*Claims)
            if !ok {
                return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
            }
            
            // コンテキストにユーザー情報を保存
            c.Set("user_id", claims.UserID)
            c.Set("email", claims.Email)
            c.Set("role", claims.Role)
            
            return next(c)
        }
    }
}
```

#### 4. AuthHandler

```go
// internal/handler/auth_handler.go
type AuthHandler struct {
    authService *service.AuthService
}

func (h *AuthHandler) Register(c echo.Context) error {
    var req dto.RegisterRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    
    if err := c.Validate(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    
    user, token, err := h.authService.Register(req.Email, req.Password, req.Name)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    
    return c.JSON(http.StatusCreated, dto.AuthResponse{
        Token: token,
        User:  user,
    })
}

func (h *AuthHandler) Login(c echo.Context) error {
    var req dto.LoginRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    
    user, token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
    }
    
    return c.JSON(http.StatusOK, dto.AuthResponse{
        Token: token,
        User:  user,
    })
}

func (h *AuthHandler) Me(c echo.Context) error {
    userID := c.Get("user_id").(string)
    user, err := h.authService.GetUserByID(userID)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "user not found")
    }
    
    return c.JSON(http.StatusOK, user)
}
```

### Frontend (TypeScript/React)

#### 1. Auth Store (Zustand)

```typescript
// store/auth-store.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => {
        localStorage.setItem('token', token);
        set({ user, token, isAuthenticated: true });
      },
      clearAuth: () => {
        localStorage.removeItem('token');
        set({ user: null, token: null, isAuthenticated: false });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

#### 2. API Client

```typescript
// lib/api-client.ts
import axios from 'axios';
import { useAuthStore } from '@/store/auth-store';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエストインターセプター: トークン自動付与
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// レスポンスインターセプター: エラーハンドリング
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 認証エラー時は自動ログアウト
      const { clearAuth } = useAuthStore.getState();
      clearAuth();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

#### 3. Auth API

```typescript
// lib/auth-api.ts
import apiClient from './api-client';

export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
}

export const authApi = {
  register: async (data: RegisterData): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/register', data);
    return response.data;
  },

  login: async (data: LoginData): Promise<AuthResponse> => {
    const response = await apiClient.post('/auth/login', data);
    return response.data;
  },

  me: async () => {
    const response = await apiClient.get('/auth/me');
    return response.data;
  },
};
```

#### 4. Auth Hooks

```typescript
// hooks/use-auth.ts
import { useMutation, useQuery } from '@tanstack/react-query';
import { authApi } from '@/lib/auth-api';
import { useAuthStore } from '@/store/auth-store';
import { useRouter } from 'next/navigation';

export function useLogin() {
  const { setAuth } = useAuthStore();
  const router = useRouter();

  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.user, data.token);
      router.push('/dashboard');
    },
  });
}

export function useRegister() {
  const { setAuth } = useAuthStore();
  const router = useRouter();

  return useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setAuth(data.user, data.token);
      router.push('/dashboard');
    },
  });
}

export function useLogout() {
  const { clearAuth } = useAuthStore();
  const router = useRouter();

  return () => {
    clearAuth();
    router.push('/login');
  };
}

export function useCurrentUser() {
  const { isAuthenticated } = useAuthStore();

  return useQuery({
    queryKey: ['currentUser'],
    queryFn: authApi.me,
    enabled: isAuthenticated,
  });
}
```

#### 5. Middleware（Next.js）

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value;
  const isAuthPage = request.nextUrl.pathname.startsWith('/login') || 
                     request.nextUrl.pathname.startsWith('/register');
  const isDashboardPage = request.nextUrl.pathname.startsWith('/dashboard') ||
                          request.nextUrl.pathname.startsWith('/projects') ||
                          request.nextUrl.pathname.startsWith('/members');

  // 未認証でダッシュボードにアクセス -> ログインページへ
  if (isDashboardPage && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // 認証済みでログインページにアクセス -> ダッシュボードへ
  if (isAuthPage && token) {
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/projects/:path*', '/members/:path*', '/login', '/register'],
};
```

## 権限管理（RBAC）

### ロール定義

1. **admin**: 全機能アクセス可能
2. **manager**: プロジェクト管理、メンバー管理可能
3. **member**: 自分のプロジェクトのみ閲覧・編集可能

### 権限チェック（Backend）

```go
// internal/middleware/role.go
func RequireRole(roles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            userRole := c.Get("role").(string)
            
            for _, role := range roles {
                if userRole == role {
                    return next(c)
                }
            }
            
            return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
        }
    }
}

// 使用例
e.DELETE("/projects/:id", projectHandler.Delete, authMiddleware, RequireRole("admin", "manager"))
```

### 権限チェック（Frontend）

```typescript
// hooks/use-permission.ts
export function usePermission() {
  const { user } = useAuthStore();

  const can = (permission: string): boolean => {
    if (!user) return false;

    const permissions = {
      admin: ['create', 'read', 'update', 'delete', 'manage_users'],
      manager: ['create', 'read', 'update', 'delete'],
      member: ['read', 'update_own'],
    };

    return permissions[user.role as keyof typeof permissions]?.includes(permission) || false;
  };

  return { can, role: user?.role };
}

// 使用例
function ProjectActions() {
  const { can } = usePermission();

  return (
    <div>
      {can('delete') && <Button onClick={handleDelete}>削除</Button>}
    </div>
  );
}
```

## セキュリティ考慮事項

### 1. パスワード要件
- 最小8文字
- 英数字含む（推奨）
- bcrypt cost: 12

### 2. JWT秘密鍵
- 最低32文字の乱数
- 環境変数で管理
- 本番環境では定期的にローテーション

### 3. トークン保存
- LocalStorage使用（シンプル、XSS対策必要）
- 将来的にはHttpOnly Cookieも検討

### 4. HTTPS必須
- 本番環境では必ずHTTPS使用
- トークンの盗聴防止

### 5. レート制限
- ログイン試行回数制限（5回/5分）
- API呼び出し制限

### 6. トークン有効期限
- Access Token: 24時間
- 将来的にRefresh Token導入検討

## 今後の拡張

1. **Refresh Token**
   - 長期間のセッション維持
   - Access Tokenの短命化

2. **OAuth 2.0**
   - Google、GitHub等でのログイン

3. **Multi-Factor Authentication (MFA)**
   - TOTP（Google Authenticator等）

4. **パスワードリセット**
   - メールによるリセットリンク送信

5. **アカウントロック**
   - 不正ログイン試行時の自動ロック

---

**Status**: Complete  
**Date**: 2025-10-11
