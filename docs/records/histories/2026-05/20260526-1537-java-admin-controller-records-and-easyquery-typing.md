## [2026-05-26 15:37] | Task: tighten Java admin controller boundaries

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> controller 的入参和出参必须是 record，并且 service 下的查询不要使用魔法值，参考 easy-query query bean 思路使用强类型约束。

### 🛠 Changes Overview

**Scope:** `backend/java/admin`, `docs/records/histories`

**Key Actions:**

- **Shared DTO record 化**: 将 `R`、`PageResult`、用户/角色/菜单相关 request/response DTO 改为 `record`，并同步修正 controller 和 converter test 的调用方式。
- **Easy Query 强类型改造**: 新增集中式 `EntityProxies` 和 `EasyQuerySupport` proxy 查询/更新入口，把主要 application query service 从字符串字段名切到 proxy 列访问。
- **守门测试补齐**: 增加 shared controller boundary DTO record contract 测试，以及 application service 中禁止 Easy Query 字段名魔法值的测试。

### 🧠 Design Intent (Why)

这次改动的重点是把 Java Admin 的接口边界和查询表达式都往“可约束、可重构、可被 Agent 安全修改”的方向收紧。controller 侧先把共享 DTO 明确成 record，减少弱结构 payload；service 侧去掉字符串字段名，让查询/更新在字段改名或重构时能依赖编译期约束而不是运行时踩坑。

### 📁 Files Modified

- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/persistence/EasyQuerySupport.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/persistence/EntityProxies.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/system/*.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/job/*.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/knowledge/*.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/storage/FileAssetQueryService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/dto/*.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/UserController.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/RoleController.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/MenuController.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/application/service/ApplicationServiceStrongTypingGuardTests.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/ControllerRecordContractTests.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/application/converter/ConverterTests.java`
