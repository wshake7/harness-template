# 登录互踢：sa-token-go + Redis 适配器兼容性

管理后台 `AuthMiddleware` 基于 sa-token-go 实现登录态校验。配置 `IsConcurrent: false` 启用互踢——新登录会挤掉旧登录的 token。

## 现象

A 浏览器登录后，B 浏览器用同一账号登录，A 的请求**没有被 `AuthMiddleware` 拦截**，而是走到了 `EncryptMiddleware` 并报 `RSA解密失败`。

## 根因

sa-token-go 的 `kickout()` 逻辑：

```go
// sa-token-go/core/manager/manager.go
func (m *Manager) kickout(loginID string, device string) error {
    accountKey := m.getAccountKey(loginID, device)
    tokenValue, err := m.storage.Get(accountKey)
    ...
    tokenStr, ok := assertString(tokenValue)  // ← 只认 string 类型
    if !ok {
        return nil  // 静默跳过，不踢出旧 token！
    }
    return m.removeTokenChain(tokenStr, false, listener.EventKickout)
}
```

`assertString()` 只做 `v.(string)` 类型断言，不处理 `[]byte`。

而项目使用的 rueidis 适配器 `Get()` 返回的是 `[]byte`（`AsBytes()`）：

```go
// sa-token/rueidis/rueidis.go (修复前)
func (s *Storage) Get(key string) (any, error) {
    result, err := s.client.Do(ctx, ...).AsBytes()
    return result, nil  // 返回 []byte
}
```

→ `kickout` 走到 `assertString([]byte(...))` → `ok = false` → `return nil` → 旧 token 从未被标记为 `KICK_OUT` → 旧 token 依然有效 → `AuthMiddleware` 放行。

## 修复

### 1. 修复 Redis 适配器（根因）

`sa-token/rueidis/rueidis.go:212`：`Get()` 返回 `string` 而非 `[]byte`：

```go
return string(result), nil
```

sa-token-go 的 `getTokenInfo` 同时兼容 `[]byte` 和 `string` 分支，其他调用路径不受影响。但这使得 `assertString` 能够正确匹配，`kickout` 不再静默失败。

### 2. AuthMiddleware 二次校验（防御层）

即使 `kickout` 正确将旧 token 标记为 `KICK_OUT`，仍有可能因 Redis TTL 竞态、`SetKeepTTL` 失败等边界情况让旧 token 滑过 `getTokenInfo` 检查。

在 `AuthMiddleware` 中追加一次校验：拿请求 token 与 `stputil.GetTokenValue(loginID)`（当前有效 token）比对，不一致即返回 `FailTokenReplaced`。

### 3. 修复误用的 fallback 逻辑

`AuthMiddleware` 原来的 "header token failed, fallback to cookie token" retry 用的是**同一个** `token` 变量，第二次调用没有任何不同。修复后仅在 header token 和 cookie token 值不同时才真正切换。

## 影响范围

| 文件 | 变更 |
|------|------|
| `backend/go/sa-token/rueidis/rueidis.go` | `Get()` 返回值从 `[]byte` 改为 `string` |
| `backend/go/admin/internal/fiberc/middleware/auth.go` | 新增 token 一致性校验；修正 cookie fallback |
