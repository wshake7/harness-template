## [2026-05-22] | Task: 接入 Milvus 基础 Service

### 🤖 Execution Context

- **Agent ID**: Sisyphus
- **Base Model**: kimi/kimi-for-coding
- **Runtime**: OpenCode

### 📥 User Query

> 按照开发计划进行开发

### 🛠 Changes Overview

**Scope:** backend/go/admin

**Key Actions:**

- **[Config]**: 新增 `MilvusConfig` 配置结构（Enabled、Address、APIKey、DBName），更新 `config.go` 和 `etc/config.yaml`
- **[Client]**: 新增 `internal/services/milvusc` 包，封装 Milvus Go SDK 客户端创建和关闭
- **[Service]**: 新增 `internal/services/milvus.go`，实现 Fiber Service 生命周期（Start、State、Terminate）
- **[Registration]**: 在 `services.New` 中注册 Milvus service
- **[Dependencies]**: 引入 `github.com/milvus-io/milvus/client/v2 v2.6.4`，排除有兼容性问题的 etcd v3.5.5
- **[Docs]**: 同步更新 `docs/BACKEND.md`、`docs/RELIABILITY.md`、`docs/SECURITY.md`

### 🧠 Design Intent (Why)

为 admin 后端添加 Milvus 向量数据库基础连接能力，纳入统一配置和 Fiber service 生命周期管理。先只做基础连接 service，不暴露业务 API，降低首轮接入范围。Milvus 默认启用，便于本地和集成环境尽早暴露外部依赖缺失问题。

### 📁 Files Modified

- `backend/go/admin/internal/config/config.go`
- `backend/go/admin/etc/config.yaml`
- `backend/go/admin/internal/services/milvusc/milvusc.go` (新增)
- `backend/go/admin/internal/services/milvus.go` (新增)
- `backend/go/admin/internal/services/init.go`
- `backend/go/admin/go.mod`
- `docs/BACKEND.md`
- `docs/RELIABILITY.md`
- `docs/SECURITY.md`

### ⚠️ Known Limitations

- 本地未启动 Milvus 时，admin 服务启动或健康检查会失败；可将 `Milvus.Enabled` 改为 `false` 关闭。
- etcd v3.5.5 与 otelgrpc v0.59.0 存在 API 不兼容问题，已通过 `go mod edit -exclude` 排除该版本。
