# Harness 框架级上下文

框架级上下文记录所有 Agent 协作任务都要遵循的流程和门禁。

必读入口：

- `docs/HARNESS_PROCESS.md`：五阶段流程和轻量门禁。
- `.harness/commands/`：稳定命令入口。
- `.harness/skills/`：可复用工作流。
- `.service-matrix/dependencies.yaml`：项目、模块、服务和路径占位符。
- `scripts/harness-sync.sh`：基于 staged diff 的文档和上下文同步检查。

维护规则：

- 流程口径只在 `docs/HARNESS_PROCESS.md` 定义。
- Skill 和 Command 只引用流程，不重新定义另一套阶段语义。
- 新增门禁时必须同时补 CI 或脚本检查。
