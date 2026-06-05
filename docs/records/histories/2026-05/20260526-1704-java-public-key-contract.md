## [2026-05-26 17:04] | Task: 对齐 Java 公钥返回协议

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> Java 服务端返回的 key 和 Go 不统一，`/api/encrypt/public/key` 这个接口参考一下 Go 代码优化一下。

### 🛠 Changes Overview

**Scope:** `backend/java/admin`，`docs/contracts`

**Key Actions:**

- **统一公钥导出格式**: 将 Java 公钥接口和登录返回中的 `publicKey` 从 PEM 改为与 Go/前端一致的 base64 SPKI 字符串。
- **补测试锁定协议**: 为公开取 key 和登录成功场景增加格式断言，并同步修正依赖该格式的加密过滤器测试。
- **修正 RSA OAEP 参数**: 将 Java `rsaEncrypt/rsaDecrypt` 改为显式 `SHA-256 + MGF1-SHA-256`，与前端 WebCrypto 和 Go `rsa.EncryptOAEP/DecryptOAEP` 保持一致，解决登录请求解密失败。
- **补齐失效公钥恢复信号**: 当客户端携带服务重启前的旧公钥导致 `X-Request-Encrypted-Key` 解密失败时，Java `EncryptFilter` 现在返回 `code=5`，让前端按既有逻辑清空本地缓存并重新拉取公钥。
- **修正加密响应写回方式**: `EncryptFilter` 不再绕过 `ContentCachingResponseWrapper` 直接写原始 `response`，避免业务失败响应在加密阶段触发 `getOutputStream()`/`getWriter()` 冲突。
- **替换认证查询 stub**: 将登录链路里的 `UserLookupRepository` 从永远返回空的 stub 替换为 Easy Query 实现，按真实 `sys_user` 和 `sys_user_role/sys_role` 表查询用户与角色。
- **同步合同文档**: 更新 Java 侧接口样例和 Go 合同说明，明确 `publicKey` 为 base64 编码的 RSA SPKI 公钥。

### 🧠 Design Intent (Why)

前端公钥缓存和 WebCrypto 导入流程只接受无头尾、无换行的 base64 公钥；同时 WebCrypto 的 `RSA-OAEP` 与 Go 的 `rsa.EncryptOAEP/DecryptOAEP` 都按 SHA-256 口径工作。Java 若继续返回 PEM、依赖 JCE 默认 OAEP 参数、在旧公钥失效时不返回标准业务码、在响应加密时绕过响应 wrapper 直接写原始流，或在认证阶段继续使用永远查不到用户的 stub 仓储，都会让同一套加密协议在 Go 和 Java 之间表现不一致，导致前端无法稳定复用同一实现。把 Java 侧公钥表示、OAEP 参数、失败恢复信号、响应写回方式和认证查询都收敛到 Go 合同可以减少分支逻辑和联调摩擦。

### 📁 Files Modified

- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/infrastructure/auth/CryptoService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/filter/ServerKeyPairProvider.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/web/controller/EncryptController.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/service/AccountService.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/port/EasyQueryUserLookupRepository.java`
- `backend/java/admin/src/main/java/cn/harnesstemplate/admin/application/port/StubUserLookupRepository.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/controller/AccountControllerTests.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/web/filter/SecurityFilterTests.java`
- `backend/java/admin/src/test/java/cn/harnesstemplate/admin/application/port/EasyQueryUserLookupRepositoryTests.java`
- `backend/java/admin/docs/contracts/frontend-api-samples.md`
- `backend/java/admin/docs/contracts/go-admin-api-contract.md`
