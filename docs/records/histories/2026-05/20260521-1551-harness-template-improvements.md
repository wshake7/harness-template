## [2026-05-21 15:51] | Task: 完善 Harness 模板闭环

### Execution Context

- **Agent ID**: `Codex`
- **Base Model**: `GPT-5`
- **Runtime**: `Codex desktop`

### User Query

> 参考 QQ 音乐 Harness Engineering 实践，把当前模板升级为可运行、可校验、可迁移的 Harness 治理层，并实现改进清单。

### Changes Overview

**Scope:** 仓库级 Harness 模板、文档、脚本和 CI 检查。

**Key Actions:**

- **补齐自检闭环**: 新增 README、根级 SECURITY 和 CODEOWNERS，修复模板基础 CI 缺口。
- **新增 Harness 治理层**: 增加五阶段主流程、三层 `context/`、服务矩阵和 `.harness/` Skill/Command 骨架。
- **增加机械校验**: 新增服务矩阵校验、初始化 smoke test 和经验记录脚本，并接入 `make ci`。
- **增加变更同步检查**: 新增 `scripts/harness-sync.sh`，并通过 `.vite-hooks/pre-commit` 在提交前检查 staged diff 是否需要同步文档和上下文。
- **修复依赖安装策略**: 为 pnpm `trustPolicy: no-downgrade` 增加 `semver@6.3.1` 精确版本例外，保留全局供应链保护。
- **修复本地 Markdown lint**: 新增 markdownlint-cli2 配置，避免安装依赖后扫描 `node_modules`。
- **同步质量与债务**: 将质量评分改为可行动雷达，并把低分项同步到技术债追踪。

### Design Intent

让模板不只是一组文档入口，而是能把 Agent 协作中的流程、上下文、路径、经验和交付检查沉淀为可版本化、可审计、可持续演进的工程资产。

### Files Modified

- `README.md`
- `docs/HARNESS_PROCESS.md`
- `context/`
- `.service-matrix/dependencies.yaml`
- `.harness/`
- `scripts/validate-service-matrix.sh`
- `scripts/check-init-project.sh`
- `scripts/new-experience.sh`
- `scripts/harness-sync.sh`
- `.vite-hooks/pre-commit`
- `pnpm-workspace.yaml`
- `.markdownlint-cli2.yaml`
- `docs/SUPPLY_CHAIN_SECURITY.md`
- `docs/QUALITY_SCORE.md`
- `docs/exec-plans/tech-debt-tracker.md`
