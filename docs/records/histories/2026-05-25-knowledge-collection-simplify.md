## 2026-05-25 | Task: 优化知识库文档导入 Temporal 调度 & 精简 KnowledgeCollection 模型

### Changes Overview

**Scope:** `appsvc`, `router/logic`, `services/orm/models`, `services/orm/query`, `front/api`, `front/routes`

**Key Actions:**

- **Temporal 调度 fallback 修复**: `enqueueIndexDocument` 原先 `h.Temporal == nil` 判断永远为 false（接口非 nil），Temporal 不可用时直接报错导致文档创建成功但索引失败。改为 `IsConnected()` 检查，Temporal 不可用时自动 fallback 到同步索引并写入 `job_execution` 记录。
- **KnowledgeCollection 模型精简**: 删除 `EmbeddingModel`、`VectorDimension`、`DocumentCount`、`IDMaxLength`、`ContentMaxLength` 五个字段。embedding 配置统一由 config 提供，仅保留 `MetricType` 和 `IndexType` 作为 per-collection Milvus 参数。
- **前端表单同步清理**: collection 创建/编辑表单、表格列、API 类型同步删除五个字段。

### Design Intent (Why)

- `EmbeddingModel` 和 `VectorDimension` 由 config 的 embedding 配置统一管理，per-collection 覆盖增加复杂度且容易出现维度不匹配的错误。
- `DocumentCount` 自始至终未被任何业务逻辑更新，是死字段。
- `IDMaxLength` 和 `ContentMaxLength` 属于 Milvus 内部参数，config 默认值已满足需求。
- Temporal fallback 确保即使 Temporal Server 未运行，文档导入仍能完成索引并产生可见的调度记录。

### Files Modified

- `backend/go/admin/internal/appsvc/temporal_service.go`
- `backend/go/admin/internal/appsvc/knowledge_document_indexer_impl.go`
- `backend/go/admin/internal/router/logic/knowledge_document.go`
- `backend/go/admin/internal/router/logic/knowledge_collection.go`
- `backend/go/admin/internal/services/orm/models/knowledge_collection.go`
- `backend/go/admin/internal/services/orm/query/knowledge_collection.gen.go`
- `backend/go/admin/internal/mock/mock_temporal_service.go`
- `front/apps/admin-react/src/api/business/knowledgeCollection.ts`
- `front/apps/admin-react/src/routes/_app/knowledge/collection.tsx`
- `docs/develop/backend.md`
