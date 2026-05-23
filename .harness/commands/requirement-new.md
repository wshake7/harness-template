# /requirement:new

用途：为新需求建立可追溯的起点。

执行步骤：

1. 收集需求名称、背景、目标、非目标和验收标准。
2. 判断是否需要 execution plan：跨多轮、跨模块、高风险或影响协议时必须创建。
3. 根据 `.service-matrix/dependencies.yaml` 判断可能影响的模块或服务。
4. 生成或更新 `docs/product/specs/` 下的需求说明。
5. 提醒后续阶段使用 `/requirement:gate-check`。

输出要求：

- 需求 ID 或文档路径。
- 当前阶段和下一步门禁。
- 已知影响面和仍待确认的问题。
