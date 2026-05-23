## [2026-05-22 23:50] | Task: 移植 AI 知识索引基础模块

### 🤖 Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### 📥 User Query

> 将 SuperBizAgent 的 `internal/ai/embedder`、`internal/ai/indexer`、`internal/ai/loader` 模块移植到当前项目，并继续使用现有的 `github.com/milvus-io/milvus/client/v2` 客户端。

### 🛠 Changes Overview

**Scope:** `backend/go/admin/internal/ai`、`backend/go/admin/internal/config`、后端配置与文档

**Key Actions:**

- **迁移基础 AI 组件**：新增 `embedder`、`loader`、`indexer` 三个包，并把现有 `knowledge_pipeline` 改为复用这些组件。
- **适配现有 Milvus 客户端**：没有引入 `milvus-sdk-go/v2`，而是改用 Eino 官方 `components/indexer/milvus2` 组件直接复用仓库已有 `milvus/client/v2` 客户端。
- **补齐配置与依赖**：新增 `AI.Embedding`、`AI.Knowledge` 配置项，修正 `go.mod` 的 Eino 依赖，并同步后端文档。

### 🧠 Design Intent (Why)

目标是把上游项目里可复用的 AI 知识索引能力带进来，同时保持当前仓库的基础设施边界不变。`knowledge_pipeline` 继续暴露 Eino workflow 入口，Milvus 的接入则统一走仓库现有客户端和 Eino 官方 `milvus2` 组件，避免在一个服务里并存两套 Milvus SDK，也减少自定义适配代码。

### 📁 Files Modified

- `backend/go/admin/internal/ai/embedder/embedder.go`
- `backend/go/admin/internal/ai/indexer/indexer.go`
- `backend/go/admin/internal/ai/loader/loader.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/embedding.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/indexer.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline/loader.go`
- `backend/go/admin/internal/config/config.go`
- `backend/go/admin/etc/config.yaml`
- `backend/go/admin/go.mod`
- `docs/ARCHITECTURE.md`
- `docs/BACKEND.md`
