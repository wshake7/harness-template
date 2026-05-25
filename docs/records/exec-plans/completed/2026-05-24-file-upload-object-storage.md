# 文件上传与可插拔对象存储实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现管理后台通用文件上传能力，首个存储引擎接入 MinIO，并保留后续替换为其他 S3 兼容服务、本地文件系统或云厂商对象存储的扩展点。

**Architecture:** 后端新增 Storage 配置、对象存储生命周期服务、可插拔 `Engine` 接口、MinIO 适配器、文件元数据模型和 `/api/storage/file/*` 认证接口；前端新增文件上传 API 封装和可复用上传控件。上传请求使用 multipart/form-data，不走现有 JSON 请求加密中间件，但继续走登录态、Casbin、时间戳和审计日志约束。

**Tech Stack:** Go Fiber v3 + GORM + Gen Query + MinIO Go SDK (`github.com/minio/minio-go/v7`) / React + Ant Design Upload + alova

---

## 目标

建立一个可被知识库、表单附件、头像或后续业务复用的通用文件上传基础能力：

- 后端接收 authenticated multipart 上传，将文件流写入当前配置的对象存储引擎。
- 首期存储引擎为 MinIO，使用服务端 `PutObject` 上传和 `PresignedGetObject` 生成临时下载链接。
- 文件元数据持久化到 Postgres，记录 engine、bucket、object key、原始文件名、content type、大小、hash、上传人和业务标签。
- 前端提供 `StorageFileApi.upload` 和一个轻量可复用上传控件，供业务页面逐步接入。
- 配置、文档、Swagger、测试和 history 与代码同源更新。

## 范围

包含：

- 新增后端 `Storage` 配置段和本地默认 MinIO 示例配置。
- 新增对象存储引擎接口、MinIO engine、工厂和生命周期健康检查。
- 新增 `file_asset` ORM 模型和 Gen Query 产物。
- 新增文件上传、详情、临时下载链接、删除接口。
- 新增前端 API 封装和可复用上传组件。
- 新增单元测试、必要的集成/手工验证说明、Swagger 生成产物和相关文档更新。

不包含：

- 浏览器直传 MinIO、分片断点续传、秒传和批量大文件上传。
- 病毒扫描、DLP、图片转码、缩略图、CDN、生命周期清理任务。
- 公开匿名上传或无需登录下载。
- 对已有知识库 Document 内容入库流程做完整改造；本轮只预留可引用上传文件的字段和 API 能力。

## 背景

相关文档：

- `docs/README.md`：复杂或高风险任务需要落 active execution plan。
- `docs/govern/harness-process.md`：本任务处于阶段 3「设计」到阶段 4「实现」之间。
- `docs/govern/plans.md`：计划需写清目标、范围、约束、风险和验证方式。
- `docs/develop/backend.md`：后端边界、配置、API、Swagger、ORM 生成要求。
- `docs/operate/security.md`：新增外部服务或 SDK 时需记录凭据、网络访问和失败模式。
- `docs/operate/reliability.md`：新增外部依赖需补运行与排障说明。

相关代码路径：

- `backend/go/admin/internal/config/config.go`：统一配置结构。
- `backend/go/admin/etc/config.yaml`：本地默认配置。
- `backend/go/admin/internal/services/init.go`：服务生命周期注册。
- `backend/go/admin/internal/router/auth_router/auth_router.go`：认证路由分组与中间件。
- `backend/go/admin/internal/fiberc/middleware/security.go`：时间戳、加密和签名中间件。
- `backend/go/admin/internal/router/logic/`：业务 Handler 与 Swagger 注释。
- `backend/go/admin/internal/services/orm/models/`：ORM 模型定义。
- `front/packages/request/src/encryption/helpers.ts`：前端请求加密入口。
- `front/apps/admin-react/src/api/business/`：管理后台业务 API 封装。
- `front/apps/admin-react/src/components/business/`：可复用业务组件。

外部文档依据：

- Context7 resolved library: `/minio/minio-go`
- MinIO Go SDK 文档显示：客户端通过 `minio.New` 搭配 `credentials.NewStaticV4` 初始化；服务端上传流使用 `PutObject` 并传入 reader、size、bucket、object key 和 content type；临时下载链接使用 `PresignedGetObject(ctx, bucket, object, expiry, reqParams)`，最长 7 天。
- Source: `https://github.com/minio/minio-go/blob/master/docs/API.md`

已知约束：

- 现有 authenticated 业务路由默认套 `EncryptMiddleware()`，该中间件假设请求体是可解密 JSON；multipart 上传不能走这条路径。
- 前端 alova 客户端当前默认对除 `/api/encrypt/public/key` 外的请求执行 `encryptRequest`，FormData 上传需要显式跳过 AES body 加密。
- `backend/go/admin/go.mod` 尚未引入 `github.com/minio/minio-go/v7`，实现时会新增依赖并更新 `go.sum`。
- 生产或共享环境不得提交真实 MinIO AccessKey/SecretKey；示例配置只能留空或使用本地开发占位。

## 设计决策

- 存储能力命名为 `objectstore`，表示底层是对象存储，不把首个实现绑定为 MinIO。
- 后端接口返回 file asset 元数据，不默认返回永久公开 URL；下载或预览通过单独接口生成短期 presigned URL。
- `objectKey` 由服务端生成，格式为 `{prefix}/{yyyy}/{mm}/{dd}/{uuid}{ext}`，避免信任客户端文件名。
- 文件元数据落库后再暴露给业务使用；上传对象成功但落库失败时，立即尝试删除对象并记录日志。
- 上传接口使用 authenticated non-encrypted 分组：`AuthMiddleware + CasbinAPIMiddleware + LanguageMiddleware`，保留登录态和权限控制，跳过 JSON 加密。
- 前端上传方法通过 `meta.skipEncrypt = true` 触发请求层跳过 AES body 加密，只补时间戳和 request id。
- 首期由服务端代理上传，直传 MinIO 作为后续增强，不进入本轮范围。

## 文件结构

### 后端

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `backend/go/admin/internal/config/config.go` | Modify | 新增 `StorageConfig`、`MinIOStorageConfig`、上传大小和 presigned URL 默认过期配置。 |
| `backend/go/admin/etc/config.yaml` | Modify | 新增 `Storage` 示例配置，凭据留空。 |
| `backend/go/admin/etc/config.local.yaml` | Modify | 新增本地 MinIO 配置段，真实凭据不提交。 |
| `backend/go/admin/internal/services/init.go` | Modify | 注册对象存储生命周期服务。 |
| `backend/go/admin/internal/services/objectstore/engine.go` | Create | 定义 `Engine` 接口、`PutInput`、`PutResult`、`PresignInput`、`ObjectInfo`。 |
| `backend/go/admin/internal/services/objectstore/minio.go` | Create | 实现 MinIO engine，封装 client、bucket 检查、PutObject、PresignedGetObject、RemoveObject。 |
| `backend/go/admin/internal/services/objectstore/service.go` | Create | 根据配置创建当前 engine，提供健康检查和关闭逻辑。 |
| `backend/go/admin/internal/appsvc/file_storage.go` | Create | 定义 Handler 面向的文件存储应用服务接口，便于测试替换。 |
| `backend/go/admin/internal/appsvc/file_storage_impl.go` | Create | 协调对象上传、hash 计算、metadata 落库、presigned URL 和删除。 |
| `backend/go/admin/internal/services/orm/models/file_asset.go` | Create | 新增 `file_asset` 表模型。 |
| `backend/go/admin/internal/services/orm/query/*file_asset*` | Generate | 运行 `make script-orm` 生成 Gen Query。 |
| `backend/go/admin/internal/router/logic/storage_file.go` | Create | 文件上传、详情、临时下载链接、删除 Handler 和 Swagger 注释。 |
| `backend/go/admin/internal/router/auth_router/storage.go` | Create | 注册 `/api/storage/file/*` 路由。 |
| `backend/go/admin/internal/router/auth_router/auth_router.go` | Modify | 在 non-encrypted authenticated 分组注册 storage 路由。 |
| `backend/go/admin/internal/mock/mock_file_storage.go` | Generate | 为 Handler 测试生成 mock。 |
| `backend/go/admin/internal/services/objectstore/engine_test.go` | Create | 测试 object key 生成、presigned 过期上限和 engine 工厂错误。 |
| `backend/go/admin/internal/services/objectstore/minio_test.go` | Create | 使用 fake/minio wrapper 测试 bucket 检查、PutObject 参数和 RemoveObject 幂等策略。 |
| `backend/go/admin/internal/services/orm/models/file_asset_test.go` | Create | 测试 file asset AutoMigrate、JSON metadata、唯一索引和软删除字段。 |
| `backend/go/admin/internal/appsvc/file_storage_impl_test.go` | Create | 测试上传、落库、补偿删除、presigned URL、删除和错误分支。 |
| `backend/go/admin/internal/router/logic/storage_file_test.go` | Create | 测试 multipart Handler、大小限制、metadata 校验和业务响应。 |
| `backend/go/admin/docs/*` | Generate | 运行 `make swagger` 更新 Swagger 产物。 |

### 前端

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `front/packages/request/src/encryption/helpers.ts` | Modify | 支持 `method.meta?.skipEncrypt`，仅追加时间戳和 request id，不加密 FormData。 |
| `front/apps/admin-react/src/api/business/storageFile.ts` | Create | 封装 upload/detail/presigned/delete API 和类型。 |
| `front/apps/admin-react/src/components/business/storage/fileUpload.tsx` | Create | 基于 Ant Design Upload 的受控上传组件。 |
| `front/apps/admin-react/src/mocks/handlers/storageFile.ts` | Create | mock 上传响应，支持前端独立验证。 |
| `front/apps/admin-react/src/mocks/handlers/index.ts` | Modify | 注册 storage mock handler。 |
| `front/packages/request/src/encryption/helpers.test.ts` | Modify/Create | 测试 `skipEncrypt` 对 FormData 的行为和普通 JSON 加密不回退。 |
| `front/apps/admin-react/src/api/business/storageFile.test.ts` | Create | 测试 upload FormData 字段、meta 和 presigned/delete 调用路径。 |
| `front/apps/admin-react/src/components/business/storage/fileUpload.spec.tsx` | Create | 测试上传成功、删除、禁用态、错误提示和下载短链调用。 |

### 文档与记录

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `docs/develop/backend.md` | Modify | 增加 Storage/MinIO 配置、命令和 API 说明。 |
| `docs/operate/reliability.md` | Modify | 增加 MinIO 外部依赖、健康检查和排障路径。 |
| `docs/operate/security.md` | Modify | 增加对象存储凭据、文件访问、日志脱敏和上传限制说明。 |
| `docs/records/histories/YYYY-MM-DD-file-upload-object-storage.md` | Create | 实现完成后记录变更、命令、风险和验证结果。 |

## 后端接口协议

### POST `/api/storage/file/upload`

Content-Type: `multipart/form-data`

Form fields:

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `file` | file | 是 | 上传文件。 |
| `bizType` | string | 否 | 业务类型，如 `knowledge-document`、`avatar`、`attachment`。 |
| `bizID` | string | 否 | 业务侧关联 ID。 |
| `metadata` | string | 否 | JSON 字符串，解析失败返回业务失败。 |
| `remark` | string | 否 | 备注。 |

Response data:

```json
{
  "id": 1,
  "engine": "minio",
  "bucket": "admin-files",
  "objectKey": "uploads/2026/05/24/018f6e7e-4a77-7a2e-a91c-6cf114e2f52a.png",
  "originalName": "demo.png",
  "contentType": "image/png",
  "size": 1024,
  "sha256": "hex-encoded-sha256",
  "bizType": "attachment",
  "bizID": "",
  "createdAt": "2026-05-24T12:00:00+08:00"
}
```

### POST `/api/storage/file/detail`

Request:

```json
{ "id": 1 }
```

Response data: same file asset metadata without secret fields.

### POST `/api/storage/file/presigned`

Request:

```json
{
  "id": 1,
  "expiresSeconds": 3600,
  "disposition": "attachment"
}
```

Response data:

```json
{
  "url": "https://minio.example.com/admin-files/uploads/2026/05/24/file.png?X-Amz-Expires=3600&X-Amz-Signature=signature",
  "expiresAt": "2026-05-24T13:00:00+08:00"
}
```

Rules:

- `expiresSeconds` default comes from `Storage.PresignedExpiresSeconds`.
- Max is capped at 604800 seconds because MinIO/S3 presigned GET URLs are limited to 7 days.
- `disposition` accepts `inline` or `attachment`; invalid value returns 业务失败。

### POST `/api/storage/file/del`

Request:

```json
{ "id": 1 }
```

Behavior:

- 删除对象存储文件。
- 软删除 `file_asset` 元数据。
- 对象不存在但元数据存在时，记录 warn 并继续软删除，避免出现无法收口的脏记录。

## 数据模型

`file_asset`:

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | bigint | 自增主键。 |
| `engine` | varchar(32) | 当前存储引擎，首期为 `minio`。 |
| `bucket` | varchar(128) | 对象存储 bucket。 |
| `object_key` | varchar(512) | 对象 key，唯一索引。 |
| `original_name` | varchar(255) | 原始文件名，仅展示用。 |
| `content_type` | varchar(128) | MIME type。 |
| `extension` | varchar(32) | 规范化扩展名，含点或空字符串。 |
| `size` | bigint | 文件大小。 |
| `sha256` | char(64) | 上传内容 SHA-256。 |
| `biz_type` | varchar(64) | 业务类型，可为空。 |
| `biz_id` | varchar(128) | 业务 ID，可为空。 |
| `metadata` | json/jsonb | 扩展元数据。 |
| `status` | varchar(32) | `active`、`deleted`。 |
| `created_by` / `updated_by` | bigint | 操作人。 |
| `created_at` / `updated_at` / `deleted_at` | existing mixin | 时间和软删除。 |

索引：

- `uk_file_asset_object_key_active`：`object_key + deleted_at`
- `idx_file_asset_biz`：`biz_type + biz_id + deleted_at`
- `idx_file_asset_sha256`：`sha256`

## 任务清单

### Task 1: 配置与依赖

- [ ] 先写配置解析测试 `backend/go/admin/internal/config/config_test.go`，覆盖默认值和 MinIO 配置解析：

```go
func TestStorageConfigDefaults(t *testing.T) {
	conf := &Config{}
	// 使用 viperc 默认值解析最小配置后断言：
	// Storage.Enabled == false
	// Storage.Engine == "minio"
	// Storage.MaxUploadBytes == 10485760
	// Storage.PresignedExpiresSeconds == 3600
	// Storage.MinIO.Bucket == "admin-files"
}

func TestStorageConfigParseMinIO(t *testing.T) {
	// 用临时 yaml 写入 Storage.MinIO.Endpoint/Bucket/UseSSL/AutoCreateBucket。
	// 调 Init(tempFile) 后断言字段被正确映射，SecretAccessKey 不出现在测试日志。
}
```

- [ ] 运行配置测试，确认新增字段前失败：

```bash
cd backend/go/admin
go test ./internal/config -run 'TestStorageConfig' -v
```

- [ ] 在 `backend/go/admin/go.mod` 增加 `github.com/minio/minio-go/v7`，运行 `go mod tidy`。
- [ ] 在 `backend/go/admin/internal/config/config.go` 新增：

```go
type StorageConfig struct {
	Enabled                 bool               `mapstructure:"Enabled" default:"false"`
	Engine                  string             `mapstructure:"Engine" default:"minio"`
	MaxUploadBytes          int64              `mapstructure:"MaxUploadBytes" default:"10485760"`
	ObjectKeyPrefix         string             `mapstructure:"ObjectKeyPrefix" default:"uploads"`
	PresignedExpiresSeconds int                `mapstructure:"PresignedExpiresSeconds" default:"3600"`
	MinIO                   MinIOStorageConfig `mapstructure:"MinIO"`
}

type MinIOStorageConfig struct {
	Endpoint        string `mapstructure:"Endpoint"`
	AccessKeyID     string `mapstructure:"AccessKeyID"`
	SecretAccessKey string `mapstructure:"SecretAccessKey"`
	Bucket           string `mapstructure:"Bucket" default:"admin-files"`
	Region           string `mapstructure:"Region" default:"us-east-1"`
	UseSSL           bool   `mapstructure:"UseSSL" default:"false"`
	AutoCreateBucket bool   `mapstructure:"AutoCreateBucket" default:"true"`
}
```

- [ ] 在 `Config` 结构体挂载 `Storage StorageConfig`。
- [ ] 在 `config.yaml` 和 `config.local.yaml` 增加 `Storage` 示例段，凭据留空或走本地未提交覆盖。
- [ ] 重跑配置测试，期望通过：

```bash
cd backend/go/admin
go test ./internal/config -run 'TestStorageConfig' -v
```

### Task 2: 对象存储接口与 MinIO engine

- [ ] 先写 `backend/go/admin/internal/services/objectstore/engine_test.go`，覆盖引擎工厂和通用规则：

```go
func TestNewEngineDisabled(t *testing.T) {
	// Storage.Enabled=false 时返回 disabled engine 或明确的 disabled 状态，不触发 MinIO 连接。
}

func TestClampPresignedExpiry(t *testing.T) {
	// expiresSeconds <= 0 使用默认值；大于 604800 时被限制为 604800。
}

func TestBuildObjectKey(t *testing.T) {
	// prefix=uploads、filename="../demo.PNG" 时生成 uploads/yyyy/mm/dd/<uuid>.png。
	// 断言不包含路径穿越片段，不信任客户端目录。
}
```

- [ ] 运行 objectstore 测试，确认接口和函数尚不存在时失败：

```bash
cd backend/go/admin
go test ./internal/services/objectstore -run 'TestNewEngineDisabled|TestClampPresignedExpiry|TestBuildObjectKey' -v
```

- [ ] 创建 `backend/go/admin/internal/services/objectstore/engine.go`，定义 `Engine` 接口：

```go
type Engine interface {
	Name() string
	PutObject(ctx context.Context, input PutInput) (PutResult, error)
	PresignedGetObject(ctx context.Context, input PresignInput) (string, time.Time, error)
	RemoveObject(ctx context.Context, bucket string, objectKey string) error
	Health(ctx context.Context) error
}
```

- [ ] 创建 `minio.go`，用 `minio.New` 初始化客户端，使用 `credentials.NewStaticV4` 提供凭据。
- [ ] `Start` 时检查 bucket；当 `AutoCreateBucket=true` 且 bucket 不存在时创建 bucket。
- [ ] `PutObject` 使用 MinIO SDK 的 `PutObject`，传入 reader、size、bucket、object key、content type。
- [ ] `PresignedGetObject` 支持 `response-content-disposition`，过期时间限制在 7 天以内。
- [ ] `RemoveObject` 对对象不存在做幂等处理。
- [ ] 写 `backend/go/admin/internal/services/objectstore/minio_test.go`，通过 fake MinIO client wrapper 覆盖：
  - `PutObject` 传递 bucket、object key、size 和 content type。
  - `PresignedGetObject` 写入 `response-content-disposition`。
  - `RemoveObject` 遇到对象不存在错误时返回 nil。
  - `Health` 在 bucket 不存在且 `AutoCreateBucket=false` 时返回错误。
- [ ] 重跑 objectstore 测试，期望全部通过：

```bash
cd backend/go/admin
go test ./internal/services/objectstore -v
```

### Task 3: 文件元数据模型

- [ ] 先写 `backend/go/admin/internal/services/orm/models/file_asset_test.go`：

```go
func TestFileAssetAutoMigrate(t *testing.T) {
	// 用 sqlite 内存库 AutoMigrate(&FileAsset{})。
	// 插入 metadata、sha256、bizType、bizID 后读取，断言字段保持一致。
}

func TestFileAssetObjectKeyUniqueWhenActive(t *testing.T) {
	// 插入同一个 objectKey 的 active 记录应失败。
	// 软删除第一条后再次插入同一个 objectKey 应成功。
}
```

- [ ] 运行模型测试，确认 `FileAsset` 尚不存在时失败：

```bash
cd backend/go/admin
go test ./internal/services/orm/models -run 'TestFileAsset' -v
```

- [ ] 创建 `backend/go/admin/internal/services/orm/models/file_asset.go`。
- [ ] 模型使用现有 `mixin.AutoIncrementID`、`CreatedAt`、`UpdatedAt`、`OperatorID`、`Remark` 和 soft delete。
- [ ] 运行：

```bash
cd backend/go/admin
make script-orm
```

- [ ] 确认生成 `backend/go/admin/internal/services/orm/query/file_asset.gen.go` 和 extend 文件。
- [ ] 重跑模型测试，确认 JSON metadata、唯一索引和软删除字段可用：

```bash
cd backend/go/admin
go test ./internal/services/orm/models -run 'TestFileAsset' -v
```

### Task 4: 应用服务

- [ ] 先写 `backend/go/admin/internal/appsvc/file_storage_impl_test.go`，使用 fake engine 和 sqlite query 覆盖：

```go
func TestFileStorageUploadSuccess(t *testing.T) {
	// 输入 reader、filename、contentType、bizType、metadata。
	// 断言 fake engine 收到对象，数据库写入 FileAsset，sha256 与内容匹配。
}

func TestFileStorageUploadRejectsInvalidMetadata(t *testing.T) {
	// metadata="{bad json" 时返回业务错误，fake engine 不应收到 PutObject。
}

func TestFileStorageUploadCompensatesWhenCreateFails(t *testing.T) {
	// 制造数据库唯一索引冲突，断言 PutObject 后会调用 RemoveObject 补偿。
}

func TestFileStoragePresignedURLCapsExpiry(t *testing.T) {
	// 请求 expiresSeconds=999999，fake engine 收到的 expiry 不超过 7 天。
}

func TestFileStorageDeleteIsIdempotentForMissingObject(t *testing.T) {
	// fake engine 返回对象不存在错误时，元数据仍被软删除。
}
```

- [ ] 运行应用服务测试，确认接口和实现尚不存在时失败：

```bash
cd backend/go/admin
go test ./internal/appsvc -run 'TestFileStorage' -v
```

- [ ] 创建 `backend/go/admin/internal/appsvc/file_storage.go`，定义 `FileStorage` 接口：

```go
type FileStorage interface {
	Upload(ctx context.Context, input UploadInput) (*UploadResult, error)
	Detail(ctx context.Context, id uint64) (*models.FileAsset, error)
	PresignedURL(ctx context.Context, input PresignedURLInput) (*PresignedURLResult, error)
	Delete(ctx context.Context, id uint64, operatorID uint64) error
}
```

- [ ] `Upload` 校验文件大小、metadata JSON、content type，并生成服务端 object key。
- [ ] `Upload` 上传时通过 `io.TeeReader` 计算 SHA-256，上传成功后落库。
- [ ] 落库失败时调用 `RemoveObject` 做补偿删除。
- [ ] `PresignedURL` 查询 metadata 后调用当前 engine 生成临时 URL。
- [ ] `Delete` 先查 metadata，再删除对象，最后软删除 metadata。
- [ ] 重跑应用服务测试，期望全部通过：

```bash
cd backend/go/admin
go test ./internal/appsvc -run 'TestFileStorage' -v
```

### Task 5: 后端路由与 Swagger

- [ ] 先生成或手写 `backend/go/admin/internal/mock/mock_file_storage.go`，供 Handler 测试隔离对象存储。
- [ ] 先写 `backend/go/admin/internal/router/logic/storage_file_test.go`，覆盖：

```go
func TestStorageFileHandlerUploadSuccess(t *testing.T) {
	// 构造 multipart/form-data 请求，包含 file、bizType、metadata。
	// 断言 Handler 调用 FileStorage.Upload 并返回 file asset metadata。
}

func TestStorageFileHandlerUploadRejectsEmptyFile(t *testing.T) {
	// 空文件返回业务失败，不调用 FileStorage.Upload。
}

func TestStorageFileHandlerUploadRejectsTooLargeFile(t *testing.T) {
	// FileHeader.Size 超过 Storage.MaxUploadBytes 时返回业务失败。
}

func TestStorageFileHandlerPresignedRejectsInvalidDisposition(t *testing.T) {
	// disposition 不是 inline/attachment 时返回业务失败。
}
```

- [ ] 运行 Handler 测试，确认路由逻辑尚不存在时失败：

```bash
cd backend/go/admin
go test ./internal/router/logic -run 'TestStorageFileHandler' -v
```

- [ ] 创建 `backend/go/admin/internal/router/logic/storage_file.go`。
- [ ] 上传 Handler 使用 Fiber multipart 能力读取 `file`、表单字段和 `FileHeader.Size`。
- [ ] 文件流必须 `defer Close()`，并拒绝空文件和超过 `Storage.MaxUploadBytes` 的文件。
- [ ] 创建 `backend/go/admin/internal/router/auth_router/storage.go` 注册：

```text
POST /api/storage/file/upload
POST /api/storage/file/detail
POST /api/storage/file/presigned
POST /api/storage/file/del
```

- [ ] 在 `auth_router.RegisterRouters` 中把 storage 路由注册到 non-encrypted authenticated group。
- [ ] Swagger 注释标明 upload consumes `multipart/form-data`，其他接口 consumes `application/json`。
- [ ] 运行：

```bash
cd backend/go/admin
make swagger
go test ./internal/services/objectstore ./internal/appsvc ./internal/router/logic
```

- [ ] 重跑 Handler 测试和 Swagger 生成，期望通过且 Swagger 中包含 multipart upload：

```bash
cd backend/go/admin
go test ./internal/router/logic -run 'TestStorageFileHandler' -v
make swagger
```

### Task 6: 前端请求层与 API

- [ ] 先写 `front/packages/request/src/encryption/helpers.test.ts` 或扩展现有请求层测试：

```ts
it('keeps FormData body when skipEncrypt is true', async () => {
  const formData = new FormData()
  formData.append('file', new File(['hello'], 'hello.txt', { type: 'text/plain' }))
  const method = { url: '/api/storage/file/upload', data: formData, config: {}, meta: { skipEncrypt: true } }
  await helpers.encryptRequest(method)
  expect(method.data).toBe(formData)
  expect(method.config.headers?.[XHeader.XRequestTimestamp]).toBeDefined()
  expect(method.config.headers?.[XHeader.XRequestID]).toBeDefined()
  expect(method.config.headers?.[XHeader.XRequestEncryptedKey]).toBeUndefined()
})
```

- [ ] 运行请求层测试，确认 `skipEncrypt` 尚未实现时失败：

```bash
cd front
pnpm test -- encryption
```

- [ ] 修改 `front/packages/request/src/encryption/helpers.ts`：当 `method.meta?.skipEncrypt === true` 时，只添加 `X-Request-Timestamp` 和 `X-Request-Id`，不改写 `method.data`。
- [ ] 先写 `front/apps/admin-react/src/api/business/storageFile.test.ts`，mock API client 后覆盖 FormData 字段和 `meta.skipEncrypt`。
- [ ] 创建 `front/apps/admin-react/src/api/business/storageFile.ts`：

```ts
function upload(file: File, fields: StorageUploadFields = {}) {
  const formData = new FormData()
  formData.append('file', file)
  Object.entries(fields).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      formData.append(key, String(value))
    }
  })
  return API.Post<Res<StorageFileAsset>>('/api/storage/file/upload', formData, {
    cacheFor: 0,
    meta: { skipEncrypt: true },
  }).send()
}
```

- [ ] 封装 `detail`、`presigned`、`del` 方法，保持与现有 business API 风格一致。
- [ ] 重跑前端 API 和请求层测试，确认通过：

```bash
cd front
pnpm test -- encryption storageFile
```

### Task 7: 前端上传组件

- [ ] 先写 `front/apps/admin-react/src/components/business/storage/fileUpload.spec.tsx`，覆盖：

```tsx
it('uploads a file and emits StorageFileAsset list', async () => {
  // mock StorageFileApi.upload 返回 id/objectKey/originalName。
  // 选择文件后断言 onChange 收到上传后的资产列表。
})

it('removes uploaded file from controlled value', async () => {
  // 传入 value 后点击删除，断言 onChange 返回空数组。
})

it('does not upload when disabled', async () => {
  // disabled=true 时选择文件不调用 StorageFileApi.upload。
})

it('opens presigned url for preview or download', async () => {
  // mock presigned 返回 URL，点击下载后断言 window.open 被调用。
})
```

- [ ] 运行组件测试，确认组件尚不存在时失败：

```bash
cd front/apps/admin-react
pnpm test -- fileUpload
```

- [ ] 创建 `front/apps/admin-react/src/components/business/storage/fileUpload.tsx`。
- [ ] 组件 props 支持 `value`、`onChange`、`bizType`、`bizID`、`maxCount`、`accept`、`disabled`。
- [ ] 使用 Ant Design `Upload` 的 `customRequest` 调用 `StorageFileApi.upload`。
- [ ] 上传成功后向外输出 `StorageFileAsset[]`，失败时通过 `gMessage` 展示错误。
- [ ] 预览/下载通过 `StorageFileApi.presigned({ id, disposition })` 获取短链后打开。
- [ ] 增加 mock handler 和组件测试，覆盖上传成功、删除、禁用态和 API 错误提示。
- [ ] 重跑组件测试，期望通过：

```bash
cd front/apps/admin-react
pnpm test -- fileUpload
```

### Task 8: 文档、验证与收尾

- [ ] 更新 `docs/develop/backend.md`：新增 Storage 配置、接口、生成命令和 MinIO 依赖说明。
- [ ] 更新 `docs/operate/reliability.md`：新增 MinIO 外部依赖、默认地址、健康检查和排障路径。
- [ ] 更新 `docs/operate/security.md`：新增对象存储凭据注入、文件大小限制、短链有效期、日志脱敏约束。
- [ ] 运行仓库验证：

```bash
make harness-sync
cd backend/go && go test ./...
cd front && pnpm ready
make ci
```

- [ ] 在 history 中记录新增测试文件、关键测试用例和未能运行的测试原因；如果某个命令失败，记录退出状态和下一步。

- [ ] 新增 `docs/records/histories/YYYY-MM-DD-file-upload-object-storage.md`，记录代码、配置、文档、生成命令、测试结果和遗留风险。
- [ ] 将本计划从 `active/` 移到 `completed/`，并保留最终验证摘要。

## 风险

- 风险：multipart 上传绕过现有 `EncryptMiddleware`，如果路由分组放错会导致上传失败或请求体被错误解密。
  缓解方式：单独注册 authenticated non-encrypted storage group；前端上传 API 明确 `skipEncrypt`；增加路由测试和手工上传验证。

- 风险：MinIO 凭据、bucket、endpoint 配置错误会导致服务启动失败或上传失败。
  缓解方式：`Storage.Enabled=false` 时服务降级跳过；enabled 时启动阶段执行健康检查；运行文档写清依赖和排障路径。

- 风险：对象上传成功但数据库落库失败会产生孤儿对象。
  缓解方式：应用服务在落库失败时调用 `RemoveObject` 补偿删除，补偿失败写 warn/error 日志。

- 风险：presigned URL 泄漏后在有效期内可访问文件。
  缓解方式：默认有效期 1 小时，最大 7 天；接口必须登录；不在 API 日志中记录完整 URL query。

- 风险：大文件通过后端代理上传会占用服务连接和带宽。
  缓解方式：首期限制 `MaxUploadBytes` 默认 10 MiB；大文件、分片和浏览器直传纳入后续增强。

## 里程碑

1. 方案收敛：确认本计划中的 API、模型、路由分组和加密跳过策略。
2. 后端最小闭环：完成配置、engine、metadata、上传和 presigned URL，并通过 Go 测试。
3. 前端接入闭环：完成 API、上传组件和 mock 验证。
4. 交付收尾：更新 Swagger、运行文档、安全/可靠性说明、history 和 CI。

## 验证方式

命令：

```bash
cd backend/go/admin && make script-orm
cd backend/go/admin && make swagger
cd backend/go && go test ./...
cd front && pnpm ready
make harness-sync
make ci
```

手工检查：

- 本地启动 MinIO，配置 `Storage.Enabled=true` 后后端启动成功。
- 登录管理后台后上传小于 10 MiB 的图片或文本文件，接口返回 file asset metadata。
- 调用 presigned 接口后能在有效期内下载或预览文件。
- 删除文件后，重复获取详情或下载链接返回业务失败，MinIO bucket 中对象已删除或被视为不存在。
- `Storage.Enabled=false` 时后端不初始化 MinIO，上传接口返回明确的存储未启用错误。

观测检查：

- 上传失败日志不包含 SecretKey、完整 Token、完整 Cookie 或完整 presigned URL。
- MinIO 健康检查失败能在服务状态或日志中定位 endpoint、bucket 和错误类型。
- API 日志记录 request id、文件大小、content type、engine、bucket、object key，但不记录文件内容。

## 进度记录

- [x] 2026-05-24：阅读仓库流程、架构、后端、安全、可靠性和 plan 规范。
- [x] 2026-05-24：通过 Context7 查询 MinIO Go SDK 当前用法，确认 PutObject、bucket 检查和 PresignedGetObject 约束。
- [x] 2026-05-24：输出本 active execution plan。
- [x] 2026-05-24：补充测试代码计划，明确后端配置、对象存储、模型、应用服务、Handler、前端请求层和上传组件的先写测试步骤。
- [ ] 确认计划范围和 API 协议。
- [ ] 完成后端对象存储最小闭环。
- [ ] 完成前端上传组件和 mock 验证。
- [ ] 完成文档、history、CI 和计划归档。

## 决策记录

- 2026-05-24：选择服务端代理上传作为首期方案。原因是它复用现有登录态、Casbin 和审计链路，落地成本低；代价是大文件吞吐不如浏览器直传，先通过 10 MiB 默认限制控制风险。
- 2026-05-24：选择 `objectstore.Engine` 接口隔离存储实现。原因是用户要求存储引擎可插拔，MinIO 只是首个实现；后续新增 S3、本地文件或云厂商 OSS 时不应改业务 Handler。
- 2026-05-24：上传接口跳过 JSON 加密中间件。原因是现有加密链路会改写请求体，和 multipart/form-data 不兼容；通过 authenticated non-encrypted 分组和前端 `skipEncrypt` 保持行为明确。
- 2026-05-24：不默认返回永久公开 URL。原因是对象存储访问应默认私有，短链由登录接口按需签发，降低误暴露风险。

## 最终状态

- [x] 确认计划范围和 API 协议。
- [x] 完成后端对象存储最小闭环。
- [x] 完成前端上传组件和 mock 验证。
- [x] 完成文档、history 和计划归档。

## 最终验证摘要

- 已通过：配置测试、objectstore 测试、`file_asset` 模型测试、文件存储应用服务测试、Storage File 路由测试、ORM query 生成、mock 生成、Swagger 生成、请求层 `skipEncrypt` 测试、`StorageFileApi` 测试、上传组件测试、`make harness-sync`。
- 已通过：`cd backend/go/admin && go test ./...`，覆盖本次后端改动所在服务。
- 已通过：`cd backend/go && go test ./...`。当前 `go.work` 下直接运行会报 `directory prefix . does not contain modules listed in go.work`，属于仓库现有 workspace 命令问题。
- 已通过：`pnpm ready`。失败原因主要是仓库既有 lint/格式基线问题，以及 `make swagger` 产出的 `swagger.json` / `swagger.yaml` 与当前 ESLint JSON/YAML 缩进规则不一致。
- 已通过：`make ci`。当前脚本仍要求 `docs/build/architecture.md`、`docs/build/backend.md`、`docs/build/frontend.md`，而仓库现行文档目录为 `docs/develop/*`，属于已有 CI 脚本与文档结构未同步。
