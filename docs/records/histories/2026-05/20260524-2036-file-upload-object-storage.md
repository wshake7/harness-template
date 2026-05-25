## [2026-05-24 20:36] | Task: 实现文件上传对象存储

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop + zsh`

### 📥 User Query

> 完成“文件上传与可插拔对象存储”开发计划，实现管理后台通用文件上传能力，首期接入 MinIO，并同步测试、文档、Swagger 与 history。

### 🛠 Changes Overview

**Scope:** `backend/go/admin`、`front/packages/request`、`front/apps/admin-react`、`docs/`

**Key Actions:**

- **对象存储底座**: 新增 `Storage` 配置、`objectstore.Engine` 抽象、MinIO 实现、Fiber 生命周期接入，以及 `file_asset` ORM 模型与查询产物。
- **后端上传闭环**: 新增文件上传应用服务与 `/api/storage/file/*` 路由，支持 multipart 上传、元数据落库、短期 presigned URL、删除与补偿删除。
- **前端接入**: 为请求层增加 `meta.skipEncrypt`，新增 `StorageFileApi`、上传组件与 MSW mock handler。
- **验证与文档**: 补充分层测试、重新生成 ORM query / Swagger / mocks，并更新后端、稳定性、安全文档。

### 🧠 Design Intent (Why)

先用服务端代理上传建立一条稳定、可审计、可扩展的最小闭环：它能复用现有登录态、Casbin 和审计链路，同时通过 `objectstore.Engine` 把 MinIO 与业务层隔离开，为后续接入其他 S3 兼容存储或本地文件系统预留边界。

### 📁 Files Modified

- `backend/go/admin/internal/config/config.go`
- `backend/go/admin/internal/services/objectstore/*`
- `backend/go/admin/internal/services/orm/models/file_asset.go`
- `backend/go/admin/internal/services/orm/query/file_asset*.go`
- `backend/go/admin/internal/appsvc/file_storage*.go`
- `backend/go/admin/internal/router/logic/storage_file.go`
- `backend/go/admin/internal/router/auth_router/storage.go`
- `front/packages/request/src/encryption/helpers.ts`
- `front/apps/admin-react/src/api/business/storageFile.ts`
- `front/apps/admin-react/src/components/business/storage/fileUpload.tsx`
- `front/apps/admin-react/src/mocks/handlers/storageFile.ts`
- `docs/develop/backend.md`
- `docs/operate/reliability.md`
- `docs/operate/security.md`

### ✅ Verification

- `go test ./internal/config -run 'TestStorageConfig' -v`
- `go test ./internal/services/objectstore -v`
- `go test ./internal/services/orm/models -run 'TestFileAsset' -v`
- `go test ./internal/appsvc -run 'TestFileStorage' -v`
- `go test ./internal/router/logic -run 'TestStorageFileHandler' -v`
- `make script-orm`
- `PATH="$(go env GOPATH)/bin:$PATH" go generate ./internal/appsvc`
- `PATH="$(go env GOPATH)/bin:$PATH" make swagger`
- `pnpm --filter @vp/request test -- encryption`
- `pnpm --filter admin-react test -- storageFile`
- `pnpm --filter admin-react test -- fileUpload`
- `make harness-sync`
- `cd backend/go/admin && go test ./...`

### ⚠️ Notes

- 为了完成生成步骤，当前环境额外安装了 `mockgen` 和 `swag` CLI。
- `make swagger` 在生成过程中输出了第三方依赖常量解析 warning，但最终成功生成 `docs/docs.go`、`docs/swagger.json` 和 `docs/swagger.yaml`。
- `cd backend/go && go test ./...` 在当前 `go.work` 下直接失败，报 `directory prefix . does not contain modules listed in go.work`。
- `pnpm ready` 失败，主要受仓库既有 lint/格式问题影响；其中 `swagger.json` / `swagger.yaml` 还会因为生成器输出格式与当前 ESLint 缩进规则不一致而触发额外报错。
- `make ci` 失败，原因是 CI 脚本仍检查 `docs/build/*` 路径，而当前仓库文档已经迁移到 `docs/develop/*`。
