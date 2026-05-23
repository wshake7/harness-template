## [2026-05-23 14:48] | Task: 将 Embedding Provider 从 DashScope 迁移到 Ark (Volcengine)

### 🤖 Execution Context

- **Agent ID**: `Sisyphus`
- **Base Model**: `kimi/kimi-for-coding`
- **Runtime**: `OpenCode`

### 📥 User Query

> knowledge_pipeline_test.go 报错 401 Unauthorized，API Key 不正确。用户想使用火山引擎（Volcengine）的豆包嵌入模型 doubao-embedding-vision-251215。

### 🛠 Changes Overview

**Scope:** `backend/go/admin/internal/ai`、`backend/go/admin/internal/config`、后端配置与文档

**Key Actions:**

- **[迁移 Embedding Provider]**：将 `embedder.go` 中的 `dashscope` 替换为 `ark` (Volcengine)，支持多模态模型自动检测。
- **[更新配置默认值]**：`config.go` 中 `AI.Embedding.Provider` 默认值从 `dashscope` 改为 `ark`。
- **[更新本地配置]**：`config.local.yaml` 中 Provider 改为 `ark`，Dimensions 改为 2048（匹配 doubao-embedding-vision-251215 默认输出维度）。
- **[修复测试代码]**：`knowledge_pipeline_test.go` 在测试开始时删除旧 collection，避免维度不匹配导致的 schema mismatch 错误。
- **[更新文档]**：同步更新 `BACKEND.md`、`ARCHITECTURE.md`、`RELIABILITY.md` 中的 Provider 引用和配置说明。

### 🧠 Design Intent (Why)

用户提供的 API Key 是火山引擎 Ark 格式（`ark-` 开头），与原有的阿里云 DashScope 不兼容。需要切换 Provider 以匹配用户的实际凭据。同时 doubao-embedding-vision-251215 是多模态模型，默认输出维度为 2048，与原有配置 1024 不匹配，需要一并调整。

### 📁 Files Modified

- `backend/go/admin/internal/ai/embedder/embedder.go`
- `backend/go/admin/internal/config/config.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline_test.go`
- `backend/go/admin/etc/config.local.yaml`
- `backend/go/admin/go.mod`
- `docs/BACKEND.md`
- `docs/ARCHITECTURE.md`
- `docs/RELIABILITY.md`

### ⚠️ Breaking Changes

- **配置变更**：`AI.Embedding.Provider` 默认值从 `dashscope` 改为 `ark`。
- **维度变更**：`AI.Embedding.Dimensions` 从 1024 改为 2048（匹配 doubao-embedding-vision-251215）。
- **依赖变更**：移除了 `github.com/cloudwego/eino-ext/components/embedding/dashscope`，添加了 `github.com/cloudwego/eino-ext/components/embedding/ark`。
- **Milvus Collection**：若已有 `biz` collection，需要删除后重新创建（测试代码已自动处理）。

### 🔧 Migration Guide

对于已有环境：
1. 更新 `AI.Embedding.Provider` 为 `ark`
2. 更新 `AI.Embedding.APIKey` 为火山引擎 Ark API Key
3. 更新 `AI.Embedding.Model` 为 `doubao-embedding-vision-251215`（或你的端点 ID）
4. 更新 `AI.Embedding.Dimensions` 为 `2048`
5. 删除 Milvus 中的旧 collection（如 `biz`），让代码重新创建
