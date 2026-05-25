# Execution Plan 使用说明

execution plan 适合用在那些超出单轮聊天上下文、需要多次推进或风险较高的任务上。

在 Harness 主流程里，execution plan 是阶段 3「设计」和阶段 4「实现」之间的重要门禁产物。完整流程见 `docs/govern/harness-process.md`。

## 什么时候该建 plan

- 任务会跨多个 commit 或多轮工作推进。
- 这次改动会影响架构、协议、数据迁移或其他高风险区域。
- 完成任务依赖阶段性验证、回滚策略或关键决策留痕。
- 可能会有多个人或多个 Agent 在一段时间内共同推进。
- 需求会触碰 `.service-matrix/dependencies.yaml` 中多个模块、服务或基础设施路径。

## 存放位置

- 进行中的 plan 放在 `docs/records/exec-plans/active/`
- 已完成的 plan 移到 `docs/records/exec-plans/completed/`
- 复用模板在 `docs/records/exec-plans/templates/execution-plan.md`
- 暂不处理但值得保留的债务放到 `docs/records/exec-plans/tech-debt-tracker.md`

## 维护要求

- 写清目标、范围、约束、风险和验证方式。
- 涉及代码、配置、接口、脚本或用户可见行为的 plan，必须单独写「测试计划」：列出要新增或修改的测试文件、每个实现切片先写哪些失败测试、运行哪些测试命令、期望失败和期望通过的标准。
- 测试计划不能只写“补测试”或“运行测试”，要能指导后续 Agent 直接落地测试代码。
- 写清影响的模块、服务、路径占位符和需要运行的门禁。
- 推进过程和关键决定要落在仓库里，不要只存在聊天记录里。
- 状态变化要同步更新。
- 过期 plan 要及时关闭、归档或清理，保证 active 目录可信。
