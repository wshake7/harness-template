## [2026-05-26 15:15] | Task: easy-query service orm refactor

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> Java admin 的 service 查询虽然用了 easy-query，但基本都在手写 SQL，希望按 `dromara/easy-query` 文档改造成 ORM 风格写法。

### 🛠 Changes Overview

**Scope:** `backend/java/admin` application service + persistence support

**Key Actions:**

- **补 ORM 入口**: 为 `EasyQuerySupport` 增加 `queryable/updatable`，让 service 直接使用 easy-query 查询和表达式更新链。
- **批量迁移 QueryService**: 将 `system/knowledge/job/storage` 下的手写 `SELECT/UPDATE/INSERT` 改成 easy-query 的 `where/orderBy/limit/count/firstOrNull/executeRows` 写法。
- **兼容 PostgreSQL jsonb**: 对 `metadata/inputJson/resultJson` 这类字段保留 ORM 链式更新，但通过 `setSQLSegment + CAST(... AS JSONB)` 维持 Postgres 可执行性。
- **补验证测试**: 新增 `EasyQuerySupportTests`，并通过现有 `SysUser/SysRole/JobSchedule` 与整套 Java admin 测试确认行为未回归。

### 🧠 Design Intent (Why)

这次改造的目标不是只把 SQL 字符串“搬个家”，而是让 Java admin 的 service 真正回到 easy-query 的 ORM 用法上，减少重复拼 SQL、手动维护分页排序条件和后续扩展成本。同时针对 PostgreSQL `jsonb` 的类型约束保留必要的显式 cast，避免为了追求纯 ORM 写法牺牲可运行性。

### 📁 Files Modified

- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/persistence/EasyQuerySupport.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysUserQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysRoleQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysApiLogQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysLoginLogQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysDictTypeQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysDictEntryQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysLanguageTypeQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysLanguageEntryQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysResourceApiQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysResourceMenuQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysResourceMenuApiQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysRoleApiQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysRoleMenuQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/SysUserRoleQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/knowledge/KnowledgeCollectionQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/knowledge/KnowledgeDocumentQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/job/JobExecutionQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/job/JobScheduleQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/storage/FileAssetQueryService.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/infrastructure/persistence/EasyQuerySupportTests.java`
