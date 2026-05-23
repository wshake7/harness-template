## [2026-05-22 22:35] | Task: rename admin appsvc

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> `backend/go/admin/internal/service` 和 `backend/go/admin/internal/services` 命名重复，想优化这两个包的区分度。

### 🛠 Changes Overview

**Scope:** `backend/go/admin`, `docs/`

**Key Actions:**

- **Rename application-facing package**: 将 `backend/go/admin/internal/service` 重命名为 `backend/go/admin/internal/appsvc`，保留原有接口与默认实现行为不变。
- **Update call sites**: 同步更新 router、logic 和测试中的 import 与包引用，从 `service.*` 切到 `appsvc.*`。
- **Sync docs**: 更新后端与架构文档，明确 `internal/appsvc` 和 `internal/services` 的职责边界。

### 🧠 Design Intent (Why)

原先 `service` 与 `services` 只差单复数，但实际承担的是两层完全不同的职责：前者是给 router/logic 注入的应用层门面，后者是启动期基础设施生命周期服务。把前者改成 `appsvc`，可以在不大规模迁移 ORM/生成代码路径的前提下，先把最常用、最容易混淆的边界名称澄清。

### 📁 Files Modified

- `backend/go/admin/internal/appsvc/*.go`
- `backend/go/admin/internal/router/**/*.go`
- `backend/go/admin/internal/mock/*.go`
- `docs/BACKEND.md`
- `docs/ARCHITECTURE.md`
