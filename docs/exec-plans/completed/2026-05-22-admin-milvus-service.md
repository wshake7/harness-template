# Admin Milvus 基础 Service 执行计划

## 目标

为 `backend/go/admin` 添加 Milvus 基础外部依赖服务，让 admin 后端可以通过统一配置初始化 Milvus Go SDK 客户端，并纳入 Fiber service 生命周期、健康状态检查和关闭释放流程。

## 范围

- 包含：
  - 新增 Milvus 配置结构和本地默认配置。
  - 新增 Milvus 客户端封装和 Fiber `Service` 生命周期适配。
  - 在 admin 服务初始化时注册 Milvus service。
  - 更新 Go 依赖、后端文档、稳定性文档、安全文档和 history。
- 不包含：
  - 不新增 HTTP API 或 Swagger 业务接口。
  - 不设计 collection schema、向量维度、索引策略。
  - 不实现 insert、search、query、collection 管理等业务能力。

## 背景

- 相关文档：
  - `docs/BACKEND.md`
  - `docs/RELIABILITY.md`
  - `docs/SECURITY.md`
  - `docs/PLANS_GUIDE.md`
- 相关代码路径：
  - `backend/go/admin/internal/config/config.go`
  - `backend/go/admin/etc/config.yaml`
  - `backend/go/admin/internal/services/init.go`
  - `backend/go/admin/internal/services/`
  - `backend/go/admin/go.mod`
- 已知约束：
  - 现有外部依赖通过 `conf.Fiber.Services` 注册，并实现 `Start`、`State`、`Terminate` 和 `String`。
  - Milvus Go SDK 文档使用 `github.com/milvus-io/milvus/client/v2/milvusclient`，通过 `milvusclient.New(ctx, &milvusclient.ClientConfig{Address, APIKey})` 建连，通过 `client.Close(ctx)` 关闭。
  - Milvus 默认启用，本地默认地址为 `127.0.0.1:19530`。
  - 示例配置中的 `APIKey` 留空，不提交默认 token 或生产凭据。

## 实现方案

- 配置：
  - 在 `config.Config` 中新增 `Milvus config.MilvusConfig`。
  - 新增 `MilvusConfig`，字段为 `Enabled bool`、`Address string`、`APIKey string`，可选支持 `DBName string`。
  - 在 `etc/config.yaml` 增加：
    - `Enabled: true`
    - `Address: 127.0.0.1:19530`
    - `APIKey:` 留空
- 客户端封装：
  - 新增 `internal/services/milvusc` 子包。
  - 保持与 `redisc`、`temporalc` 类似的包级 `Client` 风格。
  - `New(ctx, conf)` 创建 Milvus client；`Close(ctx)` 关闭并清空包级 client。
- 生命周期服务：
  - 新增 `internal/services/milvus.go`。
  - `Start(ctx)`：`Enabled=false` 时 no-op；启用时创建客户端。
  - `State(ctx)`：关闭返回 `DISABLED`；启用但未初始化或探测失败返回 `UNHEALTHY`；正常返回 `HEALTHY`。
  - `Terminate(ctx)`：关闭 Milvus client；未启用或未初始化时 no-op。
  - 健康探测优先使用 SDK 稳定、轻量的只读探测能力；如果 SDK 没有专用 ping，则使用最小元数据调用。
- 服务注册：
  - 在 `services.New(conf)` 中追加 `NewMilvus(conf.Milvus)`，放在 ORM、Redis、Temporal 等基础依赖附近。

## 风险

- 风险：默认启用后，本地未启动 Milvus 时 admin 服务启动或健康检查会失败。
- 缓解方式：文档明确默认地址、关闭方式和排障路径；本地没有 Milvus 时可把 `Milvus.Enabled` 改为 `false`。
- 风险：Milvus `APIKey` 属于外部服务凭据，误提交真实 token 会造成泄露。
- 缓解方式：示例配置留空，并在安全文档中记录生产环境必须通过安全配置注入。
- 风险：SDK 版本和健康探测 API 可能变化。
- 缓解方式：实现前用 `ctx7` 确认当前 Milvus Go SDK 文档；测试覆盖 disabled、未初始化和关闭路径。

## 里程碑

1. 配置和依赖接入：新增 `MilvusConfig`、示例配置和 Go SDK 依赖。
2. 客户端与 service 实现：完成 `milvusc` 和 `Milvus` Fiber service。
3. 初始化与文档同步：注册 service，更新后端、稳定性、安全和 history 文档。
4. 验证与收尾：运行后端测试，记录任何外部依赖不可用导致的验证限制。

## 验证方式

- 命令：
  - `cd backend/go/admin && go test ./...`
  - `cd backend/go && go test ./...`
- 单元测试：
  - `Enabled=false` 时 `Start` 不建连且不报错。
  - `Enabled=false` 时 `State` 返回 `DISABLED`。
  - enabled 但 client 未初始化时 `State` 返回 `UNHEALTHY`。
  - 未启用或未初始化时 `Terminate` 不报错。
- 手工检查：
  - 有本地 Milvus `127.0.0.1:19530` 时，admin 服务可启动并正常关闭。
  - 无本地 Milvus 时，失败原因指向 Milvus 连接或健康检查。
- 文档检查：
  - `docs/BACKEND.md`、`docs/RELIABILITY.md`、`docs/SECURITY.md` 和 history 已同步外部依赖、配置和凭据约束。

## 进度记录

- [x] 确认任务范围为基础连接 service，不包含业务 API。
- [x] 确认 Milvus 默认启用。
- [x] 确认 `APIKey` 示例留空。
- [x] 完成配置和依赖接入。
- [x] 完成 Milvus client 和 Fiber service。
- [x] 完成文档、history 和验证。

## 决策记录

- 2026-05-22：只做 Milvus 基础连接 service，先不暴露管理 API 或向量业务能力，降低首轮接入范围。
- 2026-05-22：Milvus 默认启用，便于本地和集成环境尽早暴露外部依赖缺失问题。
- 2026-05-22：`APIKey` 在示例配置中留空，避免把默认凭据提交到仓库。
