# 后端协作说明

后端代码位于 `backend/go/`，使用 Go workspace 管理多个模块。当前核心服务是 `backend/go/admin`。

## 工作区结构

- `backend/go/admin`：管理后台 API 服务，入口为 `cmd/main.go`。
- `backend/go/go-common`：共享 Go 工具库，包含日志、配置、结果封装、分页 DTO、集合、加密、ID 和字符串工具。
- `backend/go/orm-crud`：ORM CRUD、分页 proto 和 GORM 查询辅助能力。
- `backend/go/sa-token/rueidis`：sa-token Redis 适配相关模块。
- `backend/go/go.work`：串联以上模块的 workspace 入口。

## 技术栈

- Go 版本：`1.26.3`，以各 `go.mod` 和 `go.work` 为准。
- HTTP 框架：Go Fiber v3。
- ORM 与数据库：GORM，默认配置使用 Postgres；代码中也保留 MySQL/SQLite 相关依赖。
- 认证与权限：sa-token-go、Casbin。
- 任务与异步：Temporal、Asynq。
- 缓存：Redis/rueidis。
- 日志与可观测性：zap、Prometheus client、Fiber monitor。
- API 文档：Swaggo，生成产物在 `backend/go/admin/docs/`。
- AI 组件：CloudWeGo Eino 工作流（ReAct Agent、Plan-Execute-Replan Agent、RAG、ChatTemplate、Milvus Indexer/Retriever），Ark (Volcengine) embedding，OpenAI 兼容 ChatModel，MCP SSE 工具集成，Prometheus 告警查询。

## 后端代码边界

- `backend/go/admin/internal/appsvc` 放 router/logic 面向的应用层服务接口和薄实现，例如认证、缓存、权限、数据权限和 Temporal 调度门面。
- `backend/go/admin/internal/services` 放 Redis、Temporal、Casbin、ORM、HTTP client、对象存储等基础设施生命周期服务，以及它们的底层适配能力。
- `backend/go/admin/internal/workflows` 放具体 Temporal Workflow 业务实现；任务分发、执行记录和 Worker 注册仍由 `internal/services/temporaljob` 承担。
- `backend/go/admin/internal/domains` 放跨路由、服务和中间件共享的领域常量与轻量 DTO，例如加密公钥缓存 key 和 key pair 结构。

## AI Agent 模块

AI Agent 代码位于 `backend/go/admin/internal/ai/`，按职责分层：

```
internal/ai/
├── agent/
│   ├── chat_pipeline/       ReAct Agent 图编排（RAG + ChatTemplate + 工具调用）
│   ├── knowledge_pipeline/  知识索引编排（Loader → Transformer → Indexer）
│   └── plan_execute_replan/ Plan-Execute-Replan Agent（规划→执行→重规）
├── models/open_ai.go        OpenAI 兼容 ChatModel 工厂
├── tools/                   内省工具（日志 MCP、Prometheus 告警、DB CRUD、时间、文档搜索）
├── mem/                     多后端对话记忆（Memory/Redis/DB）
├── embedder/                文本向量化
├── indexer/                 Milvus 向量索引
├── loader/                  文档加载
└── retriever/               Milvus RAG 检索
```

### Agent 模式

**chat_pipeline (ReAct Agent)**：适合单轮对话。用户输入 → Milvus RAG 检索 → ChatTemplate 拼装 prompt → ReAct Agent 推理并调用工具。入口为 `BuildChatAgent`，返回 `compose.Runnable`。

**plan_execute_replan (Plan-Execute-Replan Agent)**：适合复杂多步任务。Planner 先制定执行计划 → Executor 逐步执行并调用工具 → Replanner 根据执行结果决定是否重新规划。入口为 `BuildPlanAgent`。模型和工具通过依赖注入传入，因为模型创建需要 API Key/BaseURL 等外部配置，工具创建也需要 MCP URL、Embedding 配置等参数。

### 模型配置

`models.OpenAIChatModelConfig` 包含：
- `Model` / `APIKey` / `BaseURL` — 基础连接参数
- `ExtraFields map[string]any` — 透传到 OpenAI API 请求体的额外字段，例如 `{"thinking": {"type": "disabled"}}` 可禁用 DeepSeek 思考模式，解决 thinking 模型不支持 `tool_choice` 的问题

## 本地配置

`backend/go/admin/etc/config.yaml` 是本地默认配置：

| 配置 | 当前默认 |
| --- | --- |
| `Host` | `0.0.0.0` |
| `Port` | `3001` |
| `IsSwagger` | `true` |
| `Orm.DriverName` | `postgres` |
| `Orm.DataSource` | `host=localhost port=5432 user=postgres password=postgres dbname=wshake sslmode=disable` |
| `Redis.Addr` | `127.0.0.1:6379` |
| `Temporal.HostPort` | `127.0.0.1:7233` |
| `Temporal.TaskQueue` | `admin` |
| `Temporal.WorkerEnabled` | `true` |
| `AI.Embedding.Provider` | `ark` |
| `AI.Embedding.Model` | `doubao-embedding-vision-251215` |
| `AI.Embedding.Dimensions` | `2048` |
| `AI.Knowledge.Collection` | `biz` |
| `AI.Models` | ChatModel 数组，每项含 Name/Model/APIKey/BaseURL |
| `AI.Memory.Type` | `memory`（可选 `redis`、`db`） |
| `AI.McpURL` | MCP SSE 日志工具 Server URL |
| `Storage.Enabled` | `false` |
| `Storage.Engine` | `minio` |
| `Storage.MaxUploadBytes` | `10485760`（10 MiB） |
| `Storage.ObjectKeyPrefix` | `uploads` |
| `Storage.PresignedExpiresSeconds` | `3600` |
| `Storage.MinIO.Endpoint` | `127.0.0.1:9000` |
| `Storage.MinIO.Bucket` | `admin-files` |
| `Storage.MinIO.AutoCreateBucket` | `true` |

配置里包含本地数据库密码，仅用于开发默认值；生产或共享环境必须改用安全的配置注入方式。

## 文件上传与对象存储

- 通用文件上传接口位于 `/api/storage/file/*`，首期由服务端代理上传到当前配置的对象存储引擎。
- 当前默认引擎是 MinIO，通过 `internal/services/objectstore` 暴露统一 `Engine` 接口；后续接 S3 兼容存储或本地文件系统时不需要改业务 Handler。
- 文件元数据表为 `file_asset`，记录 engine、bucket、object key、原始文件名、content type、大小、sha256、业务标签和 JSON metadata。
- 上传接口 `POST /api/storage/file/upload` 使用 `multipart/form-data`，走登录态、Casbin 和语言中间件，但不走 JSON 加密中间件。
- 详情、短链和删除接口分别为 `POST /api/storage/file/detail`、`POST /api/storage/file/presigned`、`POST /api/storage/file/del`。
- 短链默认 1 小时，最大 7 天；下载或预览必须通过登录接口按需签发 presigned URL，不默认暴露永久公开地址。

## 常用命令

在 `backend/go/admin` 目录运行：

```bash
make run
make swagger
make script-imports
make script-orm
PATH="$(go env GOPATH)/bin:$PATH" go generate ./internal/appsvc
make test-cover
make test-html
```

在 `backend/go` workspace 根目录运行：

```bash
go test ./...
```

如果只验证管理后台服务：

```bash
cd backend/go/admin
go test ./...
```

## API 与请求约定

- 服务启动后，`cmd/main.go` 会读取 `-f` 指定的配置文件，默认是 `./etc/config.yaml`。
- 路由在 `backend/go/admin/internal/router/router.go` 中统一注册到 `/api`。
- 账号接口位于 `/api/account/*`，加密公钥接口位于 `/api/encrypt/public/key`。
- 需要登录的业务路由由 `auth_router` 注册，覆盖用户、角色、资源、字典、语言、日志、任务调度、任务执行、知识库（Collection、Document）和存储文件（Storage File）。
- `/api/storage/file/upload` 为 authenticated non-encrypted 路由，前端必须携带登录态，并通过请求层 `meta.skipEncrypt = true` 跳过 AES body 加密。
- Swagger 注释说明业务响应 code：成功为 `1`，通用失败为 `2`，请求超时/重放/错误和认证授权失败使用独立 code。
- 当前代码说明“所有接口均返回 HTTP 200，通过响应体 code 区分业务状态”；前端仍会对 HTTP 非 2xx 走统一错误路径。

## 外部依赖

- Postgres：默认数据库 `wshake`，默认本地用户 `postgres`。
- Redis：默认 `127.0.0.1:6379`。
- Temporal：默认 `127.0.0.1:7233`，namespace `default`，task queue `admin`。
- MinIO：默认示例地址 `127.0.0.1:9000`，bucket `admin-files`，由 `Storage.MinIO.*` 配置控制。
- Swagger：`IsSwagger=true` 时启用，文档生成入口为 `make swagger`。
- Ark (Volcengine)：需要在 `AI.Embedding.APIKey` 注入火山方舟 embedding API key。
- Milvus：AI 知识索引默认使用 `Milvus.DBName` 指定的数据库和 `AI.Knowledge.Collection` 指定的 collection。
- OpenAI 兼容 ChatModel：需要在 `AI.Models` 中配置 Name/Model/APIKey/BaseURL，当前使用 DeepSeek。
- MCP Server：AI Agent 通过 SSE 协议连接外部 MCP 日志查询 Server，URL 由 `AI.McpURL` 配置。

## 变更要求

- 修改配置字段、默认端口、外部依赖或启动方式时，同步更新本文件、`docs/operate/reliability.md` 和 README。
- 新增公共 API 或修改响应协议时，同步更新 Swagger 注释、前端请求包或业务 API 调用说明。
- 修改 `backend/go/admin/internal/services/orm/models/` 下的 ORM model 文件时，必须在 `backend/go/admin` 运行 `make script-orm`，并提交同步生成的 `internal/services/orm/query/` 产物。
- 生成代码、ORM 查询、Swagger 和 protobuf 产物要在 history 中记录生成命令。
- 涉及认证、加密、权限、Token、Cookie 或敏感配置的改动必须同步更新 `docs/operate/security.md`。
