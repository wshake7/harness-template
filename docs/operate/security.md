# 安全与供应链

这份文档把当前前后端项目的安全默认值、依赖治理、SBOM 和 provenance 默认能力统一讲清楚。

## 认证与授权

- 后端管理服务使用 sa-token-go 处理登录态，使用 Casbin 处理权限模型。
- 登录接口位于 `/api/account/login/pwd`，登出接口位于 `/api/account/logout`。
- 后端通过中间件区分公共接口和需要认证的业务接口。
- 前端请求层会从 Cookie 或账号 store 读取 Token，并写入 `@vp/core` 定义的 Token header。
- Cookie 当前在登录逻辑中设置 `SameSite=Lax`、`HTTPOnly=true`、`Secure=false`；共享或生产环境必须重新审视 `Secure` 和域名策略。

## 加密与请求保护

- `/api/encrypt/public/key` 可无 Token 请求，用于获取公钥。
- 登录、改密等接口会走加密中间件；前端通过 `encryptRequest` 生成加密请求，并在响应头标记为加密时解密响应体。
- 文件上传接口 `/api/storage/file/upload` 使用 `multipart/form-data`，保留登录态、Casbin 和语言中间件，但必须通过请求层 `meta.skipEncrypt = true` 跳过 AES body 加密。
- 对象存储直传链路中的 `/api/storage/file/prepareUpload` 与 `/api/storage/file/completeUpload` 也属于 authenticated API；prepare 只返回短期 presigned PUT URL，不返回永久公开地址，complete 必须重新校验对象存在性和大小。
- 后端 Swagger 注释里保留请求超时、请求重放和请求错误的业务 code；当前 `NonceMiddleware` 在路由注册处仍处于注释状态，后续启用前要补测试和文档。
- 所有公共 API 响应协议变化都要同步更新前端 `@vp/request` 和 Swagger 注释。

## 配置与密钥

- `backend/go/admin/etc/config.yaml` 包含本地开发数据库连接串和密码，仅可作为本地默认值。
- 不要把生产数据库、Redis、Temporal、JWT、Token、私钥或第三方凭据提交到仓库。
- 共享环境应通过安全配置注入方式覆盖本地默认值，并记录配置来源和轮换方式。
- 新增环境变量必须同步更新 `docs/develop/backend.md`、`docs/develop/frontend.md` 或运行手册。
- Milvus `APIKey` 属于外部服务凭据，示例配置留空；生产环境必须通过安全配置注入方式提供，不得提交真实 token 到仓库。
- DashScope embedding `APIKey` 属于第三方模型凭据，示例配置必须留空，并与 Milvus 凭据一样走安全注入和轮换流程。
- `Storage.MinIO.AccessKeyID` / `Storage.MinIO.SecretAccessKey` 属于对象存储凭据；示例配置必须留空，生产或共享环境通过安全配置注入，不记录到 history、截图或排障日志里。

## 数据与日志

- 登录日志、API 日志、用户、角色、资源、字典、语言和任务相关数据由管理后台模型持久化。
- 日志中不得写入密码、Token、私钥、完整 Cookie、数据库连接串或可恢复的加密材料。
- 文件上传与下载日志可以记录 request id、文件大小、content type、engine、bucket、object key 和业务标签，但不得记录文件内容、SecretKey、完整 presigned URL query 或原始 Cookie。
- history、PR 描述和排障记录要脱敏，不直接贴本地密钥、原始凭据或完整敏感日志。

## 文件访问控制

- `file_asset` 只保存元数据，不默认返回永久公开 URL。
- 文件下载或预览必须通过 `/api/storage/file/presigned` 按需签发短期 presigned URL，默认 1 小时，最大 7 天。
- 文档导入链路只允许消费 `active` 状态的 `file_asset`，避免把未完成上传或被篡改的对象直接送入知识库索引。
- 上传默认大小限制来自 `Storage.MaxUploadBytes`，当前示例值为 10 MiB；若要放宽限制，应同时评估后端代理上传的带宽和连接占用风险。
- 对象上传成功但元数据落库失败时，服务端会立即尝试补偿删除对象，避免长期遗留孤儿文件。

## 依赖与供应链

- 前端依赖通过 `pnpm-lock.yaml` 和 workspace catalog 管理。
- 后端依赖通过各 Go module 的 `go.mod`、`go.sum` 和 `backend/go/go.work` 管理。
- 新增外部服务或 SDK 时，同步记录安全边界、凭据来源、网络访问和失败模式。本轮新增对象存储依赖为 `github.com/minio/minio-go/v7`。
- 在 Pull Request 上做依赖变更审查。
- 在 PR、定时任务和手动触发时，用 OSV 对仓库中的依赖声明和 lockfile 做漏洞扫描。
- 为 release 产物生成 SBOM。
- 为 release 产物生成 build provenance attestation。
- 所有 GitHub Actions 都固定到不可变的 commit SHA，而不是漂移的版本标签。
- pnpm 保持 `trustPolicy: no-downgrade`；如果遇到已确认的传递依赖信任降级误报，必须用精确版本写入 `trustPolicyExclude`，不要关闭全局策略。

## 供应链控制项

- `actions/dependency-review-action`：阻止 PR 引入高风险依赖变更。
- `google/osv-scanner-action`：根据仓库里的依赖文件扫描已知漏洞。
- `anchore/sbom-action`：生成 SPDX 格式的 SBOM。
- `actions/attest-build-provenance`：为 release artifact 生成签名 provenance。
- `scripts/check-action-pinning.sh`：如果 workflow 里出现浮动 tag 而不是 SHA，直接让 CI 失败。
- `pnpm-workspace.yaml`：记录 pnpm trust policy 和精确版本例外。

## pnpm trust policy 例外

`trustPolicyExclude` 只能使用精确包名和版本，新增例外时要记录原因，不能扩大到包级或关闭 `trustPolicy: no-downgrade`。

| 包 | 原因 |
| --- | --- |
| `semver@6.3.1` | 由 `@babel/core` 经 `eslint-plugin-react-hooks` 传递引入；仓库保留全局 no-downgrade，只放行该精确旧版本。 |
| `@trickfilm400/rollup-plugin-off-main-thread@3.0.0-pre1` | `workbox-build@7.4.1` 传递引入；该 prerelease 版本缺少 provenance attestation，触发 `ERR_PNPM_TRUST_DOWNGRADE`，当前只放行报错精确版本以恢复安装。 |

## 供应链限制和前提

- Dependency Review 在 public repo 可以直接使用；private repo 通常需要 GitHub Advanced Security 或对应的代码安全能力。
- OSV 和 SBOM 的效果依赖仓库里存在可识别的依赖清单或 lockfile。
- 只有当 `scripts/release-package.sh` 真的代表项目的构建产物时，provenance 才真正有意义。
- OpenSSF Scorecard 默认不启用，因为新模板仓库还没有真实分支保护、release 历史和 SAST 姿态可以评分；等仓库规则配置完成后再按需加回。

## 后续建议

- 锁定并提交项目真实依赖的 lockfile。
- 定期复查 `trustPolicyExclude`，依赖上游恢复 provenance 或替换传递依赖后及时删除例外。
- 让构建过程尽量可重复、可验证。
- 如果条件允许，在部署链路里增加对 provenance 的校验。
- 把 attestation 校验继续下沉到部署平台或准入层。

## 当前短板

- 本地配置仍含开发密码示例，尚未提供 `.env.example` 或统一密钥注入规范。
- Nonce/replay 防护在路由注册处尚未启用。
- 生产 Cookie、CORS、CSRF、权限模型初始化和审计日志保留策略还没有定稿。
- CI 尚未强制运行前后端安全扫描之外的项目级测试。
