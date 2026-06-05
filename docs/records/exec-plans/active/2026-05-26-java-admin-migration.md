# Java Admin Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `backend/java/admin` as a Java implementation of the current Go Admin service while preserving the existing frontend API contract, authentication, encrypted request protocol, database model, and staged rollback path.

**Architecture:** The Java service is a parallel Spring Boot 3.x service, initially running on port `3002` while the Go service keeps port `3001`. The implementation follows the migration guide in `docs/develop/java-admin-migration.md`: Controller handles HTTP protocol, application services own business orchestration, Easy Query owns persistence, infrastructure packages isolate Redis, Sa-Token, jCasbin, Temporal, MinIO, Milvus, and AI integrations. Each stage must be independently testable and must not require deleting or disabling `backend/go/admin`.

**Tech Stack:** Spring Boot 3.x latest stable line, Java 25, Gradle, Flyway, Lombok, linpeilie/mapstruct-plus, dromara/easy-query v3, Sa-Token Java, jCasbin, Spring Data Redis, Temporal Java SDK, LangChain4j, MinIO Java SDK, Knife4j, Micrometer, Actuator.

---

## 范围

包含：

- 在 `backend/java/admin` 新建 Java Admin 工程骨架。
- 迁移 Go Admin 的配置、统一响应、异常处理、认证、加密、权限、系统 CRUD、任务、文件、知识库和 AI 链路。
- 固化 Go 与 Java 的 API contract、加密协议 golden cases、数据库 schema 和前端回归路径。
- 同步更新 `docs/develop/java-admin-migration.md`、`docs/develop/backend.md`、`docs/operate/reliability.md`、`docs/operate/security.md`、`.service-matrix/dependencies.yaml` 和 history。

不包含：

- 删除 `backend/go/admin`。
- 重做 `front/apps/admin-react` 的请求层、页面结构或加密协议。
- 大改现有数据库表名、字段名、索引和初始化数据。
- 一次性强行迁移全部 AI Agent 能力；AI 按最后阶段独立验收。

## 背景

相关文档：

- `docs/develop/java-admin-migration.md`
- `docs/develop/backend.md`
- `docs/develop/architecture.md`
- `docs/operate/reliability.md`
- `docs/operate/security.md`
- `docs/govern/harness-process.md`
- `docs/govern/plans.md`
- `context/project/admin/INDEX.md`

相关代码路径：

- Go 服务：`backend/go/admin`
- Go model：`backend/go/admin/internal/services/orm/models`
- Go 路由：`backend/go/admin/internal/router`
- Go 中间件：`backend/go/admin/internal/fiberc/middleware`
- Go 业务 Handler：`backend/go/admin/internal/router/logic`
- Java 目标目录：`backend/java/admin`
- 管理后台前端：`front/apps/admin-react`
- 服务矩阵：`.service-matrix/dependencies.yaml`

已知约束：

- Java 服务迁移期默认端口为 `3002`，Go 服务保持 `3001`。
- API 前缀保持 `/api`。
- 登录、Token、Cookie、RSA/AES 加密、业务响应 `code` 和权限失败语义必须兼容现有前端。
- 数据库 schema 先由 Go model 和当前数据库快照固定，Java 后续迁移由 Flyway 接管。
- Java 25 需要在第一阶段验证 Gradle toolchain、CI 镜像和依赖库兼容性。

## 文件结构计划

### 新建 Java 工程

- Create: `backend/java/admin/settings.gradle`
- Create: `backend/java/admin/build.gradle`
- Create: `backend/java/admin/gradle/wrapper/gradle-wrapper.properties`
- Create: `backend/java/admin/gradlew`
- Create: `backend/java/admin/gradlew.bat`
- Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/AdminApplication.java`
- Create: `backend/java/admin/src/main/resources/application.yml`
- Create: `backend/java/admin/src/main/resources/application-local.yml`
- Create: `backend/java/admin/src/main/resources/db/migration/V1__baseline_schema.sql`
- Create: `backend/java/admin/src/test/resources/application-test.yml`

### Java 包边界

- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/config`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/dto`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/interceptor`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/advice`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/port`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/converter`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/domain/entity`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/domain/model`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/persistence`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/redis`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/auth`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/permission`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/temporal`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/objectstore`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/milvus`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/http`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/workflow`
- Create directory: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai`

### 仓库文档和矩阵

- Modify: `.service-matrix/dependencies.yaml`
- Modify: `docs/develop/architecture.md`
- Modify: `docs/develop/backend.md`
- Modify: `docs/operate/reliability.md`
- Modify: `docs/operate/security.md`
- Modify: `context/project/admin/INDEX.md`
- Create stage history files under `docs/records/histories/2026-05/`, following `docs/govern/histories.md` naming rules at the time each stage closes.

## 风险

| 风险 | 影响 | 缓解方式 |
| --- | --- | --- |
| Java 25 与依赖不兼容 | 工程无法构建或运行时异常 | 阶段 1 先跑 Gradle toolchain、`./gradlew test` 和 smoke test；失败时记录具体依赖并决定降级或等待兼容。 |
| 加密协议不兼容 | 前端登录后接口全部不可用 | 阶段 2 先写 golden tests，固定 RSA、AES、timestamp、request id、query、body 和期望结果。 |
| 业务响应 code 不兼容 | 前端统一错误处理失效 | 阶段 1 建统一响应 contract tests；每个 Controller 切片都先补 contract test。 |
| ORM 行为差异 | 分页、空值、时间、JSON、软删除行为变化 | 阶段 3 以 Go model 和数据库快照生成 Flyway baseline，并为 Easy Query repository 写集成测试。 |
| 权限模型偏差 | 菜单/API 误放行或误拒 | 阶段 2 建 jCasbin policy 加载测试和路由权限拒绝测试。 |
| Temporal workflow 跨语言冲突 | Go/Java worker 抢任务或反序列化失败 | 迁移期 Java 使用独立 task queue，例如 `admin-java`，切流前显式记录 workflow 命名和版本策略。 |
| 文件状态机不一致 | 直传完成失败或脏数据 | 阶段 5 覆盖 `pending_upload`、`active`、`deleted` 和 presigned URL 过期行为。 |
| AI 生态差异 | Go Eino 编排无法直接翻译 | 阶段 6 先保应用接口和 RAG 最小链路，再补 Plan-Execute-Replan。 |

## 里程碑

### Milestone 0: 基线冻结

- [ ] **Step 0.1: 固化 Go API contract 清单**

  Files:

  - Create: `backend/java/admin/docs/contracts/go-admin-api-contract.md`

  Commands:

  ```bash
  cd backend/go/admin
  make swagger
  ```

  Expected:

  - `backend/go/admin/docs/swagger.json` 和 `backend/go/admin/docs/swagger.yaml` 是当前最新生成结果。
  - `backend/java/admin/docs/contracts/go-admin-api-contract.md` 记录账号、加密、系统管理、任务、文件、知识库和 AI 接口清单。

- [ ] **Step 0.2: 固化数据库 schema baseline**

  Files:

  - Create: `backend/java/admin/docs/contracts/go-admin-schema-baseline.md`
  - Create: `backend/java/admin/src/main/resources/db/migration/V1__baseline_schema.sql`

  Commands:

  ```bash
  cd backend/go/admin
  make script-orm
  ```

  Expected:

  - Go ORM 生成产物与 model 同步。
  - `V1__baseline_schema.sql` 包含当前 Java 需要承接的表、索引和初始化约束。

- [ ] **Step 0.3: 固化前端关键调用样例**

  Files:

  - Create: `backend/java/admin/docs/contracts/frontend-api-samples.md`

  Content requirements:

  - 登录成功和失败样例。
  - 加密 POST 样例。
  - 分页列表样例。
  - 权限拒绝样例。
  - 文件 prepare/complete 样例。
  - 知识文档索引失败样例。

### Milestone 1: Java 工程骨架

- [ ] **Step 1.1: 创建 Gradle Spring Boot 工程**

  Files:

  - Create: `backend/java/admin/settings.gradle`
  - Create: `backend/java/admin/build.gradle`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/AdminApplication.java`

  Test first:

  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/AdminApplicationTests.java`
  - Test asserts Spring context starts with `test` profile.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test
  ```

  Expected before implementation:

  - FAIL because Gradle project or Spring application class does not exist.

  Expected after implementation:

  - PASS with context load test.

- [ ] **Step 1.2: 接入配置绑定和启动端口**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/config/AdminProperties.java`
  - Modify: `backend/java/admin/src/main/resources/application.yml`
  - Modify: `backend/java/admin/src/main/resources/application-local.yml`

  Test first:

  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/config/AdminPropertiesTests.java`
  - Assert default `restPrefix=/api`, `server.port=3002`, `auth.tokenName=token`, `defaultLanguage=zh-CN`.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*AdminPropertiesTests'
  ```

  Expected:

  - The test fails before properties exist and passes after configuration binding is implemented.

- [ ] **Step 1.3: 建统一响应、异常处理和 smoke endpoint**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/dto/Response.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/advice/GlobalExceptionHandler.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/HealthController.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/HealthControllerTests.java`

  Test first:

  - `GET /api/health/smoke` returns HTTP 200.
  - JSON body contains `code=1`, `msg=success`, and non-null `data`.
  - A controlled exception returns HTTP 200 with `code=2`.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*HealthControllerTests'
  ```

  Expected:

  - Tests fail before endpoint and response wrapper exist, then pass.

- [ ] **Step 1.4: 接入 Actuator、Micrometer 和 Knife4j**

  Files:

  - Modify: `backend/java/admin/build.gradle`
  - Modify: `backend/java/admin/src/main/resources/application.yml`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/config/OpenApiAndActuatorTests.java`

  Test first:

  - `GET /actuator/health` returns UP in test profile.
  - OpenAPI docs endpoint is exposed in local profile.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*OpenApiAndActuatorTests'
  ```

  Expected:

  - Tests fail before dependencies/config exist and pass after wiring.

### Milestone 2: 认证、加密和权限基础

- [ ] **Step 2.1: 迁移加密协议 golden tests**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/auth/CryptoService.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/auth/CryptoServiceTests.java`
  - Create: `backend/java/admin/src/test/resources/golden/crypto-request.json`

  Test first:

  - Use fixed RSA private key, RSA public key, AES key, timestamp, request id, query string and request body.
  - Assert RSA decrypts request AES key.
  - Assert AES decrypts encrypted body.
  - Assert AAD string sorting matches Go behavior.
  - Assert response encryption can be decrypted by the same AES key.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*CryptoServiceTests'
  ```

  Expected:

  - Tests fail before `CryptoService` exists and pass only when Go-compatible encryption is implemented.

- [ ] **Step 2.2: 实现 Sa-Token 登录会话适配**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/AccountService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/auth/SaTokenSessionService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/AccountController.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/AccountControllerTests.java`

  Test first:

  - `POST /api/account/login/pwd` invalid username returns HTTP 200 and business failure code.
  - Valid login returns `token` and `publicKey`.
  - Session stores private key, user id, username, role codes and role ids.
  - `GET /api/account/logout` clears session.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*AccountControllerTests'
  ```

  Expected:

  - Tests fail before controller/service exist and pass after Sa-Token integration.

- [ ] **Step 2.3: 实现请求过滤器链**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter/TimestampFilter.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter/EncryptFilter.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/interceptor/LanguageInterceptor.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/filter/SecurityFilterTests.java`

  Test first:

  - Expired timestamp returns business request-expired code.
  - Missing encrypted AES key on encrypted route returns request error code.
  - Multipart upload path is authenticated but not JSON-body encrypted.
  - Language header overrides default language in request context.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*SecurityFilterTests'
  ```

  Expected:

  - Tests fail before filters/interceptor exist and pass after wiring.

- [x] **Step 2.4: 实现 jCasbin 权限拦截**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/permission/CasbinPermissionService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/interceptor/PermissionInterceptor.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/permission/CasbinPermissionServiceTests.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/interceptor/PermissionInterceptorTests.java`

  Test first:

  - Load sample `casbin_rule` policies.
  - Role with API permission can access route.
  - Role without API permission receives authorization failure business code.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*CasbinPermissionServiceTests' --tests '*PermissionInterceptorTests'
  ```

  Expected:

  - Tests fail before permission service exists and pass after jCasbin integration.

### Milestone 3: 数据模型和系统管理 CRUD

- [x] **Step 3.1: 建 Easy Query 实体和 Flyway baseline**

  Files:

  - Create entity files under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/domain/entity`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/persistence/SchemaBaselineTests.java`
  - Modify: `backend/java/admin/src/main/resources/db/migration/V1__baseline_schema.sql`

  Test first:

  - Flyway migrates an empty test database.
  - Tables include `sys_user`, `sys_role`, `sys_resource_menu`, `sys_resource_api`, `sys_dict_type`, `sys_dict_entry`, `job_schedule`, `job_execution`, `file_asset`, `knowledge_collection`, `knowledge_document`, `agent_session`, `agent_message`, `casbin_rule`.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*SchemaBaselineTests'
  ```

  Expected:

  - Tests fail before migration and entity metadata exist and pass after baseline is complete.

- [x] **Step 3.2: 建 mapstruct-plus 转换层**

  Files:

  - Create converter files under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/converter`
  - Create DTO files under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/dto`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/application/converter/ConverterTests.java`

  Test first:

  - Entity to response DTO keeps JSON-facing field names.
  - Request DTO to entity ignores server-owned fields such as id, created time and updated time where the route owns them.
  - Sensitive fields such as password, token and private key are not present in public response DTOs.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*ConverterTests'
  ```

  Expected:

  - Tests fail before converters exist and pass after mapstruct-plus annotations and generated converters are wired.

- [ ] **Step 3.3: 迁移用户、角色、菜单和 API 资源 CRUD**

  Files:

  - Create controllers under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/system`
  - Create services under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system`
  - Create repositories/query services under `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/persistence/system`
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/system`

  Test first:

  - User list supports pagination and username filtering.
  - Role create/update/delete preserves role-menu and role-api relations.
  - Menu tree returns stable parent-child structure.
  - API resource list matches permission route naming.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.web.controller.system.*'
  ```

  Expected:

  - Tests fail route by route before implementation and pass after each CRUD slice.

- [ ] **Step 3.4: 迁移字典、语言、API 日志和登录日志**

  Files:

  - Create controllers/services/repositories for dict, language, api log and login log.
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/system`

  Test first:

  - Dict entry match returns entries by type and key.
  - Language entry update preserves type relation.
  - API log list filters by route, method and time range.
  - Login log records success and failure status.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.web.controller.system.*'
  ```

  Expected:

  - Tests fail before endpoints exist and pass after implementation.

### Milestone 4: 任务调度和执行

- [ ] **Step 4.1: 接入 Temporal Java SDK**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/temporal/TemporalProperties.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/temporal/TemporalClientFactory.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/temporal/TemporalClientFactoryTests.java`

  Test first:

  - Worker disabled profile does not create a worker.
  - Default task queue is `admin-java`.
  - Client configuration carries namespace and host port from properties.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*TemporalClientFactoryTests'
  ```

  Expected:

  - Tests fail before properties/factory exist and pass after implementation.

- [ ] **Step 4.2: 迁移 job schedule 和 execution API**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/job/JobScheduleController.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/job/JobExecutionController.java`
  - Create services/repositories under `application/service/job` and `infrastructure/persistence/job`
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/job`

  Test first:

  - Schedule create/update/switch/trigger works against test DB.
  - Execution list/detail/cancel/retry records expected state transitions.
  - Temporal unavailable returns controlled business failure and persists failure reason.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.web.controller.job.*'
  ```

  Expected:

  - Tests fail before job endpoints exist and pass after staged implementation.

### Milestone 5: 文件上传和知识库

- [ ] **Step 5.1: 接入 MinIO 文件资产状态机**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/objectstore/ObjectStorageService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/storage/FileAssetService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/storage/StorageFileController.java`
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/storage`

  Test first:

  - `prepareUpload` creates `pending_upload`.
  - `completeUpload` verifies object exists and switches to `active`.
  - `presigned` respects expiry bounds.
  - `del` marks asset deleted and deletes or invalidates object according to configured policy.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.web.controller.storage.*'
  ```

  Expected:

  - Tests fail before storage services exist and pass after MinIO integration.

- [ ] **Step 5.2: 迁移 knowledge collection/document CRUD**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/knowledge/KnowledgeCollectionController.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/knowledge/KnowledgeDocumentController.java`
  - Create services/repositories under `application/service/knowledge` and `infrastructure/persistence/knowledge`
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/knowledge`

  Test first:

  - Collection list/detail/create/update/delete matches Go response shape.
  - Document create/update/delete records indexing status.
  - `importFile` reads completed `file_asset` and rejects unsupported content type.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.web.controller.knowledge.*'
  ```

  Expected:

  - Tests fail before knowledge endpoints exist and pass after CRUD and import flow are implemented.

- [ ] **Step 5.3: 接入 Milvus 和 embedding 配置**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/milvus/MilvusVectorStore.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai/EmbeddingService.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/ai/KnowledgeIndexingTests.java`

  Test first:

  - Missing embedding API key fails with explicit configuration error.
  - Document indexing failure records `indexing_error`.
  - Milvus collection dimension mismatch returns actionable error.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*KnowledgeIndexingTests'
  ```

  Expected:

  - Tests fail before indexing services exist and pass after explicit success/failure behavior is implemented.

### Milestone 6: AI Agent

- [ ] **Step 6.1: 建 AI 应用接口和最小 RAG 对话**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai/ChatModelFactory.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai/RagChatService.java`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/ai/AiChatController.java`
  - Create: `backend/java/admin/src/test/java/cn/harnesstemplate/admin/ai/RagChatServiceTests.java`

  Test first:

  - Missing model alias returns controlled configuration error.
  - RAG service injects retrieved document snippets into prompt context.
  - MCP URL missing disables log tool without failing chat model creation.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests '*RagChatServiceTests'
  ```

  Expected:

  - Tests fail before AI service exists and pass after minimum RAG path is implemented.

- [ ] **Step 6.2: 迁移 memory 和工具调用**

  Files:

  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai/memory`
  - Create: `backend/java/admin/src/main/java/cn/harnesstemplate/admin/ai/tools`
  - Create tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/ai`

  Test first:

  - Memory supports in-memory and Redis backends.
  - DB memory persists `agent_session` and `agent_message`.
  - Tool registry tolerates unavailable optional tools and logs a warning.

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test --tests 'cn.harnesstemplate.admin.ai.*'
  ```

  Expected:

  - Tests fail before memory/tools exist and pass after staged implementation.

### Milestone 7: 切流、文档和交付

- [ ] **Step 7.1: 纳入服务矩阵和文档入口**

  Files:

  - Modify: `.service-matrix/dependencies.yaml`
  - Modify: `docs/develop/architecture.md`
  - Modify: `docs/develop/backend.md`
  - Modify: `docs/operate/reliability.md`
  - Modify: `docs/operate/security.md`
  - Modify: `context/project/admin/INDEX.md`

  Required content:

  - `admin-java-backend` service entry with `repo_path: backend/java/admin`.
  - Java startup command.
  - Java test command.
  - External dependency defaults.
  - Authentication, encryption and secret handling notes.

- [ ] **Step 7.2: 前端联调和切流验证**

  Files:

  - Modify only if needed: `front/apps/admin-react/.env.local.example`
  - Create: `backend/java/admin/docs/contracts/frontend-regression-report.md`

  Commands:

  ```bash
  cd backend/java/admin
  ./gradlew test
  cd ../../../front/apps/admin-react
  VITE_API_URL=http://127.0.0.1:3002 pnpm e2e:test
  ```

  Expected:

  - Java tests pass.
  - Frontend login, system management, task, file and knowledge pages pass against Java backend.
  - If external dependencies are unavailable, the report records skipped cases, dependency names and retry commands.

- [ ] **Step 7.3: 记录 history 并关闭阶段**

  Files:

  - Create a stage history file under `docs/records/histories/2026-05/`, following `docs/govern/histories.md` naming rules at completion time.
  - Modify: `docs/records/exec-plans/active/2026-05-26-java-admin-migration.md`

  Required content:

  - User request summary.
  - Stage delivered.
  - Files changed.
  - Commands run and results.
  - Known risks and follow-up actions.

## 测试计划

### 测试文件

- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/AdminApplicationTests.java`: Spring context and profile boot.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/config/AdminPropertiesTests.java`: configuration binding.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/HealthControllerTests.java`: unified response and smoke endpoint.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/auth/CryptoServiceTests.java`: RSA/AES/AAD golden cases.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/AccountControllerTests.java`: login, logout, change password.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/filter/SecurityFilterTests.java`: timestamp, encryption and multipart skip behavior.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/permission/CasbinPermissionServiceTests.java`: policy loading and enforcement.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/persistence/SchemaBaselineTests.java`: Flyway schema and Easy Query entity compatibility.
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/application/converter/ConverterTests.java`: mapstruct-plus conversion behavior.
- Controller tests under `web/controller/system`, `web/controller/job`, `web/controller/storage`, `web/controller/knowledge`, and `web/controller/ai`: API contract compatibility.
- AI tests under `backend/java/admin/src/test/java/cn/harnesstemplate/admin/ai`: RAG, memory and tools.

### 先写失败测试

Each implementation slice begins with the test file named in its milestone. The first run must fail for the expected reason:

- Missing class or bean for new package slices.
- Missing endpoint for Controller slices.
- Missing migration table for schema slices.
- Missing policy or permission service for jCasbin slices.
- Missing external client factory for Temporal, MinIO, Milvus and AI slices.

Run the narrow command listed in the milestone first, then implement the smallest code needed for that command to pass.

### 通过标准

- Stage-level narrow tests pass before moving to the next stage.
- Before any stage is marked complete, run:

  ```bash
  cd backend/java/admin
  ./gradlew test
  ```

- Before cutover validation, run:

  ```bash
  cd backend/go
  go test ./...
  cd ../../front/apps/admin-react
  VITE_API_URL=http://127.0.0.1:3002 pnpm e2e:test
  ```

- If Postgres, Redis, Temporal, MinIO, Milvus, embedding provider or ChatModel provider is unavailable, record the skipped test group, dependency, command output summary and retry command in the stage history.

## 验证方式

Commands:

```bash
make ci
cd backend/java/admin && ./gradlew test
cd backend/go && go test ./...
cd front/apps/admin-react && VITE_API_URL=http://127.0.0.1:3002 pnpm e2e:test
```

Manual checks:

- Java service starts on `3002`; Go service remains available on `3001`.
- `GET /actuator/health` reports healthy when dependencies are available.
- Knife4j/OpenAPI shows account, encryption and migrated business routes.
- Frontend can log in, read menus, call encrypted endpoints and complete basic CRUD.

Observability checks:

- Authentication failure, permission failure, encryption failure and dependency failure emit clear logs.
- Prometheus metrics endpoint is exposed through Actuator.
- Temporal, MinIO, Milvus and AI initialization failures are explicit and actionable.

## 进度记录

- [x] 2026-05-26: 移植指南已建立在 `docs/develop/java-admin-migration.md`。
- [x] 2026-05-26: Java 技术栈已固定为 Spring Boot 3.x、Java 25、Gradle、Flyway、Lombok、mapstruct-plus、easy-query v3、Sa-Token、jCasbin、Spring Data Redis、Temporal、LangChain4j、MinIO、Knife4j、Micrometer/Actuator。
- [x] 2026-05-26: Java admin 已移除主代码中的 JdbcTemplate 业务持久化，改为经 Easy Query client 适配层访问；同时补齐 Flyway baseline、Easy Query 方言和 LangChain4j 依赖/调用边界。
- [ ] Milestone 0: 基线冻结。
- [ ] Milestone 1: Java 工程骨架。
- [ ] Milestone 2: 认证、加密和权限基础。
- [ ] Milestone 3: 数据模型和系统管理 CRUD。
- [ ] Milestone 4: 任务调度和执行。
- [ ] Milestone 5: 文件上传和知识库。
- [ ] Milestone 6: AI Agent。
- [ ] Milestone 7: 切流、文档和交付。

## 决策记录

- 2026-05-26: Java 版 admin 作为并行服务迁移，不替换或删除 Go 服务；原因是认证、加密、数据和任务链路风险高，需要可回退路径。
- 2026-05-26: Spring Boot 固定在 3.x 最新稳定线，不跟随 Spring Boot 4.x；原因是当前迁移目标优先稳定兼容。
- 2026-05-26: Java 运行时选 Java 25；实现第一阶段必须验证 Gradle toolchain、CI 镜像和依赖兼容性。
- 2026-05-26: ORM 使用 `dromara/easy-query` v3 最新稳定版；原因是用户指定，并且其 entity/proxy/query 模式适合承接现有 GORM model 和查询。
- 2026-05-26: 类转换使用 `linpeilie/mapstruct-plus`；原因是统一 DTO、VO、request、response 和 entity 转换边界，减少手写拷贝。
- 2026-05-26: API 文档使用 Knife4j；原因是用户指定，迁移时仍需保持 OpenAPI/Swagger contract 可核对。
- 2026-05-26: AI 编排使用 LangChain4j；原因是用户指定，后续 RAG、ChatModel、工具调用和 Agent 迁移以 LangChain4j 为边界。
- 2026-05-26: 对已有 Go schema 启用 Flyway `baseline-on-migrate`；原因是 `public` schema 非空但没有 `flyway_schema_history` 时，Flyway 默认会拒绝启动。
- 2026-05-26: Easy Query 主键字段统一标记 `generatedKey = true`，插入通过 `executeRows(true)` 回填；原因是 PostgreSQL `BIGSERIAL` id 需要由 ORM 显式回写到实体，后续更新、删除和查询依赖该 id。
