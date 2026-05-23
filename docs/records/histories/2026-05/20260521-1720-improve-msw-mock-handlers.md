## [2026-05-21 17:20] | Task: 完善前端 MSW Mock Handlers 覆盖

### 🤖 Execution Context

- **Agent ID**: claude-opus-4-7
- **Base Model**: Opus 4.7 (1M context)
- **Runtime**: darwin / zsh

### 📥 User Query

> 现在前端的mocks的handlers里msw的mock没有覆盖完全，并且对应接口也有些问题，给我完善一下。admin进去后菜单没有内容。

### 🛠 Changes Overview

**Scope:** `front/apps/admin-react`, `front/apps/app-react`, `front/apps/app-react-ssr`, `front/packages/request`

**Key Actions:**

- **[Fix account handler]** `admin-react/src/mocks/handlers/account.ts` login 响应补充 `publicKey` 字段，与 `ResLogin` 类型定义对齐
- **[Add admin-react handlers]** 新增 12 个业务域 handler 文件，覆盖约 57 个接口：
  - `encrypt.ts` — 返回 RSA-SPKI 公钥（mock 明文）
  - `dict.ts` — 字典类型/条目 CRUD + match + batch/copy
  - `sysUser.ts` — 用户管理 CRUD
  - `role.ts` — 角色 CRUD + tree + permissions
  - `resourceMenu.ts` — 菜单资源 CRUD + tree
  - `resourceApi.ts` — API 资源 CRUD
  - `apiLog.ts` / `loginLog.ts` — 日志列表/详情
  - `jobSchedule.ts` — 任务调度 CRUD + options + switch/sync/trigger
  - `jobExecution.ts` — 任务执行 列表/详情 + cancel/retry
  - `language.ts` — 语言类型/条目 CRUD + batch/create
  - `events.ts` — SSE EventStream
- **[Add app-react / app-react-ssr handlers]** 各新增 3 个 handler（encrypt、dict match、events）
- **[Fix mock mode encryption]** `api/index.ts` 三应用均在 `VITE_MOCK=true` 时跳过请求加密与响应解密；`@vp/request` login 回调兼容无 `aesKey` 的明文响应，使 mock 模式下 login 流程正常解析 token
- **[Fix menu paths]** `resourceMenu.ts` mock 菜单路径与实际 TanStack Router 路由表对齐（`/system/*`、`/account/*`、`/job/*`、`/logger/*`、`/dashboard`）

### 🧠 Design Intent (Why)

1. **覆盖率不足** — admin-react 原有 3 个 handler，仅占 60 个接口的 5%；app-react / app-react-ssr 为 0%。在 `VITE_MOCK=true` 开发模式下，大量页面因 fallback 到真实后端而无法正常工作。
2. **请求体加密冲突** — 前端通过 `encryptRequest` 对请求体做 RSA+AES 加密，MSW handler 用 `request.json()` 解析时遇到加密二进制数据，导致 JSON 解析异常。Mock 模式下直接跳过加密，保持 handler 与前端都走明文 JSON。
3. **菜单路径不匹配** — 前端 `toDynamicMenuItems` 会校验 `routesByPath[node.path]`，mock 数据中的 `/system/user`、`/system/role` 等路径在路由表中不存在（实际为 `/account/user`、`/account/role`），导致所有菜单项被过滤为空。

### 📁 Files Modified

- `front/apps/admin-react/src/mocks/handlers/account.ts`
- `front/apps/admin-react/src/mocks/handlers/encrypt.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/dict.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/sysUser.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/role.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/resourceMenu.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/resourceApi.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/apiLog.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/loginLog.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/jobSchedule.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/jobExecution.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/language.ts` *(new)*
- `front/apps/admin-react/src/mocks/handlers/events.ts` *(new)*
- `front/apps/admin-react/src/api/index.ts`
- `front/apps/app-react/src/mocks/handlers/encrypt.ts` *(new)*
- `front/apps/app-react/src/mocks/handlers/dict.ts` *(new)*
- `front/apps/app-react/src/mocks/handlers/events.ts` *(new)*
- `front/apps/app-react/src/api/index.ts`
- `front/apps/app-react-ssr/src/mocks/handlers/encrypt.ts` *(new)*
- `front/apps/app-react-ssr/src/mocks/handlers/dict.ts` *(new)*
- `front/apps/app-react-ssr/src/mocks/handlers/events.ts` *(new)*
- `front/apps/app-react-ssr/src/api/index.ts`
- `front/packages/request/src/client/alova.ts`
