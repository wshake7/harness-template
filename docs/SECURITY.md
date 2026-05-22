# 安全默认约束

这份文档把当前前后端项目的安全默认值讲清楚。仓库级依赖、SBOM 和 provenance 默认能力统一写在 `docs/SUPPLY_CHAIN_SECURITY.md`。

## 认证与授权

- 后端管理服务使用 sa-token-go 处理登录态，使用 Casbin 处理权限模型。
- 登录接口位于 `/api/account/login/pwd`，登出接口位于 `/api/account/logout`。
- 后端通过中间件区分公共接口和需要认证的业务接口。
- 前端请求层会从 Cookie 或账号 store 读取 Token，并写入 `@vp/core` 定义的 Token header。
- Cookie 当前在登录逻辑中设置 `SameSite=Lax`、`HTTPOnly=true`、`Secure=false`；共享或生产环境必须重新审视 `Secure` 和域名策略。

## 加密与请求保护

- `/api/encrypt/public/key` 可无 Token 请求，用于获取公钥。
- 登录、改密等接口会走加密中间件；前端通过 `encryptRequest` 生成加密请求，并在响应头标记为加密时解密响应体。
- 后端 Swagger 注释里保留请求超时、请求重放和请求错误的业务 code；当前 `NonceMiddleware` 在路由注册处仍处于注释状态，后续启用前要补测试和文档。
- 所有公共 API 响应协议变化都要同步更新前端 `@vp/request` 和 Swagger 注释。

## 配置与密钥

- `backend/go/admin/etc/config.yaml` 包含本地开发数据库连接串和密码，仅可作为本地默认值。
- 不要把生产数据库、Redis、Temporal、JWT、Token、私钥或第三方凭据提交到仓库。
- 共享环境应通过安全配置注入方式覆盖本地默认值，并记录配置来源和轮换方式。
- 新增环境变量必须同步更新 `docs/BACKEND.md`、`docs/FRONTEND.md` 或运行手册。
- Milvus `APIKey` 属于外部服务凭据，示例配置留空；生产环境必须通过安全配置注入方式提供，不得提交真实 token 到仓库。

## 数据与日志

- 登录日志、API 日志、用户、角色、资源、字典、语言和任务相关数据由管理后台模型持久化。
- 日志中不得写入密码、Token、私钥、完整 Cookie、数据库连接串或可恢复的加密材料。
- history、PR 描述和排障记录要脱敏，不直接贴本地密钥、原始凭据或完整敏感日志。

## 依赖与供应链

- 前端依赖通过 `pnpm-lock.yaml` 和 workspace catalog 管理。
- 后端依赖通过各 Go module 的 `go.mod`、`go.sum` 和 `backend/go/go.work` 管理。
- GitHub Actions 需要继续保持 SHA pinning；依赖审查、OSV、SBOM 和 provenance 约束见 `docs/SUPPLY_CHAIN_SECURITY.md`。
- 新增外部服务或 SDK 时，同步记录安全边界、凭据来源、网络访问和失败模式。

## 当前短板

- 本地配置仍含开发密码示例，尚未提供 `.env.example` 或统一密钥注入规范。
- Nonce/replay 防护在路由注册处尚未启用。
- 生产 Cookie、CORS、CSRF、权限模型初始化和审计日志保留策略还没有定稿。
- CI 尚未强制运行前后端安全扫描之外的项目级测试。
