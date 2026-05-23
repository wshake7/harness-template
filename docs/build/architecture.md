# 架构总览

这个仓库已经从纯 Harness 模板演进为前后端一体的协作仓库。`docs/` 仍是长期知识源，`.service-matrix/dependencies.yaml` 是模块、服务和路径归属的单一真相源。

## 顶层结构

- `backend/go/`：Go workspace，包含管理后台 API、共享 Go 库、ORM/分页工具和 sa-token 适配。
- `front/`：pnpm workspace，包含三个 React 应用和一组 `@vp/*` 前端共享包。
- `docs/`：面向人和评审的架构、流程、质量、运行、安全和交付文档。
- `context/`：面向 Agent 执行的稳定上下文入口。
- `.service-matrix/`：项目、模块、服务和路径占位符的单一真相源。
- `.harness/`：可复用 Skill、Agent、Command 的源目录。
- `scripts/`：仓库级自动化脚本，供人和 Agent 直接调用。

## 后端拓扑

- `backend/go/admin` 是当前核心后端服务，模块名为 `admin`，使用 Go Fiber v3 提供 HTTP API，默认入口为 `cmd/main.go`。
- `backend/go/admin/etc/config.yaml` 是本地默认配置：服务默认监听 `0.0.0.0:3001`，开启 Swagger，使用 Postgres、Redis 和 Temporal。
- 后端路由统一注册到 `/api` 下，当前包括账号登录/登出、加密公钥、用户、角色、资源菜单/API、字典、语言、日志、任务调度、任务执行和知识库（Collection、Document）。
- `backend/go/admin/internal/appsvc` 承担 router/logic 面向的应用层服务接口和薄适配实现，和 `internal/services` 的基础设施生命周期层分开。
- `backend/go/admin/internal/ai` 存放管理后台内部 AI 组件与 Eino workflow：`agent/knowledge_pipeline` 负责编排，`embedder`、`loader`、`indexer` 提供可复用底层组件。
- `backend/go/admin/internal/workflows` 存放具体 Temporal Workflow 实现；`internal/services/temporaljob` 只保留调度分发、执行记录和 Worker 注册等任务基础能力。
- `backend/go/go-common` 提供通用 DTO、结果封装、日志、配置读取、集合、加密、ID、字符串、时间等 Go 共享能力。
- `backend/go/orm-crud` 提供 ORM CRUD、分页 proto 和 GORM 查询辅助能力；`backend/go/orm-crud/api` 含 Buf/protobuf 相关生成入口。
- `backend/go/go.work` 统一串起 `go-common`、`orm-crud/*`、`admin` 和 `sa-token/rueidis`，后台开发应优先在 workspace 根目录下验证。

## 前端拓扑

- `front/apps/admin-react` 是管理后台 SPA，使用 React、Vite/vite-plus、TanStack Router、Ant Design、zustand、i18next、alova 和 MSW。
- `front/apps/app-react` 是客户端 React 应用模板，复用同一套 `@vp/*` 包和请求/环境能力。
- `front/apps/app-react-ssr` 是 SSR 应用模板，额外接入 TanStack Start/Nitro。
- `front/packages/build-config` 沉淀共享构建配置。
- `front/packages/core` 定义跨应用领域类型、HTTP 常量、字典、分页、账号、通知和加密类型。
- `front/packages/react-core` 提供 React 侧环境变量、i18n、mock 和状态工厂。
- `front/packages/request` 封装 alova 请求客户端、Token 注入、登录响应处理、加密响应解密和统一错误处理。
- `front/packages/utils` 提供日期、CSS class、WebCrypto 等通用工具。

## 依赖与数据流

- 前端应用通过 workspace 依赖引用 `@vp/core`、`@vp/react-core`、`@vp/request`、`@vp/utils` 和 `@vp/build-config`。
- 前端开发服务器通过 Vite proxy 将 `/api` 转发到 `VITE_API_URL`；mock 模式由 `VITE_MOCK=true` 控制。
- 后端 API 统一返回业务响应体，Swagger 注释说明成功和失败 code；HTTP 状态不等同于业务状态。
- 登录接口使用加密中间件，前端在请求层处理公钥、AES key、Token Cookie 和响应解密。
- 管理后台数据主要经 Postgres 持久化，Redis 用于缓存/会话相关能力，Temporal 用于任务调度与执行。
- AI 知识索引流程通过 Ark (Volcengine) embedding 生成向量，并使用 Eino 官方 `components/indexer/milvus2` 组件接入现有 `github.com/milvus-io/milvus/client/v2` 客户端写入 Milvus。

## 边界约定

- 服务和包路径先更新 `.service-matrix/dependencies.yaml`，再同步相关文档。
- 前端共享逻辑优先放入 `front/packages/*`，应用目录只保留页面、路由、业务 API 和应用状态。
- 后端通用工具优先沉淀到 `backend/go/go-common`；`backend/go/admin/internal` 保持为管理后台服务内部实现。
- 生成代码、Swagger、ORM 查询和 protobuf 相关产物更新时，要同步记录生成命令和验证方式。
- 任何影响启动、依赖、配置、路由、安全或可观测性的变更，都要在同一轮更新对应文档和 history。
