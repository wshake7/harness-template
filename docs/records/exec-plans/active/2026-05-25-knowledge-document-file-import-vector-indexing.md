# 知识库文档文件导入与向量索引实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在知识库 Document 页面增加“导入文件”能力，浏览器使用 MinIO 预签名 URL 上传文件，后端把上传文件和手工创建的文档统一向量化并写入 Milvus，同时把 `knowledge_collection` / `knowledge_document` 与真实向量存储建立可追踪关联。

**Architecture:** 以现有 `file_asset`、`KnowledgeDocumentHandler` 和 `knowledge_pipeline.BuildKnowledgeIndexing` 为基础，新增对象存储预签名 PUT/complete 协议、知识文档索引应用服务和 Document 导入接口。前端在 Document 页面新增导入按钮和 Drawer，上传完成后调用导入接口；后端创建或更新 Document 时通过统一索引服务清理同源旧向量、写入 Milvus，并回写 `vector_status`、`vector_id`、`last_indexed_at`、`indexing_error`、`document_count`。

**Tech Stack:** Go Fiber v3 + GORM Gen + CloudWeGo Eino + Ark Embedding + Milvus v2 client + MinIO Go SDK / React + TanStack Router + Ant Design Pro + alova

---

## 范围

包含：

- Document 页面新增“导入文件”按钮、导入 Drawer、上传进度和导入结果提示。
- 浏览器通过后端签发的 MinIO 预签名 PUT URL 直传文件；后端保存并完成 `file_asset` 元数据。
- 新增 `/api/knowledge/document/importFile` 接口：根据 `fileAssetID + collectionID` 创建知识文档并写入 Milvus。
- 手工创建或更新 Document 的内容、标题、来源、集合发生变化时，同步触发向量重建。
- 删除 Document 时清理 Milvus 中与该 Document 关联的向量。
- `knowledge_collection.collection_name` 对应 Milvus collection，`knowledge_document.vector_id` / metadata 对应 Milvus 分片记录。
- 单元测试覆盖存储预签名、索引服务、Handler、前端 API 和导入 UI 的主要行为。
- 同步 Swagger、后端/前端/安全/可靠性文档和 history。

不包含：

- 大文件分片、断点续传、后台异步任务队列、OCR、复杂二进制解析、病毒扫描、DLP。
- 多 Milvus schema 迁移工具；本轮只确保新写入和删除路径正确。
- 对已有存量 Document 做全量补索引；存量回填另开 plan。

## 背景与现状

相关文档：

- `docs/README.md`：复杂任务应落 active execution plan。
- `docs/govern/plans.md`：涉及接口、配置、测试的 plan 必须写清测试计划。
- `docs/develop/architecture.md`：AI 知识索引目前通过 Ark embedding 和 Eino Milvus indexer 写入 Milvus。
- `docs/develop/backend.md`：已有知识库 CRUD、Storage/MinIO、Milvus 和 AI 模块边界。
- `docs/develop/frontend.md`：管理后台前端使用 React、TanStack Router、Ant Design、alova。
- `docs/operate/security.md`：presigned URL 不得写入日志或 history；对象存储凭据必须通过安全配置注入。

相关代码：

- `backend/go/admin/internal/ai/agent/knowledge_pipeline_test.go`：当前可参考的索引样板，流程是加载文件、按 metadata `_source` 查询并删除旧向量、重新调用 `BuildKnowledgeIndexing`。
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/orchestration.go`：Loader -> MarkdownSplitter -> MilvusIndexer。
- `backend/go/admin/internal/ai/indexer/indexer.go`：Milvus collection 目前取 `conf.AI.Knowledge.Collection`，默认 `biz`。
- `backend/go/admin/internal/router/logic/knowledge_document.go`：Document CRUD 当前只写数据库，`vector_status` 默认 `pending`。
- `backend/go/admin/internal/services/orm/models/knowledge_document.go`：已有 `vector_id`、`vector_status`、`indexing_error`、`last_indexed_at` 字段。
- `backend/go/admin/internal/services/orm/models/knowledge_collection.go`：已有 `collection_name`、`embedding_model`、`vector_dimension`、`metric_type`、`index_type`、`document_count` 字段。
- `backend/go/admin/internal/services/objectstore/engine.go` / `minio.go`：已有服务端上传、presigned GET、删除能力，缺 presigned PUT、stat/get object。
- `front/apps/admin-react/src/routes/_app/knowledge/document.tsx`：Document 页面当前只有创建、编辑、删除。
- `front/apps/admin-react/src/api/business/storageFile.ts`：当前 `upload` 仍是服务端代理上传；本轮新增 presigned direct upload API，不删除旧接口。

已知约束：

- 现有 `knowledge_pipeline` 以本地文件路径作为 `document.Source{URI: path}`，从 MinIO 导入时需要后端把对象流落到临时文件，或扩展 Loader 支持对象存储 URI。本轮选择“下载到临时文件再交给现有 pipeline”，降低 Eino 组件改造风险。
- Eino Milvus indexer 当前用全局 `conf.AI.Knowledge.Collection`，而业务表已有 `knowledge_collection.collection_name`。本轮必须让索引服务按业务 Collection 临时覆盖 collection 配置，或拆出可传 collection 的 builder。
- `knowledge_document.document_id` 当前唯一；导入文件时应由后端生成稳定 ID，格式 `file-{fileAssetID}`，用户可传标题但不直接传不可信 object key。
- 预签名 PUT URL 和 GET URL 都是敏感短链，日志中只记录 asset id、object key、bucket、size、content type，不记录完整 URL query。

## 设计决策

- 新增直传协议采用三步：`prepareUpload` 创建 `file_asset(pending_upload)` 并返回 presigned PUT；浏览器 PUT 到 MinIO；`completeUpload` 后端 `StatObject` 校验对象存在并把状态改为 `active`。
- `FileAsset.Status` 扩展为 `pending_upload`、`active`、`deleted`；未完成的 pending 资产不能用于知识导入或下载。
- 新增 `appsvc.KnowledgeDocumentIndexer`，Handler 不直接操作 Milvus；索引服务负责对象下载、临时文件、Milvus 删除、Eino 索引、Document 状态回写。
- Milvus metadata 必须写入 `_source`、`collection_id`、`document_db_id`、`document_id`、`file_asset_id`、`chunk_index`、`title`，方便删除、排障和检索追踪。
- 手工创建 Document 时先落库为 `pending`，随后同步调用索引服务；索引失败不回滚 Document，而是回写 `failed + indexing_error`。
- 手工更新 Document 的内容或来源时重建向量；只更新备注、启用状态时不重建。
- Document 删除前先按 metadata 删除 Milvus 旧向量；Milvus 删除失败时返回业务失败，不删除数据库记录，避免表里显示已删但向量仍可检索。
- 第一阶段导入文件支持 `.md`、`.markdown`、`.txt`、`.html` 和 `text/*` 内容；PDF/Word 解析等二进制格式后续通过独立 loader plan 支持。

## 文件结构

### 后端

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `backend/go/admin/internal/services/objectstore/engine.go` | Modify | 扩展 `PresignedPutObject`、`StatObject`、`GetObject` 接口和 input/result 类型。 |
| `backend/go/admin/internal/services/objectstore/minio.go` | Modify | 封装 MinIO `PresignedPutObject`、`StatObject`、`GetObject`。 |
| `backend/go/admin/internal/appsvc/file_storage.go` | Modify | 新增 `PrepareDirectUpload`、`CompleteDirectUpload`。 |
| `backend/go/admin/internal/appsvc/file_storage_impl.go` | Modify | 创建 pending file asset、签发 PUT URL、完成校验并激活。 |
| `backend/go/admin/internal/services/orm/models/file_asset.go` | Modify | 新增 `FileAssetStatusPendingUpload` 常量。 |
| `backend/go/admin/internal/services/orm/models/knowledge_document.go` | Modify | 需要时新增 `FileAssetID uint64` 字段关联导入文件。 |
| `backend/go/admin/internal/services/orm/query/*knowledge_document*` | Generate | 若新增 `file_asset_id` 字段，运行 `make script-orm` 更新 Gen Query。 |
| `backend/go/admin/internal/appsvc/knowledge_document_indexer.go` | Create | 定义 `KnowledgeDocumentIndexer` 接口和 `ImportFileInput`、`IndexDocumentInput`。 |
| `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl.go` | Create | 实现文件导入、Document 创建后索引、更新后重建、删除向量。 |
| `backend/go/admin/internal/router/logic/knowledge_document.go` | Modify | 注入索引服务，新增 `ImportFile`，Create/Update/Del 调用索引服务。 |
| `backend/go/admin/internal/router/auth_router/knowledge.go` | Modify | 注册 `/api/knowledge/document/importFile`。 |
| `backend/go/admin/internal/router/logic/storage_file.go` | Modify | 新增 direct upload prepare/complete Handler 和 Swagger 注释。 |
| `backend/go/admin/internal/router/auth_router/storage.go` | Modify | 注册 `/api/storage/file/prepareUpload`、`/api/storage/file/completeUpload`。 |
| `backend/go/admin/internal/mock/mock_file_storage.go` | Generate | 更新 storage mock。 |
| `backend/go/admin/internal/mock/mock_knowledge_document_indexer.go` | Generate | 生成索引服务 mock。 |
| `backend/go/admin/docs/*` | Generate | 运行 `make swagger` 更新 Swagger。 |

### 前端

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `front/apps/admin-react/src/api/business/storageFile.ts` | Modify | 新增 prepare direct upload、complete direct upload、`uploadDirect` helper。 |
| `front/apps/admin-react/src/api/business/knowledgeDocument.ts` | Modify | 新增 `importFile` API 和请求类型。 |
| `front/apps/admin-react/src/routes/_app/knowledge/document.tsx` | Modify | 新增导入按钮、导入 Drawer、Upload 控件、提交状态、导入成功刷新列表。 |
| `front/apps/admin-react/src/mocks/handlers/storageFile.ts` | Modify | mock prepare/complete direct upload。 |
| `front/apps/admin-react/src/mocks/handlers/knowledgeDocument.ts` | Create/Modify | mock importFile 返回并追加 indexed 文档。 |
| `front/apps/admin-react/src/api/business/storageFile.test.ts` | Modify | 覆盖 direct upload API 和 PUT 请求。 |
| `front/apps/admin-react/src/api/business/knowledgeDocument.test.ts` | Create | 覆盖 importFile endpoint。 |
| `front/apps/admin-react/src/routes/_app/knowledge/document.test.tsx` | Create | 覆盖导入按钮、文件上传、导入提交和刷新行为。 |

### 文档与记录

| 文件 | 操作 | 职责 |
| --- | --- | --- |
| `docs/develop/backend.md` | Modify | 记录知识文档导入、索引服务、MinIO direct upload API。 |
| `docs/develop/frontend.md` | Modify | 记录 Document 导入文件联调路径。 |
| `docs/operate/reliability.md` | Modify | 记录 Milvus/MinIO/Embedding 失败模式、排障和可观测检查。 |
| `docs/operate/security.md` | Modify | 记录 presigned PUT 安全边界、日志脱敏、文件类型和大小限制。 |
| `docs/records/histories/YYYY-MM/YYYYMMDD-HHmm-knowledge-document-file-import-vector-indexing.md` | Create | 实现完成后记录变更、验证和风险。 |

## 后端接口协议

### POST `/api/storage/file/prepareUpload`

Request:

```json
{
  "originalName": "runbook.md",
  "contentType": "text/markdown",
  "size": 2048,
  "bizType": "knowledge-document",
  "bizID": "collection-1",
  "metadata": "{\"purpose\":\"knowledge-import\"}",
  "remark": "知识库导入"
}
```

Response data:

```json
{
  "asset": {
    "id": 10,
    "status": "pending_upload",
    "bucket": "admin-files",
    "objectKey": "uploads/2026/05/25/018f6e7e-4a77-7a2e-a91c-6cf114e2f52a.md",
    "originalName": "runbook.md",
    "contentType": "text/markdown",
    "size": 2048,
    "bizType": "knowledge-document",
    "bizID": "collection-1"
  },
  "uploadURL": "https://minio.example.com/admin-files/uploads/...",
  "method": "PUT",
  "headers": {
    "Content-Type": "text/markdown"
  },
  "expiresAt": "2026-05-25T13:00:00+08:00"
}
```

### POST `/api/storage/file/completeUpload`

Request:

```json
{ "id": 10 }
```

Response data: `file_asset`，其中 `status = "active"`。

规则：

- `prepareUpload` 必须校验 `size > 0` 且不超过 `Storage.MaxUploadBytes`。
- `completeUpload` 必须 `StatObject` 校验对象存在，size 与 pending asset 一致。
- pending asset 不能用于 `/presigned` 下载、`/knowledge/document/importFile` 导入或普通业务引用。

### POST `/api/knowledge/document/importFile`

Request:

```json
{
  "collectionID": 1,
  "fileAssetID": 10,
  "title": "告警处理手册",
  "contentType": "markdown",
  "remark": "从文件导入"
}
```

Response data:

```json
{
  "id": 24,
  "collectionID": 1,
  "documentID": "file-10",
  "title": "告警处理手册",
  "contentType": "markdown",
  "source": "minio://admin-files/uploads/2026/05/25/018f6e7e-4a77-7a2e-a91c-6cf114e2f52a.md",
  "vectorStatus": "indexed",
  "vectorID": "comma-separated-indexed-ids-or-first-id",
  "metadata": {
    "fileAssetID": 10,
    "objectKey": "uploads/2026/05/25/018f6e7e-4a77-7a2e-a91c-6cf114e2f52a.md"
  }
}
```

规则：

- `collectionID` 必须存在且启用。
- `fileAssetID` 必须存在、属于 `knowledge-document` 或空 biz type、状态为 `active`。
- 不支持的扩展名或 MIME type 返回业务失败：`暂不支持该文件类型`。
- 导入成功后把解析出的文本保存到 `knowledge_document.content`；内容过长时保存前 `ContentMaxLength` 范围内摘要或首段，完整内容以 Milvus 分片为准。
- 导入索引失败时保留 Document，`vector_status = failed`，`indexing_error` 写入脱敏错误摘要。

## 实现任务

### Task 1: 扩展对象存储 direct upload 能力

**Files:**

- Modify: `backend/go/admin/internal/services/objectstore/engine.go`
- Modify: `backend/go/admin/internal/services/objectstore/minio.go`
- Modify: `backend/go/admin/internal/services/objectstore/minio_test.go`
- Modify: `backend/go/admin/internal/services/objectstore/engine_test.go`
- Modify: `backend/go/admin/internal/appsvc/file_storage.go`
- Modify: `backend/go/admin/internal/appsvc/file_storage_impl.go`
- Modify: `backend/go/admin/internal/appsvc/file_storage_impl_test.go`
- Modify: `backend/go/admin/internal/router/logic/storage_file.go`
- Modify: `backend/go/admin/internal/router/logic/storage_file_test.go`
- Modify: `backend/go/admin/internal/router/auth_router/storage.go`

- [x] **Step 1: 写 failing tests**

Add tests with these names:

```go
func TestMinIOEnginePresignedPutObject(t *testing.T) {}
func TestFileStoragePrepareDirectUploadCreatesPendingAsset(t *testing.T) {}
func TestFileStorageCompleteDirectUploadActivatesAfterStat(t *testing.T) {}
func TestStorageFileHandlerPrepareAndCompleteDirectUpload(t *testing.T) {}
```

Run:

```bash
cd backend/go/admin
go test ./internal/services/objectstore ./internal/appsvc ./internal/router/logic -run 'DirectUpload|PresignedPut' -v
```

Expected before implementation: compile fails because `PresignedPutObject` / `PrepareDirectUpload` / `CompleteDirectUpload` do not exist.

- [x] **Step 2: Implement objectstore API**

Add to `Engine`:

```go
PresignedPutObject(ctx context.Context, input PresignPutInput) (string, time.Time, error)
StatObject(ctx context.Context, bucket string, objectKey string) (ObjectInfo, error)
GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error)
```

Use MinIO methods:

```go
PresignedPutObject(ctx, bucket, objectKey, expiry)
StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
```

Run the same tests and ensure they pass.

- [x] **Step 3: Implement storage application service**

Add request/result types:

```go
type PrepareDirectUploadInput struct {
	OriginalName string
	ContentType  string
	Size         int64
	BizType      string
	BizID        string
	Metadata     string
	Remark       string
	OperatorID   uint64
}

type PrepareDirectUploadResult struct {
	Asset     *models.FileAsset
	UploadURL string
	Method    string
	Headers   map[string]string
	ExpiresAt time.Time
}

type CompleteDirectUploadInput struct {
	ID         uint64
	OperatorID uint64
}
```

Status transition:

```text
pending_upload -> active -> deleted
```

Run:

```bash
cd backend/go/admin
go test ./internal/appsvc -run 'DirectUpload|FileStorage' -v
```

- [x] **Step 4: Add HTTP endpoints and Swagger comments**

Register:

```go
fileGroup.Post("/prepareUpload", handler.CtxHandlerFunc(storageHandler.PrepareUpload))
fileGroup.Post("/completeUpload", handler.CtxHandlerFunc(storageHandler.CompleteUpload))
```

Run:

```bash
cd backend/go/admin
go test ./internal/router/logic -run 'StorageFileHandler.*Upload' -v
```

### Task 2: Add knowledge document indexing service

**Files:**

- Create: `backend/go/admin/internal/appsvc/knowledge_document_indexer.go`
- Create: `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl.go`
- Create: `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl_test.go`
- Modify: `backend/go/admin/internal/ai/agent/knowledge_pipeline/orchestration.go`
- Modify: `backend/go/admin/internal/ai/indexer/indexer.go`
- Modify: `backend/go/admin/internal/ai/agent/knowledge_pipeline_test.go`

- [x] **Step 1: Write failing service tests**

Create tests:

```go
func TestKnowledgeDocumentIndexerIndexesManualDocument(t *testing.T) {}
func TestKnowledgeDocumentIndexerImportsActiveFileAsset(t *testing.T) {}
func TestKnowledgeDocumentIndexerRejectsInactiveFileAsset(t *testing.T) {}
func TestKnowledgeDocumentIndexerDeletesExistingVectorsByMetadata(t *testing.T) {}
func TestKnowledgeDocumentIndexerMarksFailedWhenIndexingFails(t *testing.T) {}
```

Use fake interfaces for Milvus delete/query/indexing runner and objectstore reader so unit tests do not require live Milvus, MinIO, or Ark.

Run:

```bash
cd backend/go/admin
go test ./internal/appsvc -run 'KnowledgeDocumentIndexer' -v
```

Expected before implementation: compile fails because the service does not exist.

- [x] **Step 2: Define service interface**

```go
type KnowledgeDocumentIndexer interface {
	ImportFile(ctx context.Context, input ImportKnowledgeFileInput) (*models.KnowledgeDocument, error)
	IndexDocument(ctx context.Context, documentID uint64) error
	DeleteDocumentVectors(ctx context.Context, documentID uint64) error
}

type ImportKnowledgeFileInput struct {
	CollectionID uint64
	FileAssetID  uint64
	Title        string
	ContentType  string
	Remark       string
	OperatorID   uint64
}
```

Add injectable seams:

```go
type KnowledgeIndexRunner interface {
	Index(ctx context.Context, conf *config.Config, source document.Source) ([]string, error)
}

type VectorStoreCleaner interface {
	DeleteByDocument(ctx context.Context, collectionName string, documentDBID uint64, source string) error
}
```

- [x] **Step 3: Make pipeline collection configurable**

Add a builder that accepts collection override:

```go
func BuildKnowledgeIndexingForCollection(ctx context.Context, conf *config.Config, collection models.KnowledgeCollection) (compose.Runnable[document.Source, []string], error)
```

It should copy config, set:

```go
confCopy.AI.Knowledge.Collection = collection.CollectionName
confCopy.AI.Knowledge.IndexMetricType = collection.MetricType
confCopy.AI.Knowledge.IndexType = collection.IndexType
confCopy.AI.Knowledge.IDMaxLength = collection.IDMaxLength
confCopy.AI.Knowledge.ContentMaxLength = collection.ContentMaxLength
confCopy.AI.Embedding.Model = collection.EmbeddingModel
confCopy.AI.Embedding.Dimensions = collection.VectorDimension
```

Run:

```bash
cd backend/go/admin
go test ./internal/ai/agent -run TestKnowledgeRun -v
```

If Ark/Milvus credentials are unavailable, record as skipped/manual in the final history; unit tests must still pass.

- [x] **Step 4: Implement import/index/delete semantics**

Index metadata filter should target either document DB id or source:

```go
metadata["document_db_id"] == "<id>"
```

Fallback for old data:

```go
metadata["_source"] == "<source>"
```

On successful indexing:

```go
vector_status = "indexed"
vector_id = strings.Join(ids, ",")
last_indexed_at = time.Now().UnixMilli()
indexing_error = ""
```

On failure:

```go
vector_status = "failed"
indexing_error = sanitizedError(err)
last_indexed_at = 0
```

Run:

```bash
cd backend/go/admin
go test ./internal/appsvc -run 'KnowledgeDocumentIndexer' -v
```

### Task 3: Wire Document Handler to the indexing service

**Files:**

- Modify: `backend/go/admin/internal/router/logic/knowledge_document.go`
- Modify: `backend/go/admin/internal/router/logic/knowledge_document_test.go`
- Modify: `backend/go/admin/internal/router/auth_router/knowledge.go`
- Generate: `backend/go/admin/internal/mock/mock_knowledge_document_indexer.go`

- [x] **Step 1: Write failing Handler tests**

Add tests:

```go
func TestKnowledgeDocumentHandler_CreateIndexesDocument(t *testing.T) {}
func TestKnowledgeDocumentHandler_UpdateReindexesWhenContentChanges(t *testing.T) {}
func TestKnowledgeDocumentHandler_UpdateDoesNotReindexWhenOnlyRemarkChanges(t *testing.T) {}
func TestKnowledgeDocumentHandler_DelDeletesVectorsBeforeDBRecord(t *testing.T) {}
func TestKnowledgeDocumentHandler_ImportFileDelegatesToIndexer(t *testing.T) {}
```

Run:

```bash
cd backend/go/admin
go test ./internal/router/logic -run 'KnowledgeDocumentHandler.*(Index|Import|Del)' -v
```

Expected before implementation: constructor signature and `ImportFile` do not exist.

- [x] **Step 2: Update constructor and routes**

Constructor:

```go
func NewKnowledgeDocumentHandler(q *query.Query, indexer appsvc.KnowledgeDocumentIndexer) *KnowledgeDocumentHandler
```

Route:

```go
documentGroup.Post("/importFile", handler.CtxHandlerFunc(documentHandler.ImportFile))
```

- [x] **Step 3: Implement create/update/delete indexing hooks**

Rules:

- `Create`: create row first, then `IndexDocument(ctx, item.ID)`.
- `Update`: fetch before update; reindex only when `collectionID`、`documentID`、`title`、`content`、`contentType`、`source`、`metadata` changes.
- `Del`: call `DeleteDocumentVectors(ctx, req.ID)` before `Delete()`.
- If `IndexDocument` fails, return success only if the service has already marked the row failed; otherwise return `res.FailDefault`.

Run:

```bash
cd backend/go/admin
go test ./internal/router/logic -run 'KnowledgeDocumentHandler' -v
```

### Task 4: Frontend direct upload API and Document import UI

**Files:**

- Modify: `front/apps/admin-react/src/api/business/storageFile.ts`
- Modify: `front/apps/admin-react/src/api/business/storageFile.test.ts`
- Modify: `front/apps/admin-react/src/api/business/knowledgeDocument.ts`
- Create: `front/apps/admin-react/src/api/business/knowledgeDocument.test.ts`
- Modify: `front/apps/admin-react/src/routes/_app/knowledge/document.tsx`
- Create: `front/apps/admin-react/src/routes/_app/knowledge/document.test.tsx`
- Modify: `front/apps/admin-react/src/mocks/handlers/storageFile.ts`

- [x] **Step 1: Write failing frontend API tests**

Expected API calls:

```ts
await StorageFileApi.prepareUpload({
  originalName: 'runbook.md',
  contentType: 'text/markdown',
  size: 12,
  bizType: 'knowledge-document',
  bizID: 'collection-1',
})

await StorageFileApi.completeUpload({ id: 10 })
await KnowledgeDocumentApi.importFile({ collectionID: 1, fileAssetID: 10, title: 'runbook.md' })
```

Run:

```bash
cd front/apps/admin-react
pnpm vitest run src/api/business/storageFile.test.ts src/api/business/knowledgeDocument.test.ts
```

Expected before implementation: missing methods.

- [x] **Step 2: Implement `uploadDirect` helper**

Flow:

```ts
const prepared = await prepareUpload(fields)
await fetch(prepared.data.uploadURL, {
  method: prepared.data.method,
  headers: prepared.data.headers,
  body: file,
})
return await completeUpload({ id: prepared.data.asset.id })
```

Do not send auth token to MinIO presigned URL.

- [x] **Step 3: Add Document import Drawer**

UI behavior:

- Toolbar shows primary `导入文件` button only when `collectionId > 0`.
- Drawer fields: file, title, content type, remark.
- Accepted extensions: `.md,.markdown,.txt,.html`.
- Submit disables while uploading or importing.
- Success closes Drawer, shows `导入成功`, refreshes table.
- Failure leaves Drawer open and shows `导入失败`.

Run:

```bash
cd front/apps/admin-react
pnpm vitest run src/routes/_app/knowledge/document.test.tsx
```

### Task 5: Docs, generated artifacts, and verification

**Files:**

- Modify: `backend/go/admin/docs/docs.go`
- Modify: `backend/go/admin/docs/swagger.json`
- Modify: `backend/go/admin/docs/swagger.yaml`
- Modify: `front/apps/admin-react/src/auto-imports.d.ts` if auto-import generation changes it.
- Modify: `docs/develop/backend.md`
- Modify: `docs/develop/frontend.md`
- Modify: `docs/operate/reliability.md`
- Modify: `docs/operate/security.md`
- Create: `docs/records/histories/YYYY-MM/YYYYMMDD-HHmm-knowledge-document-file-import-vector-indexing.md`

- [x] **Step 1: Regenerate backend artifacts**

Run if ORM model fields changed:

```bash
cd backend/go/admin
make script-orm
```

Run Swagger generation:

```bash
cd backend/go/admin
make swagger
```

- [x] **Step 2: Run backend verification**

```bash
cd backend/go/admin
go test ./internal/services/objectstore ./internal/appsvc ./internal/router/logic ./internal/ai/agent -v
```

For full workspace:

```bash
cd backend/go
go test ./...
```

- [x] **Step 3: Run frontend verification**

```bash
cd front/apps/admin-react
pnpm vitest run src/api/business/storageFile.test.ts src/api/business/knowledgeDocument.test.ts src/routes/_app/knowledge/document.test.tsx
pnpm build
```

- [x] **Step 4: Manual integration check**

Prerequisites:

- Postgres, Redis, MinIO, Milvus are running.
- `Storage.Enabled = true` and MinIO credentials are injected locally.
- Ark embedding API key is injected locally.

Manual flow:

1. Start backend: `cd backend/go/admin && make run`.
2. Start frontend: `cd front/apps/admin-react && pnpm dev`.
3. Open `/knowledge/collection`, create or choose a collection whose vector dimension matches embedding config.
4. Open `/knowledge/document?collectionId=<id>`.
5. Click `导入文件`, choose a `.md` file, submit.
6. Confirm MinIO contains the object, `file_asset.status = active`, `knowledge_document.vector_status = indexed`.
7. Query Milvus collection and confirm metadata contains `collection_id` and `document_db_id`.
8. Delete the document and confirm related Milvus rows are removed.

## 测试计划

- `backend/go/admin/internal/services/objectstore/minio_test.go`
  - `TestMinIOEnginePresignedPutObject`：验证签发 PUT URL 时使用 bucket、object key、expiry。
  - `TestMinIOEngineStatAndGetObject`：验证 stat/get 参数传递和 not found 映射。
- `backend/go/admin/internal/appsvc/file_storage_impl_test.go`
  - `TestFileStoragePrepareDirectUploadCreatesPendingAsset`：验证 pending asset、object key、upload URL。
  - `TestFileStorageCompleteDirectUploadActivatesAfterStat`：验证对象存在且 size 一致才激活。
  - `TestFileStorageCompleteDirectUploadRejectsSizeMismatch`：验证 size 不一致不激活。
- `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl_test.go`
  - `TestKnowledgeDocumentIndexerIndexesManualDocument`：手工文档写入 Milvus 后回写 indexed。
  - `TestKnowledgeDocumentIndexerImportsActiveFileAsset`：active 文件导入创建 Document 并索引。
  - `TestKnowledgeDocumentIndexerRejectsInactiveFileAsset`：pending/deleted 文件不可导入。
  - `TestKnowledgeDocumentIndexerDeletesExistingVectorsByMetadata`：重建前先按 metadata 删除旧向量。
  - `TestKnowledgeDocumentIndexerMarksFailedWhenIndexingFails`：索引失败回写 failed 和脱敏错误。
- `backend/go/admin/internal/router/logic/knowledge_document_test.go`
  - `TestKnowledgeDocumentHandler_CreateIndexesDocument`
  - `TestKnowledgeDocumentHandler_UpdateReindexesWhenContentChanges`
  - `TestKnowledgeDocumentHandler_UpdateDoesNotReindexWhenOnlyRemarkChanges`
  - `TestKnowledgeDocumentHandler_DelDeletesVectorsBeforeDBRecord`
  - `TestKnowledgeDocumentHandler_ImportFileDelegatesToIndexer`
- `front/apps/admin-react/src/api/business/storageFile.test.ts`
  - direct upload prepare、PUT、complete 调用顺序。
- `front/apps/admin-react/src/api/business/knowledgeDocument.test.ts`
  - importFile endpoint 和 payload。
- `front/apps/admin-react/src/routes/_app/knowledge/document.test.tsx`
  - 导入按钮、文件选择、提交、成功刷新、失败保留 Drawer。

## 验证方式

命令：

```bash
cd backend/go/admin
go test ./internal/services/objectstore ./internal/appsvc ./internal/router/logic -v
make swagger
```

```bash
cd backend/go
go test ./...
```

```bash
cd front/apps/admin-react
pnpm vitest run src/api/business/storageFile.test.ts src/api/business/knowledgeDocument.test.ts src/routes/_app/knowledge/document.test.tsx
pnpm build
```

手工检查：

- 导入文件不会经过后端代理上传大文件正文，只从浏览器 PUT 到 MinIO presigned URL。
- 前端不会把登录 token、Cookie 或内部 API header 带到 MinIO presigned URL。
- 失败时 Document 行能显示 `failed` 和可读错误摘要。
- 删除 Document 后 Milvus 无对应 `document_db_id` metadata 行。

观测检查：

- 后端日志记录 request id、asset id、document id、collection name、chunk count、耗时。
- 后端日志不记录完整 presigned URL、对象内容、SecretAccessKey、Cookie、Token。
- MinIO、Milvus、Embedding 任一失败时，错误路径可从日志定位到阶段：prepare、upload complete、download object、load file、index vector、delete vector。

## 风险与缓解

- 风险：Milvus collection 维度与当前 embedding 维度不一致导致写入失败。
  - 缓解：导入前校验 `knowledge_collection.vector_dimension == conf.AI.Embedding.Dimensions`，失败时返回清晰错误并阻止索引。
- 风险：presigned PUT 上传完成但未调用 complete，产生 pending asset。
  - 缓解：本轮先禁止 pending 被业务引用；后续另补清理任务。
- 风险：同步索引导致创建/导入请求耗时较长。
  - 缓解：本轮保持同步以保证状态准确；若耗时不可接受，再以 Temporal/Asynq 做异步化。
- 风险：文件解析能力与用户预期不一致。
  - 缓解：首期 UI 和后端都限制文本类文件，错误提示明确“不支持该文件类型”。
- 风险：删除数据库记录后向量清理失败会留下脏检索数据。
  - 缓解：删除顺序固定为先清 Milvus，再软删数据库。

## 进度记录

- [x] 2026-05-25：确认范围、现状和技术约束，创建 active execution plan。
- [x] 扩展对象存储 direct upload 协议。
- [x] 新增知识文档索引服务。
- [x] Document Handler 接入索引服务和 importFile。
- [x] 前端 Document 页面接入导入文件。
- [x] 补齐文档、Swagger、history 和验证记录。

## 决策记录

- 2026-05-25：采用 MinIO presigned PUT 直传，而不是复用现有服务端代理上传；原因是用户明确要求“使用预签名上传”，且后续更适合大文件与浏览器直传。
- 2026-05-25：导入文件先落 MinIO，再由后端下载到临时文件复用现有 Eino file loader；这样能最大化复用 `knowledge_pipeline_test.go` 已验证的 Loader -> Splitter -> MilvusIndexer 流程。
- 2026-05-25：Document 删除采用先删 Milvus 再删数据库；优先避免检索到已删除文档。
- 2026-05-25：首期只支持文本类文件；PDF/Word 解析不进入本轮，避免把导入和文档解析两个问题揉在一起。

