# Harness 主流程

这份文档是模板级 Harness 流程和门禁的单一真相源。`AGENTS.md`、`.harness/skills/`、`.harness/commands/` 和 CI 检查都应围绕这里保持一致。

## 五阶段

| 阶段 | 目标 | 主要产物 |
| --- | --- | --- |
| 1. 初始化 | 建立项目名称、上下文入口和矩阵骨架。 | `README.md`、`context/`、`.service-matrix/dependencies.yaml` |
| 2. 需求定义 | 明确背景、目标、非目标和验收标准。 | `docs/product/specs/` 或需求系统链接 |
| 3. 设计 | 记录方案、影响面、风险、回滚和验证方式。 | `docs/records/exec-plans/active/` 或轻量设计说明 |
| 4. 实现 | 修改代码、脚本、文档和测试，并保持同源更新。 | diff、测试、文档、必要的经验记录 |
| 5. 交付 | 验证、审查、发布记录和知识沉淀。 | `make ci` 输出、`docs/records/histories/`、release note |

## 轻量门禁

| 门禁 | 位置 | 阻塞条件 | 可机读检查 |
| --- | --- | --- | --- |
| 需求门禁 | 阶段 2 末尾 | 没有目标、非目标或验收标准。 | 需求文档必须包含 `目标` 和 `验收标准` 段落。 |
| 设计门禁 | 阶段 3 末尾 | 缺少实现方案、影响面、风险或验证方式。 | 复杂任务必须有 `docs/records/exec-plans/active/` 计划。 |
| 实现入口门禁 | 阶段 4 开始 | 任务没有可执行切分，或服务矩阵路径不合法。 | `scripts/validate-service-matrix.sh` 必须通过。 |
| 交付门禁 | 阶段 5 | 未运行 CI，或代码/流程变更没有 history。 | `make ci` 必须通过，实质变更补 `docs/records/histories/`。 |

## 执行原则

- 任何影响流程、目录、脚本或规则的改动，都要同步更新本文件或它引用的上下文。
- 需求和设计可以很轻，但不能只存在聊天里。
- 门禁优先放在错误代价最低的位置，避免实现后再发现方向错误。
- 如果用户纠正了模式性错误，Agent 应主动提议沉淀到 `context/` 或经验文档。

## 变更同步检查

`scripts/harness-sync.sh` 用 staged diff 判断手动修改是否需要同步文档和上下文，并由 `.vite-hooks/pre-commit` 在提交前执行。

默认规则：

- 脚本、CI、hooks、包管理或工具链变化，需要同步流程或 CI 说明，并补 history。
- `apps/`、`packages/`、`infra/` 或 `.service-matrix/` 变化，需要同步架构、服务矩阵或项目上下文，并补 history。
- `.service-matrix/dependencies.yaml` 变化会额外运行服务矩阵校验。

手动检查当前工作区可运行 `make harness-sync`。确认为一次性例外时，可以临时设置 `HARNESS_SYNC_SKIP=1` 跳过 pre-commit 检查，但应在评审里说明原因。

## 多运行时适配

`.harness/` 是 Skill、Agent、Command 的真相源。后续如果需要支持 `.claude/`、`.codex/`、`.gemini/` 或 `.continue/`，应通过脚本从 `.harness/` 渲染镜像目录，而不是在多个运行时目录里手工维护重复规则。
