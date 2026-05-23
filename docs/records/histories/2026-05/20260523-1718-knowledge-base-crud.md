## [2026-05-23 17:18] | Task: 实现知识库前后端 CRUD

### 🤖 Execution Context

- **Agent ID**: `Claude Code`
- **Base Model**: `Opus 4.7`
- **Runtime**: `Claude Code desktop`

### 📥 User Query

> 继续完成知识库前后端 CRUD 任务。Collection 和 Document 的完整前后端功能，前端展示 Collection 列表，点击后跳转新 Tab 展示对应 Document 列表，并在 Document 页面顶部显示所属 Collection 信息。

### 🛠 Changes Overview

**Scope:** `backend/go/admin`, `front/apps/admin-react`, `docs/`

**Key Actions:**

- **[后端 Collection Handler]**: 实现列表、详情、创建、更新、删除，含默认值填充和重复键检测。
- **[后端 Document Handler]**: 实现列表、按集合筛选、详情、创建、更新、删除，含元数据 JSON 解析。
- **[后端路由]**: 注册 `/api/knowledge/collection/*` 和 `/api/knowledge/document/*`。
- **[前端 API]**: 封装 `KnowledgeCollectionApi` 和 `KnowledgeDocumentApi`。
- **[前端页面]**: Collection 列表页（ProTable + Drawer）、Document 列表页（顶部显示 Collection 信息）。
- **[单元测试]**: 20 个后端 Handler 测试全部通过。
- **[E2E 测试]**: 新增 `knowledge.spec.ts` 覆盖页面加载、抽屉校验和跳转。
- **[初始化 SQL]**: 补充 `init.sql` 知识库菜单、API、权限规则、角色关联。
- **[前端表单校验]**: `collection.tsx` / `document.tsx` 引入 zod schema 表单校验（`useZodForm` + `superRefine`）。
- **[文档沉淀]**: 新增 `context/project/admin/experience/zod-form-validation.md` 记录前端表单校验模式。
- **[登录互踢修复]**: 修复 sa-token-go + rueidis 适配器类型不兼容（`Get` 返 `[]byte`，`assertString` 只认 `string`）；`AuthMiddleware` 追加 token 一致性校验。沉淀到 `context/project/admin/experience/auth-kickout-rueidis.md`。

### 🧠 Design Intent (Why)

- 状态切换不单独开 Switch 接口，统一走 Update（仅传 `id` + `isEnabled`），减少 API 面。
- `ListByCollection` 通过 query filter 合并 `collectionID` 条件后走 `PageWithPaging`，避免 gen query 链式调用限制。
- 前端 `isEnabled` 用 `ProFormSelect`（value: 1/0），提交时通过 `Boolean()` 转换，与后端 bool 字段对齐。
- 枚举字段使用 `sys_dict` 管理：`MetricType`、`IndexType`、`ContentType`、`VectorStatus` 改用字典而非硬编码。优势是选项可在管理后台动态维护、复用翻译/渲染基础设施、与系统字典风格一致。
- 删除 `Status` 字段：与 `IsEnabled` 语义重复（active/archived ≈ true/false），保留 `IsEnabled` 减少双状态维护成本。
- sa-token-go 的 `kickout()` 通过 `assertString()` 校验 token 值类型，但 rueidis 适配器 `Get()` 返回 `[]byte` 不兼容，导致互踢静默失败。修复适配器返回值 + AuthMiddleware 双重校验。

### 📁 Files Modified

- `backend/go/admin/internal/services/orm/models/knowledge_collection.go`
- `backend/go/admin/internal/services/orm/query/knowledge_collection.gen.go`
- `backend/go/admin/internal/services/orm/query/knowledge_collection_extend.gen.go`
- `backend/go/admin/internal/router/logic/knowledge_collection.go`
- `backend/go/admin/internal/router/logic/knowledge_document.go`
- `backend/go/admin/internal/router/auth_router/knowledge.go`
- `backend/go/admin/internal/router/auth_router/auth_router.go`
- `backend/go/admin/internal/router/logic/knowledge_collection_test.go`
- `backend/go/admin/internal/router/logic/knowledge_document_test.go`
- `backend/go/admin/internal/router/logic/sqlite_test.go`
- `front/apps/admin-react/src/api/business/knowledgeCollection.ts`
- `front/apps/admin-react/src/api/business/knowledgeDocument.ts`
- `front/apps/admin-react/src/routes/_app/knowledge.tsx`
- `front/apps/admin-react/src/routes/_app/knowledge/collection.tsx`
- `front/apps/admin-react/src/routes/_app/knowledge/document.tsx`
- `front/apps/admin-react/tests/business/knowledge.spec.ts`
- `front/apps/admin-react/src/routeTree.gen.ts`
- `front/apps/admin-react/src/domains/dict.ts`
- `backend/go/sa-token/rueidis/rueidis.go`
- `backend/go/admin/internal/fiberc/middleware/auth.go`
- `backend/go/admin/cmd/scripts/init.sql`
- `context/project/admin/INDEX.md`
- `context/project/admin/experience/zod-form-validation.md`
- `docs/build/architecture.md`
- `docs/build/backend.md`
- `docs/records/exec-plans/active/2026-05-23-knowledge-base-crud.md`
