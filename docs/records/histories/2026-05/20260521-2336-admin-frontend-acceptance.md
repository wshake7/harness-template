## [2026-05-21 23:36] | Task: 完善 admin 前端测试验收

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 完善前端 admin 的测试和验收，并按已确认计划实现稳定 e2e、补齐核心页面冒烟覆盖、同步验收文档和 history。

### 🛠 Changes Overview

**Scope:** `front/apps/admin-react`

**Key Actions:**

- **稳定 E2E**: 固定 admin Playwright e2e 单 worker 执行，抽取统一 mock 登录 helper，减少多 worker 登录态和 mock 数据干扰。
- **补齐核心验收**: 新增 API 资源、语言、任务配置、执行记录、API 日志和登录日志的 e2e 冒烟与关键交互覆盖。
- **同步验收说明**: 更新 admin-react 验收清单，记录当前测试范围、稳定性策略和可直接运行的验证命令。

### 🧠 Design Intent (Why)

admin-react 已有单测、组件测试和部分 e2e，但业务页面覆盖不完整，且 e2e 在默认并发下存在不稳定现象。此次变更先把 e2e 变成可复现门禁，再用轻量冒烟测试覆盖核心管理页面，避免把完整 CRUD 深度测试一次性做重。

### 📁 Files Modified

- `front/apps/admin-react/playwright.config.ts`
- `front/apps/admin-react/tests/helpers/auth.ts`
- `front/apps/admin-react/tests/business/*.spec.ts`
- `front/apps/admin-react/docs/acceptance-checklist.md`
