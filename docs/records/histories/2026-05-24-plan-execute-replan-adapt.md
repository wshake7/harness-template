# 2026-05-24: plan_execute_replan 改造 & 周边修复

## 变更摘要

从 SuperBizAgent 参考实现迁移 `plan_execute_replan` 包到 harness-template，适配本地配置模式，并修复周边编译/运行时问题。

## 涉及文件

### 1. `plan_execute_replan` 包 — 依赖注入改造

**planner.go / executor.go / replan.go / plan_execute_replan.go**

原参考实现通过无参函数 (`models.OpenAIForDeepSeekV31Think()`, `tools.GetLogMcpTool()`) 创建模型和工具。harness-template 的模型和工具需要外部配置（API Key、MCP URL、Embedding 配置等），因此改为依赖注入模式：

- `NewPlanner(ctx, planModel)` — 接受 `model.ToolCallingChatModel`
- `NewExecutor(ctx, execModel, toolList)` — 接受模型和工具列表
- `NewRePlanAgent(ctx, chatModel)` — 接受 `model.ToolCallingChatModel`
- `BuildPlanAgent(ctx, planModel, execModel, replanModel, tools, query)` — 接受所有依赖

同时将 `prints.Event(event)` 替换为 `fmt.Printf("%+v\n", event)`（`prints` 来自 `eino-examples`，不在本项目依赖中）。

### 2. `models/open_ai.go` — 新增 ExtraFields

`OpenAIChatModelConfig` 新增 `ExtraFields map[string]any` 字段，透传到 `openai.ChatModelConfig.ExtraFields`。

用途：DeepSeek thinking 模型 (deepseek-reasoner / deepseek-v3.1) 不支持 `tool_choice`，导致 planner 报 "Thinking mode does not support this tool_choice"。通过 `ExtraFields: {"thinking": {"type": "disabled"}}` 禁用思考模式来兼容。

### 3. `services/httpc/httpc.go` — 默认初始化

`var Client *HttpcClient` → `var Client = &HttpcClient{resty.New()}`

原因：`NewPrometheusAlertsQueryTool` 内部调用 `httpc.Client.R()`。在测试中 `httpc.New()` 未被调用，导致空指针 panic。现在默认即可用，`New(logger)` 仍可覆盖。

### 4. 测试文件

- **plan_execute_test.go** — 测试完整 Plan-Execute-Replan 流程，加载配置、创建模型和工具、调用 `BuildPlanAgent`
- **tool_call_test.go** — 新增，测试 ChatModel 工具绑定和链式调用（对应 SuperBizAgent `llm_tool_cmd/main.go`）
- **recall_test.go** — 测试 RAG 检索

### 5. chat_pipeline/prompt.go — 微小修正

system prompt 中多余空格修复。
