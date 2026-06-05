## [2026-05-26 16:02] | Task: fix Java admin login encrypted request handling

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> Java 服务端登录失败，`/api/account/login/pwd` 报 `Required request body is missing`，需要直接修好。

### 🛠 Changes Overview

**Scope:** `backend/java/admin`, `docs/records/histories`

**Key Actions:**

- **过滤器兼容前端真实头名**: `TimestampFilter` 和 `EncryptFilter` 同时兼容 `X-Request-Timestamp`，避免 Java Admin 继续依赖旧的 `X-Timestamp` 口径。
- **请求体缓存后再解密/透传**: `EncryptFilter` 改成先缓存 body，再根据签名决定解密还是原样透传，修复带加密头时 request body 被提前消费、最终在 controller 层变成 missing body 的问题。
- **公钥接口与合同文档同步**: 为 Java Admin 补上匿名可访问的 `/api/encrypt/public/key`，新增回归测试，并把 Java Admin 合同文档里的旧头名更新为当前前端实现。

### 🧠 Design Intent (Why)

这次问题的根因不在登录业务本身，而在 Java Admin 的加密协议边界有缺口：缺少匿名可访问的公钥接口，且安全过滤链仍依赖旧头名。修复重点是把这些协议入口补齐并与当前前端真实头名对齐，让登录请求能够进入正确的解密和控制器处理流程。

### 📁 Files Modified

- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter/EncryptFilter.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter/TimestampFilter.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/EncryptController.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/config/WebConfig.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/filter/SecurityFilterTests.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/AccountControllerTests.java`
- `backend/java/admin/docs/contracts/frontend-api-samples.md`
- `backend/java/admin/docs/contracts/go-admin-api-contract.md`
