# 稳定性与可运维性

这份文档定义当前前后端项目的运行质量底线。CI/CD 流程结构和 release 自动化统一写在 `docs/operate/cicd.md`。

## 启动与关键路径

- 后端核心服务是 `backend/go/admin`，默认从 `./etc/config.yaml` 读取配置，监听 `0.0.0.0:3001`。
- 后端 API 前缀为 `/api`，Swagger 注释和生成产物位于 `backend/go/admin/docs/`。
- 管理后台前端是 `front/apps/admin-react`，本地通过 `VITE_API_URL` 将 `/api` proxy 到后端。
- `VITE_MOCK=true` 时前端启用 MSW，可用于缺少后端依赖时的页面验证。
- 本地联调默认依赖 Postgres、Redis、Temporal 和可选的 MinIO；缺少 Postgres/Redis/Temporal 时后端主链路会受影响，缺少 MinIO 时可保持 `Storage.Enabled=false` 或前端改走 mock。
- AI 知识索引链路额外依赖 Ark (Volcengine) embedding 凭据和 Milvus collection 配置；缺配置时知识入库流程会初始化失败。
- AI Chat Agent 额外依赖 OpenAI 兼容 ChatModel API 凭据（`AI.Models`）和可选的 MCP SSE 日志 Server（`AI.McpURL`）；缺 ChatModel 配置时 Agent 图编译失败，缺 MCP URL 时日志查询工具以 Warn 级别降级跳过。

## 外部依赖

| 依赖 | 默认地址 | 用途 |
| --- | --- | --- |
| Postgres | `localhost:5432` | 管理后台业务数据和权限数据。 |
| Redis | `127.0.0.1:6379` | 缓存、会话或加密相关临时数据。 |
| Temporal | `127.0.0.1:7233` | 任务调度、任务执行和 workflow 管理。 |
| MinIO | `127.0.0.1:9000` | 文件上传对象存储，默认 bucket `admin-files`。 |
| Milvus | `127.0.0.1:19530` | 向量数据库，用于向量检索和语义搜索。 |
| OpenAI 兼容 ChatModel API | 由 `AI.Models[].BaseURL` 配置 | LLM 推理，当前接入 DeepSeek。 |
| MCP SSE Server | 由 `AI.McpURL` 配置 | 日志查询 MCP 工具 Server（可选）。 |

本地默认值来自 `backend/go/admin/etc/config.yaml`。生产或共享环境必须提供环境隔离后的配置，不应直接复用本地默认连接串。

## 可观测性

- 后端依赖 zap 日志；服务初始化在 `backend/go/admin/cmd/main.go` 中调用 `log.Init(log.DevConfig())`。
- 后端依赖 Prometheus client 和 Fiber monitor，但当前仓库尚未把指标入口、采集路径和仪表盘写成稳定运行手册。
- 对象存储健康状态通过 Fiber service 生命周期暴露：`Storage.Enabled=false` 时状态为 `DISABLED`，启用后会在启动阶段执行 bucket 检查，并按配置决定是否自动创建 bucket。
- 前端请求层使用 NProgress、统一错误提示和控制台 API error 输出，适合本地调试，不等同于生产观测方案。
- API 文档由 Swaggo 维护，变更公共接口时运行 `make swagger` 并提交生成产物。

## 验证命令

仓库级：

```bash
make ci
```

前端：

```bash
pnpm ready
cd front/apps/admin-react && pnpm e2e:test
```

后端：

```bash
cd backend/go
go test ./...
```

如果项目级验证失败，应在 history 或 PR 描述中记录命令、退出状态、失败原因和后续动作。

## 常见排障路径

- 前端请求失败：检查 `VITE_API_URL`、Vite proxy、后端是否监听 `3001`、浏览器请求路径是否以 `/api` 开头。
- 登录或加密请求失败：检查 `/api/encrypt/public/key`、Cookie 中 Token、前端 `encryptRequest` 和后端加密中间件。
- 上传接口失败：检查 `Storage.Enabled`、`Storage.MinIO.Endpoint`、`Storage.MinIO.Bucket`、`Storage.MaxUploadBytes`，并确认上传请求走 `meta.skipEncrypt=true`，没有被 JSON 加密链路改写。
- 后端启动失败：检查 `backend/go/admin/etc/config.yaml` 中 Postgres、Redis、Temporal、MinIO、Milvus 是否可访问。
- AI 知识索引失败：检查 `AI.Embedding.APIKey`（Ark API key）、`AI.Embedding.Model`（如 `doubao-embedding-vision-251215`）、`AI.Embedding.Dimensions`（该模型默认 2048）、`Milvus.DBName` 和 `AI.Knowledge.Collection` 是否与目标环境一致。若报 "collection schema mismatch"，说明 Milvus collection 的向量维度与当前配置不匹配，需要删除旧 collection 后重新创建。
- AI Chat Agent 失败：检查 `AI.Models` 中是否配置了目标模型别名（如 `deepseek-v4-flash`）及对应的 APIKey/BaseURL；确认 `AI.McpURL` 指向的 MCP SSE Server 可访问（不可访问时 Agent 会降级跳过日志工具）。若 ReAct Agent 达到 `MaxStep=25` 仍未完成，会返回达到最大步数的错误。
- 任务调度失败：检查 Temporal 地址、namespace、task queue `admin` 和 worker 是否启用。
- Swagger 不一致：在 `backend/go/admin` 运行 `make swagger`，确认 `docs/` 生成产物同步。

## 当前短板

- 还没有 Docker Compose 或等价脚本一键拉起 Postgres、Redis、Temporal。
- 还没有稳定的健康检查端点和指标访问约定。
- `make ci` 目前仍是仓库基础门禁，没有强制执行 `pnpm ready` 或 `go test ./...`。
- 生产级日志、metrics、traces、告警和 dashboard 尚未定义。
