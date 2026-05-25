## [2026-05-25 21:06] | Task: 实现知识文档文件导入与向量入库

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop app`

### 📥 User Query

> 根据 execution plan 实现知识库 document 文件导入、MinIO 预签名上传，以及 knowledge collection/document 与 Milvus 向量存储的实际关联。

### 🛠 Changes Overview

**Scope:** `backend/go/admin`、`front/apps/admin-react`、`docs/`

**Key Actions:**

- **[对象存储直传]**: 为 `file_asset` 增加 `pending_upload` 状态，补齐 MinIO `PresignedPutObject` / `StatObject` / `GetObject` 能力，并新增 `prepareUpload`、`completeUpload` 路由与应用服务。
- **[知识文档索引服务]**: 新增 `KnowledgeDocumentIndexer`，把文档创建、更新、导入、删除与 Milvus 向量写入/清理打通，索引流程对齐 `knowledge_pipeline_test.go` 的 loader + splitter + indexer 设计。
- **[前端导入入口]**: 管理后台知识文档页新增“导入文件”按钮和抽屉表单，前端先走 MinIO 预签名直传，再调用 `importFile` 完成建文档与向量化。
- **[验证与文档]**: 增补后端、前端测试，并同步更新后端/前端/可靠性/安全文档。

### 🧠 Design Intent (Why)

把文件对象存储、文档元数据和 Milvus 索引串成一条闭环，避免 collection/document 只停留在业务表层却没有实际向量数据；同时用预签名直传减轻后端代理大文件上传的带宽占用。

### 📁 Files Modified

- `backend/go/admin/internal/services/objectstore/engine.go`
- `backend/go/admin/internal/services/objectstore/minio.go`
- `backend/go/admin/internal/services/objectstore/minio_test.go`
- `backend/go/admin/internal/services/orm/models/file_asset.go`
- `backend/go/admin/internal/appsvc/file_storage.go`
- `backend/go/admin/internal/appsvc/file_storage_impl.go`
- `backend/go/admin/internal/appsvc/file_storage_impl_test.go`
- `backend/go/admin/internal/appsvc/knowledge_document_indexer.go`
- `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl.go`
- `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl_test.go`
- `backend/go/admin/internal/router/logic/storage_file.go`
- `backend/go/admin/internal/router/logic/storage_file_test.go`
- `backend/go/admin/internal/router/logic/knowledge_document.go`
- `backend/go/admin/internal/router/logic/knowledge_document_test.go`
- `backend/go/admin/internal/router/auth_router/storage.go`
- `backend/go/admin/internal/router/auth_router/knowledge.go`
- `backend/go/admin/internal/router/auth_router/auth_router.go`
- `backend/go/admin/cmd/scripts/init.sql`
- `front/apps/admin-react/src/api/business/storageFile.ts`
- `front/apps/admin-react/src/api/business/storageFile.test.ts`
- `front/apps/admin-react/src/api/business/knowledgeDocument.ts`
- `front/apps/admin-react/src/api/business/knowledgeDocument.test.ts`
- `front/apps/admin-react/src/components/business/storage/fileUpload.tsx`
- `front/apps/admin-react/src/components/business/knowledge/documentImport.test.ts`
- `front/apps/admin-react/src/routes/_app/knowledge/document.tsx`
- `front/apps/admin-react/src/stores/resourceMenu.test.ts`
- `front/apps/admin-react/src/utils/zod.test.ts`
- `docs/develop/backend.md`
- `docs/develop/frontend.md`
- `docs/operate/reliability.md`
- `docs/operate/security.md`
