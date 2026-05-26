# Go Admin 到 Java Admin 移植指南

这份文档用于把当前 `backend/go/admin` 管理后台服务移植到 `backend/java/admin`。当前 Java 目录还是空目录，因此本文先定义移植边界、分层映射、阶段顺序和验收口径，后续真正落 Java 工程时应以本文为执行基线。

## 目标

- 在 `backend/java/admin` 建立一份功能等价的 Java 管理后台服务。
- 保持前端 `front/apps/admin-react` 的 API 协议、登录态、加密请求、业务响应码和错误处理尽量不变。
- 复用现有 Postgres、Redis、Temporal、MinIO、Milvus、OpenAI 兼容模型和 MCP SSE Server 等外部依赖语义。
- 让 Go 与 Java 两个服务可以在迁移期并行存在，方便按模块灰度验证和回退。
- 把一次性翻译拆成可验收阶段，优先迁移基础框架、认证权限和核心 CRUD，再迁移异步任务、对象存储、知识库和 AI Agent。

## 非目标

- 不在第一阶段重做前端请求协议或页面交互。
- 不把 Go 共享库直接暴露成 Java 服务依赖；需要在 Java 侧重建等价能力。
- 不要求第一版 Java 服务同时覆盖所有 AI Agent 能力；AI 链路可以最后迁移。
- 不在移植时顺手大改数据库表结构。确需调整时必须提供兼容迁移脚本和回滚方案。
- 不把 Go 服务删掉。Java 服务通过独立端口或配置启动，验证通过后再讨论切流。

## 当前 Go Admin 基线

Go 服务入口为 `backend/go/admin/cmd/main.go`，启动顺序是：

1. 读取 `etc/config.local.yaml` 或 `-f` 指定配置。
2. 初始化 zap 日志。
3. 初始化基础设施服务，包含 ORM、Redis、Temporal、Casbin、对象存储、Milvus、HTTP client 等。
4. 创建 Fiber 应用，注册恢复、日志、trace、CORS、pprof、健康检查、metrics、Swagger。
5. 在配置的 `RestPrefix` 下注册 `/api` 路由。

主要模块如下：

| Go 路径 | 职责 |
| --- | --- |
| `cmd/main.go` | 服务启动入口。 |
| `etc/config.yaml` | 本地默认配置。 |
| `internal/config` | 配置结构和解析。 |
| `internal/fiberc` | Fiber 应用封装、中间件、响应处理。 |
| `internal/router` | `/api` 根路由、账号和加密路由。 |
| `internal/router/auth_router` | 登录后业务路由注册。 |
| `internal/router/logic` | HTTP Handler 与业务编排。 |
| `internal/appsvc` | 面向 Handler 的应用服务接口和薄适配。 |
| `internal/services` | ORM、Redis、Temporal、Casbin、对象存储、Milvus 等基础设施生命周期。 |
| `internal/services/orm/models` | GORM model，当前数据库表定义来源。 |
| `internal/services/orm/query` | GORM Gen 生成查询代码。 |
| `internal/workflows` | Temporal Workflow 实现。 |
| `internal/ai` | Eino Agent、RAG、embedding、Milvus indexer/retriever、工具和记忆。 |

## 已定 Java 技术基线

Java 版 admin 采用下面的技术栈。Spring Boot 固定在 3.x 最新稳定线，不跟随 Spring Boot 4.x 升级；当前 ctx7 可查到的 Spring Boot 3 参考版本为 `3.5.9`。Easy Query 固定采用 `dromara/easy-query` v3 最新稳定版。

| 能力 | Java 建议 | 迁移注意 |
| --- | --- | --- |
| JDK | Java 25 | 实现前确认 Gradle toolchain、CI 镜像和依赖库对 JDK 25 的兼容性。 |
| 构建工具 | Gradle | Java admin 使用独立 Gradle 工程；后续如纳入仓库 CI，需要补 Gradle 验证命令。 |
| HTTP 服务 | Spring Boot 3.x 最新稳定版 + Spring Boot Web | 保持 `/api` 路径和 JSON 协议兼容，不升级到 Spring Boot 4.x。 |
| 配置 | Spring Boot configuration properties | 字段命名可 Java 化，但要提供 Go 配置到 Java 配置的映射。 |
| 样板代码生成 | Lombok | 用于 entity、DTO、配置类和构造器等样板代码；密码、token、private key 等敏感字段不要被 `toString` 输出。 |
| 类转换 | linpeilie/mapstruct-plus | 用于 entity、DTO、VO、request、response 的显式转换，避免在 Controller 和业务层手写重复拷贝逻辑。 |
| ORM | dromara/easy-query v3 最新稳定版 | 以 Easy Query entity/proxy/query API 复刻现有 GORM model、分页和条件查询。 |
| 数据库迁移 | Flyway | 不继续依赖运行时 AutoMigrate 作为生产迁移手段。 |
| 认证会话 | Sa-Token Java | 前端 Cookie/Token 名称和登录响应必须兼容。 |
| 权限 | jCasbin | 要兼容现有 `sys_role`、`sys_resource_*`、Casbin rule 数据。 |
| Redis | Spring Data Redis | 用于会话、缓存、加密 key、nonce 等。 |
| Temporal | Temporal Java SDK | namespace、task queue、workflow 名称要与 Go 侧兼容或显式分版本。 |
| 对象存储 | MinIO Java SDK | 保持上传、预签名、完成上传语义。 |
| API 文档 | Knife4j | 路由、请求体、响应体和业务 code 要对齐现有 Swagger/OpenAPI。 |
| 可观测性 | Micrometer + Actuator | 对齐 `/metrics`、健康检查和 Prometheus 采集目标。 |

如果后续调整以上选型，需要在本文和 `docs/develop/backend.md` 中同步记录原因。

## 分层映射

Java 服务不要按 Go 文件名机械翻译，应按职责映射：

| Go 分层 | Java 建议包 | 说明 |
| --- | --- | --- |
| `cmd/main.go` | `cn.harnesstemplate.admin.AdminApplication` | Spring Boot 启动类，实际包名确定后同步更新本文。 |
| `internal/config` | `config` | 配置 properties、外部依赖配置、条件开关。 |
| `internal/fiberc/middleware` | `web/filter`、`web/interceptor`、`web/advice` | timestamp、auth、casbin、language、encrypt、全局异常和响应包装。 |
| `internal/router/*` | `web/controller` | Controller 只负责 HTTP 协议和参数绑定。 |
| `internal/router/logic` | `application` 或 `service` | 业务编排层，承接现有 Handler 逻辑。 |
| `internal/appsvc` | `application/port`、`application/adapter` | 以接口隔离认证、调度、索引、文件存储等能力。 |
| DTO/entity 转换 | `application/converter` 或 `web/converter` | 统一使用 mapstruct-plus，避免转换逻辑散落在 Controller、Service 和 Repository 中。 |
| `internal/services/orm/models` | `domain/entity` | 数据库实体。 |
| `internal/services/orm/query` | `infrastructure/persistence` | Mapper、Repository、Query helper。 |
| `internal/services/*` | `infrastructure/*` | Redis、Temporal、Casbin、MinIO、Milvus、HTTP client。 |
| `internal/workflows` | `workflow` | Temporal workflow 和 activity。 |
| `internal/ai` | `ai` | LLM、RAG、工具、记忆和 Agent 编排。 |

推荐 Java 初始目录：

```text
backend/java/admin/
  build.gradle
  src/main/java/cn/harnesstemplate/admin/
    AdminApplication.java
    config/
    web/
      controller/
      filter/
      interceptor/
      advice/
      dto/
    application/
      port/
      service/
    domain/
      entity/
      model/
    infrastructure/
      persistence/
      redis/
      auth/
      permission/
      temporal/
      objectstore/
      milvus/
      http/
    workflow/
    ai/
  src/main/resources/
    application.yml
    db/migration/
  src/test/java/
```

## API 兼容要求

第一版 Java 服务应把“前端无感切换”作为核心验收标准。

- API 前缀保持 `/api`。
- 账号接口保持：
  - `POST /api/account/login/pwd`
  - `GET /api/account/logout`
  - `POST /api/account/changePwd`
- 加密公钥接口保持 `/api/encrypt/public/key`。
- 登录后业务路由覆盖用户、角色、字典、语言、资源菜单、资源 API、API 日志、登录日志、任务调度、任务执行、知识库、对象存储。
- 响应体继续使用业务 `code` 区分状态，成功为 `1`，通用失败为 `2`，登录/权限相关失败保持现有约定。
- HTTP 状态默认仍按现有服务兼容策略处理：业务失败不应随意改成前端未适配的非 2xx。
- 请求和响应字段名保持现有 JSON 名称，避免前端 DTO 和请求层同步大改。
- Swagger/OpenAPI 文档要能覆盖前端调用到的全部接口。

## 认证、加密和权限

这部分是移植风险最高的基础协议，必须优先迁移和测试。

1. 登录接口校验用户名和密码，返回 token 和本次会话 public key。
2. 服务端保存 session 信息，至少包含 private key、user id、username、role codes、role ids。
3. 登录后加密路由需要读取请求头中的加密 AES key、签名、时间戳和 request id。
4. Java 侧必须复刻 AES/RSA 加解密、AAD 拼接、响应加密和 `X-Response-Is-Encrypt` 语义。
5. `/api/storage/file/upload` 这类 multipart 上传接口保持 authenticated 但跳过 JSON body 加密。
6. Casbin 或等价权限中间件必须兼容当前角色、菜单、API 和策略数据。
7. 语言中间件继续保留默认语言和请求语言解析能力。

建议先为加密协议写跨语言 golden tests：使用固定 RSA key、AES key、timestamp、request id、query 和 body，确认 Go 与 Java 的密文、签名校验和解密结果一致。

## 数据模型迁移

Go 当前 model 位于 `backend/go/admin/internal/services/orm/models`，Java 实体应以这些表为第一来源：

| 表/实体范围 | 说明 |
| --- | --- |
| `sys_user`、`sys_role`、`sys_user_role` | 用户、角色和关系。 |
| `sys_resource_menu`、`sys_resource_api`、`sys_resource_menu_api` | 菜单、API 资源和关联。 |
| `sys_role_menu`、`sys_role_api` | 角色资源授权。 |
| `sys_data_permission` | 数据权限。 |
| `sys_dict_type`、`sys_dict_entry` | 字典类型与条目。 |
| `sys_language_type`、`sys_language_entry` | 多语言类型与文案。 |
| `sys_api_log`、`sys_login_log` | API 和登录日志。 |
| `job_schedule`、`job_execution` | 任务调度和执行记录。 |
| `file_asset` | 文件资产、上传状态和对象存储元数据。 |
| `knowledge_collection`、`knowledge_document` | 知识库 collection 和文档索引状态。 |
| `agent_session`、`agent_message` | AI Agent 会话和消息记忆。 |
| `casbin_rule` | Casbin 策略表。 |

迁移顺序：

1. 固化 schema：从 Go model 和当前数据库生成 DDL 快照。
2. Java 建 Easy Query entity、proxy 和 repository/query service，不改字段名、类型和索引。
3. 用 Flyway 接管后续 schema 迁移。
4. 为关键表补 repository 层测试，覆盖分页、软删除、唯一键、JSON 字段和时间字段。
5. 与 Go 服务共享同一套本地测试库跑只读查询，再逐步开启写路径。

## 功能迁移阶段

### 阶段 0：基线冻结

- 固定 Go 服务当前可用配置、Swagger、数据库 schema 和前端调用清单。
- 导出接口清单、请求/响应样例和关键错误样例。
- 确认 Java 服务端口，例如先使用 `3002`，避免和 Go 的 `3001` 冲突。

验收：

- Go 服务现有后端测试可运行。
- 前端主要页面的 API 调用清单已保存。
- 数据库 schema 快照和初始化数据可复现。

### 阶段 1：Java 工程骨架

- 建立 Java 项目、配置文件、日志、统一响应、全局异常、健康检查、OpenAPI。
- 提供 `/api/encrypt/public/key` 和一个无需复杂依赖的 smoke endpoint。
- 接入本地 Postgres 和 Redis 配置，但业务写入可以先不开放。

验收：

- Java 服务能独立启动。
- 健康检查和 OpenAPI 可访问。
- 配置缺失时有清晰错误或降级说明。

### 阶段 2：认证、加密和权限

- 迁移登录、登出、改密。
- 迁移 session 存储、Token Cookie、RSA/AES 加解密协议。
- 迁移 Casbin 权限校验和语言中间件。

验收：

- 前端可通过 Java 服务登录。
- 加密请求和响应能被现有前端请求层正常处理。
- 权限失败和未登录失败的业务 code 与 Go 侧兼容。

### 阶段 3：系统管理 CRUD

- 迁移用户、角色、菜单、API、字典、语言、API 日志、登录日志。
- 保持分页、排序、批量操作、树结构和翻译查询语义。
- 优先补 Controller contract tests 和 repository tests。

验收：

- 前端系统管理页面可切到 Java 服务完成日常操作。
- Java 与 Go 对同一查询样例返回兼容 JSON。

### 阶段 4：任务调度和执行

- 迁移 job schedule、job execution、Temporal client、worker 注册和 workflow/activity。
- 明确 Go worker 和 Java worker 是否共用 task queue；迁移期建议分 task queue，避免同名 workflow 的跨语言行为不清晰。
- 保持执行记录、取消、重试和手动触发语义。

验收：

- Java 可创建、触发、取消和重试任务。
- Temporal 不可用时的错误路径清楚，不吞掉任务状态。

### 阶段 5：对象存储和知识库

- 迁移文件上传、预签名直传、完成上传、详情、短链和删除。
- 迁移 knowledge collection/document CRUD、文件导入、索引状态和失败原因。
- 接入 MinIO、Milvus、embedding 配置。

验收：

- 前端可完成传统上传和直传流程。
- 文档导入后能写入 `knowledge_document` 并触发索引。
- 索引失败可在文档状态或错误字段中排查。

### 阶段 6：AI Agent

- 迁移 ChatModel 工厂、embedding、retriever、indexer、memory、工具调用和 Agent 编排。
- Java 生态不一定有 Eino 的直接等价物，可以先以应用服务接口隔离，实现最小对话和 RAG，再迁移 Plan-Execute-Replan。
- MCP SSE 工具、Prometheus 告警、数据库 CRUD 工具和内部文档搜索要分批接入。

验收：

- RAG 对话能基于 Milvus 检索结果回答。
- memory 后端至少支持 memory 和 Redis，DB 后端可按需求排期。
- 工具不可用时有可观测降级，而不是让整个 Agent 初始化失败。

## 配置映射

Java 配置应覆盖 Go 当前默认配置的同等含义：

| Go 配置 | Java 建议 | 说明 |
| --- | --- | --- |
| `Host`、`Port` | `server.address`、`server.port` | Java 迁移期建议先用 `3002`。 |
| `RestPrefix` | `admin.rest-prefix` | 默认 `/api`。 |
| `IsSwagger`、`SwaggerPrefix` | `knife4j.*` 或 `admin.openapi.*` | 控制 Knife4j/OpenAPI 暴露。 |
| `Orm.*` | `spring.datasource.*`、`admin.database.*` | 数据源和迁移配置。 |
| `Redis.*` | `spring.data.redis.*` | Redis 连接。 |
| `Temporal.*` | `admin.temporal.*` | namespace、host、task queue、worker 开关。 |
| `Storage.*` | `admin.storage.*` | enabled、engine、max upload bytes、presigned 过期时间。 |
| `Milvus.*` | `admin.milvus.*` | Milvus 地址、DB、collection。 |
| `AI.*` | `admin.ai.*` | embedding、model、memory、MCP URL。 |
| `Auth.TokenName` | `admin.auth.token-name` | 必须与前端 Cookie/请求层兼容。 |
| `DefaultLanguage` | `admin.default-language` | 语言中间件默认值。 |

配置文件应拆成本地默认、测试和生产模板，禁止把生产密钥写入仓库。

## 并行运行和切流

- Go 服务继续监听 `3001`。
- Java 服务迁移期监听 `3002`。
- 前端通过 `VITE_API_URL` 或反向代理选择目标后端。
- 模块级验证时可以只把部分路由转发到 Java，但必须保证 session、加密 key 和权限数据共享策略清楚。
- 切流前应保留 Go 回退路径，并记录当前数据库迁移版本。

推荐切流顺序：

1. 本地手动切 `VITE_API_URL` 到 Java。
2. 测试环境按路由组切流。
3. 系统管理主链路全量 Java。
4. 文件、知识库和任务链路分批 Java。
5. AI Agent 最后切流。

## 测试与验收

最低验证集合：

- Java 单元测试：配置绑定、加密协议、认证服务、权限服务、repository。
- Java 集成测试：Postgres、Redis、MinIO、Temporal 可用路径。
- API contract tests：对比 Go 和 Java 的关键接口响应结构。
- 前端回归：`front/apps/admin-react` 登录、系统管理、任务、文件、知识库页面。
- 安全测试：未登录、权限不足、过期 timestamp、错误签名、错误 AES key、multipart 跳过加密。
- 可观测性测试：健康检查、metrics、关键错误日志、启动失败日志。

建议保留一组跨语言兼容样例：

| 样例 | 用途 |
| --- | --- |
| 登录成功/失败请求响应 | 校验 token、public key、业务 code。 |
| 加密 POST 请求 | 校验 AES/RSA/AAD/签名兼容。 |
| 权限拒绝请求 | 校验业务 code 和错误响应。 |
| 分页列表请求 | 校验分页字段和排序语义。 |
| 文件直传 prepare/complete | 校验对象存储状态机。 |
| 知识文档索引失败 | 校验错误记录和排障路径。 |

## 风险清单

| 风险 | 影响 | 缓解 |
| --- | --- | --- |
| 加密协议不兼容 | 前端所有登录后接口不可用 | 先写 golden tests，再接业务接口。 |
| ORM 行为差异 | 分页、空值、软删除或事务行为变化 | 用真实表结构和关键查询样例做 contract tests。 |
| 权限模型偏差 | 菜单/API 权限错误放行或误拒 | 优先迁移 Casbin 数据兼容测试。 |
| Temporal workflow 名称冲突 | Go/Java worker 抢占任务或反序列化失败 | 迁移期使用独立 task queue 或版本化 workflow。 |
| 对象存储状态机不一致 | 文件资产脏数据或直传完成失败 | 用状态流测试覆盖 pending/active/deleted。 |
| AI 生态差异 | Agent 能力无法直接翻译 | 先保接口和配置边界，分批替换底层编排实现。 |
| 文档和代码脱节 | 后续 Agent 无法稳定执行 | 每完成一个阶段同步更新本文、`docs/develop/backend.md` 和 history。 |

## 交付检查清单

每个阶段完成前至少确认：

- 相关 Java 代码、配置、测试和文档同源更新。
- 新增或变化的外部依赖已记录默认值、启动方式和故障表现。
- 公共 API 变化已同步 OpenAPI。
- 涉及认证、加密、权限、Cookie、密钥或外部凭据的改动已同步 `docs/operate/security.md`。
- 实质代码变更已补 `docs/records/histories/`。
- 若服务路径正式纳入仓库矩阵，已更新 `.service-matrix/dependencies.yaml`。
