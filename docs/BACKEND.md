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

## 后端代码边界

- `backend/go/admin/internal/appsvc` 放 router/logic 面向的应用层服务接口和薄实现，例如认证、缓存、权限、数据权限和 Temporal 调度门面。
- `backend/go/admin/internal/services` 放 Redis、Temporal、Casbin、ORM、HTTP client 等基础设施生命周期服务，以及它们的底层适配能力。
- `backend/go/admin/internal/workflows` 放具体 Temporal Workflow 业务实现；任务分发、执行记录和 Worker 注册仍由 `internal/services/temporaljob` 承担。
- `backend/go/admin/internal/domains` 放跨路由、服务和中间件共享的领域常量与轻量 DTO，例如加密公钥缓存 key 和 key pair 结构。

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

配置里包含本地数据库密码，仅用于开发默认值；生产或共享环境必须改用安全的配置注入方式。

## 常用命令

在 `backend/go/admin` 目录运行：

```bash
make run
make swagger
make script-imports
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
- 需要登录的业务路由由 `auth_router` 注册，覆盖用户、角色、资源、字典、语言、日志、任务调度和任务执行。
- Swagger 注释说明业务响应 code：成功为 `1`，通用失败为 `2`，请求超时/重放/错误和认证授权失败使用独立 code。
- 当前代码说明“所有接口均返回 HTTP 200，通过响应体 code 区分业务状态”；前端仍会对 HTTP 非 2xx 走统一错误路径。

## 外部依赖

- Postgres：默认数据库 `wshake`，默认本地用户 `postgres`。
- Redis：默认 `127.0.0.1:6379`。
- Temporal：默认 `127.0.0.1:7233`，namespace `default`，task queue `admin`。
- Swagger：`IsSwagger=true` 时启用，文档生成入口为 `make swagger`。

## 变更要求

- 修改配置字段、默认端口、外部依赖或启动方式时，同步更新本文件、`docs/RELIABILITY.md` 和 README。
- 新增公共 API 或修改响应协议时，同步更新 Swagger 注释、前端请求包或业务 API 调用说明。
- 生成代码、ORM 查询、Swagger 和 protobuf 产物要在 history 中记录生成命令。
- 涉及认证、加密、权限、Token、Cookie 或敏感配置的改动必须同步更新 `docs/SECURITY.md`。
