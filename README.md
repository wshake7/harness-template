# harness-template-cn

面向 Agent 协作开发的 Harness 模板。这个仓库不替代 Claude Code、Codex、Gemini CLI、Cursor 等执行工具，而是在它们上方提供一层可版本化的工程治理：上下文、流程、门禁、服务矩阵、经验沉淀和基础 CI。

## 快速开始

```bash
make ci
make init PROJECT=your-project-name
```

每轮协作从 `AGENTS.md` 开始。长期稳定的仓库知识在 `docs/` 和 `context/` 中维护：

- `docs/`：面向人和评审的说明、指南、质量记录和发布记录。
- `context/`：面向 Agent 执行的稳定上下文入口。
- `.service-matrix/`：项目、模块、服务和路径占位符的单一真相源。
- `.harness/`：可复用 Skill、Agent、Command 的源目录。

## Harness 闭环

本模板默认采用轻量五阶段流程：

1. 初始化：建立项目、上下文和矩阵骨架。
2. 需求定义：写清目标、非目标和验收标准。
3. 设计：记录方案、影响面、风险和验证方式。
4. 实现：按任务修改代码、文档和测试。
5. 交付：运行 `make ci`，补 history 或 release note，沉淀可复用经验。

流程与门禁的单一真相源是 `docs/HARNESS_PROCESS.md`。

## 常用命令

```bash
make check-docs # 检查必需文档、Harness 入口和脚本是否齐全
make check-repo # 检查文档骨架和仓库基础卫生
make ci # 运行完整模板级 CI 门禁
make new-plan SLUG=example-plan # 创建一份 execution plan
make new-history SLUG=example-change # 创建一份交付 history 记录
make new-experience SCOPE=project/example SLUG=example-lesson # 创建一条可复用经验
make validate-matrix # 校验服务/项目矩阵
```

## 设计边界

- 保持技术栈中立，不预设前端、后端、云平台或部署方式。
- 优先沉淀可机读规则，不把关键约束只留在聊天记录里。
- 默认先做小而清晰的 Harness 治理层，不重造 IDE、模型网关或通用 Agent runtime。
