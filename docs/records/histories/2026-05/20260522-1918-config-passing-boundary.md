# 配置传递边界清理：services 与内部组件不再依赖全局 config.Conf

## 背景

后端代码中 `config.Conf` 全局变量被多处直接引用，导致配置边界模糊：
- `services/` 层（基础设施生命周期）直接读取全局配置
- `middleware/`（语言、认证）在运行时直接访问全局配置
- `ai/` 组件（embedder、indexer）在初始化时直接读取全局配置
- `router/logic/`（业务 handler）直接读取全局配置

这使得 `services/` 与外部配置强耦合，不利于后续抽取通用代码。

## 目标

将全局配置引用改为**显式参数传递**，明确边界：
- `services/` 层：由 `main.go` 在启动时传递所需配置片段
- `middleware/`：由 router 注册时传入配置值
- `ai/` 组件：由调用方传入配置对象
- `router/logic/`：由 handler 构造函数传入配置值

## 改动概览

### 1. middleware 层（显式参数传递）

| 文件 | 变更 |
|------|------|
| `internal/fiberc/middleware/language.go` | `LanguageMiddleware()` → `LanguageMiddleware(defaultLanguage string)` |
| `internal/fiberc/middleware/auth.go` | `AuthMiddleware()` → `AuthMiddleware(tokenName string)` |

### 2. router 层（配置向下传递）

| 文件 | 变更 |
|------|------|
| `internal/router/router.go` | `Router` 增加 `Conf *config.Config` 字段；`RegisterRouters` 传入配置 |
| `internal/router/account.go` | `registerAccountRouters` 增加 `conf *config.Config` 参数 |
| `internal/router/auth_router/auth_router.go` | `RegisterRouters` 增加 `conf *config.Config` 参数 |
| `internal/router/auth_router/job_schedule.go` | `registerJobScheduleRouters` 增加 `conf *config.Config` 参数 |

### 3. logic 层（handler 持有配置）

| 文件 | 变更 |
|------|------|
| `internal/router/logic/job_schedule.go` | `JobScheduleHandler` 增加 `DefaultTaskQueue string`；`NewJobScheduleHandler` 增加参数 |

### 4. AI 组件层（配置注入）

| 文件 | 变更 |
|------|------|
| `internal/ai/embedder/embedder.go` | `New(ctx)` → `New(ctx, conf config.AIEmbeddingConfig)` |
| `internal/ai/indexer/indexer.go` | `New(ctx)` → `New(ctx, conf *config.Config)`；内部函数全部接收配置参数 |
| `internal/ai/agent/knowledge_pipeline/indexer.go` | `newIndexer(ctx)` → `newIndexer(ctx, conf *config.Config)` |
| `internal/ai/agent/knowledge_pipeline/embedding.go` | `newEmbedding(ctx)` → `newEmbedding(ctx, conf config.AIEmbeddingConfig)` |
| `internal/ai/agent/knowledge_pipeline/orchestration.go` | `BuildKnowledgeIndexing(ctx)` → `BuildKnowledgeIndexing(ctx, conf *config.Config)` |

### 5. main.go（启动时传递）

| 文件 | 变更 |
|------|------|
| `cmd/main.go` | `router.Router{Conf: conf}` 初始化 |

### 6. 测试修复

| 文件 | 变更 |
|------|------|
| `internal/router/logic/job_schedule_test.go` | 所有 `NewJobScheduleHandler` 调用增加 `"admin"` 参数 |

## 验证

- `go build ./admin/... ./go-common/...` ✅
- `go test ./admin/...` ✅（全部通过）

## 影响

- **编译通过**：无行为变更，纯重构
- **测试通过**：`job_schedule_test.go` 已同步更新
- **边界清晰**：`services/` 层仍由 `main.go` 注入完整配置（因其负责基础设施生命周期），但内部组件不再直接访问全局变量
- **后续可抽取**：`services/` 层的基础设施初始化逻辑现在只依赖传入参数，可更容易抽取为通用库

## 关键文件

- `backend/go/admin/internal/fiberc/middleware/language.go`
- `backend/go/admin/internal/fiberc/middleware/auth.go`
- `backend/go/admin/internal/router/router.go`
- `backend/go/admin/internal/router/account.go`
- `backend/go/admin/internal/router/auth_router/auth_router.go`
- `backend/go/admin/internal/router/auth_router/job_schedule.go`
- `backend/go/admin/internal/router/logic/job_schedule.go`
- `backend/go/admin/internal/ai/embedder/embedder.go`
- `backend/go/admin/internal/ai/indexer/indexer.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/orchestration.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/indexer.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/embedding.go`
- `backend/go/admin/cmd/main.go`
- `backend/go/admin/internal/router/logic/job_schedule_test.go`
