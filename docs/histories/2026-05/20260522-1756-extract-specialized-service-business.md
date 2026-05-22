## [2026-05-22 17:56] | Task: 抽离 services 中的特化业务

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 在 services 里不要放这种特化场景的业务，给抽出来找个地方放。

### 🛠 Changes Overview

**Scope:** `backend/go/admin`

**Key Actions:**

- **[Boundary cleanup]**: 将加密公钥缓存 key 和 key pair DTO 从 `services/redisc` 移到 `internal/domains`，让 Redis 包只保留客户端能力。
- **[Workflow extraction]**: 将示例 `PrintCountWorkflow` 从 `services/temporaljob` 移到 `internal/workflows/example`，并通过 `internal/workflows/registry.go` 汇总具体业务 Workflow。
- **[Dependency direction]**: `services/temporaljob` 只注册分发 Workflow 和执行记录 Activity，具体 Workflow 注册器由 `cmd/main.go` 注入，避免 `services` 反向依赖业务场景。
- **[Docs]**: 更新后端和架构文档，记录 `services`、`domains` 与 `workflows` 的边界。

### 🧠 Design Intent (Why)

`services` 目录应承载基础服务适配和可复用服务层能力，不应混入加密缓存结构、示例 Workflow 这类特定业务场景。抽离后，基础设施包和业务实现的职责更清楚，也便于后续新增真实 Workflow 时按业务包组织。

### 📁 Files Modified

- `backend/go/admin/internal/domains/encrypt.go`
- `backend/go/admin/internal/workflows/example/print_count.go`
- `backend/go/admin/internal/service/redis_cache.go`
- `backend/go/admin/internal/fiberc/middleware/auth.go`
- `backend/go/admin/internal/services/temporaljob/job.go`
- `docs/ARCHITECTURE.md`
- `docs/BACKEND.md`
