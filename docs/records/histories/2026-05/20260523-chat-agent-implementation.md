## [2026-05-23] | Task: 实现 Chat Agent 图编排与工具调用

### 🤖 Execution Context

- **Agent ID**: Claude Code
- **Base Model**: claude-sonnet-4-6
- **Runtime**: Claude Code CLI

### 📥 User Query

> 实现 agent 分支的 chat agent 流程，包括 ReAct 推理、RAG 检索、多工具调用和对话记忆

### 🛠 Changes Overview

**Scope:** backend/go/admin

**Key Actions:**

- **[Chat Pipeline]**: 新增 `internal/ai/agent/chat_pipeline` 包，基于 Eino `compose.Graph` 构建 ChatAgent 图编排：`InputToRag` → `MilvusRetriever` → `ChatTemplate` → `ReactAgent`，以及 `InputToChat` 注入历史与日期到模板。ReAct Agent 最大步数 25，使用 `AllPredecessor` 触发模式。
- **[ChatModel]**: 新增 `internal/ai/models` 包，封装 OpenAI 兼容 ChatModel 工厂，模型名/APIKey/BaseURL 由配置注入，支持 DeepSeek 等实现。
- **[Retriever]**: 新增 `internal/ai/retriever` 包，封装 Milvus 向量检索器，支持 IP/L2/COSINE 度量类型，输出 `id`/`content`/`metadata` 字段。
- **[Tools]**: 新增 `internal/ai/tools` 包，实现 5 个 Agent 可调用工具：
  - `query_internal_docs`：基于 Milvus RAG 检索内部文档
  - `db_crud`：直接执行 SQL 查询/写入
  - `get_current_time`：返回多精度系统时间
  - `query_log`：通过 MCP SSE 协议连接外部日志查询 Server
  - `query_prometheus_alerts`：查询 Prometheus 活跃告警并去重
- **[Memory]**: 新增 `internal/ai/mem` 包，提供多后端对话记忆（`InMemMemory`/`RedisMemory`/`DBMemory`），支持滑动窗口裁剪。
- **[ORM Models]**: 新增 `AgentSession`（会话标识/标题/状态）和 `AgentMessage`（角色/内容/Token 计数）模型，支持软删除和 JSON 扩展元数据。
- **[Config]**: 扩展 `AIConfig`，新增 `ChatModelConfig`（Name/Model/APIKey/BaseURL 数组）、`AIMemoryConfig`（Type/MaxWindowSize/TTL）、`McpURL`，以及 `ResolveModel` 按别名查找模型配置。
- **[Indexer]**: 重构 `internal/ai/indexer`，移除内联 Milvus 客户端创建，改为复用全局 `milvusc.Client`，减少重复连接。
- **[Test]**: 新增 `chat_pipeline_test.go` 端到端测试（两次连续对话验证 RAG 检索和 ReAct 推理），`knowledge_pipeline_test.go` 测试函数重命名以区分。

### 🧠 Design Intent (Why)

在已有的知识索引管线（embedding + Milvus indexer）之上构建完整的对话 Agent 能力。使用 Eino 图编排实现固定流程（用户输入 → RAG 检索 → 模板注入 → ReAct 推理），将 Milvus 检索结果注入 ChatTemplate，让模型在生成回复前获得相关上下文。工具层内聚在 `tools` 包下，每个工具通过 `utils.InferOptionableTool` 定义类型安全的输入输出。Memory 层抽象为 `Memory`/`Session` 接口，支持运行时切换后端。

### 📁 Files Modified

- `backend/go/admin/internal/ai/agent/chat_pipeline/orchestration.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/flow.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/lambda_func.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/model.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/prompt.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/retriever.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/embedding.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/tools_node.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline/types.go` (新增)
- `backend/go/admin/internal/ai/agent/chat_pipeline_test.go` (新增)
- `backend/go/admin/internal/ai/models/open_ai.go` (新增)
- `backend/go/admin/internal/ai/retriever/retriever.go` (新增)
- `backend/go/admin/internal/ai/mem/mem.go` (新增)
- `backend/go/admin/internal/ai/mem/memory_inmem.go` (新增)
- `backend/go/admin/internal/ai/mem/memory_db.go` (新增)
- `backend/go/admin/internal/ai/mem/memory_redis.go` (新增)
- `backend/go/admin/internal/ai/tools/query_internal_docs.go` (新增)
- `backend/go/admin/internal/ai/tools/db_crud.go` (新增)
- `backend/go/admin/internal/ai/tools/get_current_time.go` (新增)
- `backend/go/admin/internal/ai/tools/query_log.go` (新增)
- `backend/go/admin/internal/ai/tools/query_metrics_alerts.go` (新增)
- `backend/go/admin/internal/ai/indexer/indexer.go`
- `backend/go/admin/internal/ai/agent/knowledge_pipeline_test.go`
- `backend/go/admin/internal/config/config.go`
- `backend/go/admin/internal/services/orm/models/agent_session.go` (新增)
- `backend/go/admin/internal/services/orm/models/agent_message.go` (新增)
- `backend/go/admin/internal/services/orm/query/agent_session.gen.go` (新增)
- `backend/go/admin/internal/services/orm/query/agent_message.gen.go` (新增)
- `backend/go/admin/etc/config.yaml`
- `backend/go/admin/go.mod`
- `docs/build/architecture.md`
- `docs/build/backend.md`

### ⚠️ Known Limitations

- ChatModel 配置需要有效的 API Key，本地开发需在 `etc/config.local.yaml` 中覆盖。
- MCP 日志 Server 连接失败时 Agent 仍可运行（降级跳过该工具），日志中以 Warn 级别记录。
- InMem memory 在服务重启后丢失，生产环境应使用 Redis 或 DB 后端。
- RAG 检索当前 TopK=1，仅返回最相关的一条文档。
