## [2026-05-26 16:55] | Task: 前端公钥缓存兜底

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 登录前请求 `/api/encrypt/public/key` 后，前端会进入必须手动清理 `localhost:3000` / `3002` 站点数据才能继续请求的状态；先从前端修复，不要求兼容后端当前返回格式，只要非法格式不再污染缓存并直接提示请求错误即可。

### 🛠 Changes Overview

**Scope:** `front/packages/react-core`

**Key Actions:**

- **[非法公钥拦截]**: 在共享 `device-store` 中校验公钥缓存格式，检测到 PEM 或非 base64 值时不再写入持久化状态。
- **[坏缓存自清理]**: 读取旧缓存并尝试导入公钥时，如果格式非法或导入失败，立即清空本地缓存，避免用户再手动删浏览器站点数据。
- **[回归测试]**: 新增共享 store 测试，覆盖“非法格式不持久化”和“历史坏缓存自动清理”两个场景。

### 🧠 Design Intent (Why)

后端公钥格式后续会统一，因此前端不额外兼容 PEM；本次修复聚焦于状态自愈，避免一次异常返回把本地缓存永久污染，导致后续所有登录尝试都卡死在前端本地。

### 📁 Files Modified

- `front/packages/react-core/src/stores/factories.ts`
- `front/packages/react-core/tests/stores/factories.test.ts`
